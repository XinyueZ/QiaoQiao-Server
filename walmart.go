package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine/log"
)

type WalmartResult struct {
	Items  []WalmartItem `json:"items"`
	Errors []WalmartError `json:"errors"`
}

type WalmartItem struct {
	Name             string  `json:"name"`
	Upc              string  `json:"upc"`
	LongDescription  *string `json:"longDescription"`
	ShortDescription *string `json:"shortDescription"`
	Brand            *string `json:"brandName"`
	Thumbnail        *string  `json:"thumbnailImage"`
	MediumImage      *string  `json:"mediumImage"`
	LargeImage       *string  `json:"largeImage"`
	ProductUrl       *string  `json:"productUrl"`
	Rating           string   `json:"customerRating"`
	ImageEntities    []WalmartImage  `json:"imageEntities"`
}

type WalmartImage struct {
	Thumbnail   *string  `json:"thumbnailImage"`
	MediumImage *string  `json:"mediumImage"`
	LargeImage  *string  `json:"largeImage"`
}

type WalmartError struct {
	Code    int   `json:"code"`
	Message string   `json:"message"`
}

func (p *WalmartResult) parse(productQuery *ProductQuery) IProductResult {
	cxt := appengine.NewContext(productQuery.r)
	chBytes := make(chan []byte)
	go get(productQuery.r, fmt.Sprintf(productQuery.targetUrl, productQuery.key, productQuery.params.Keyword), chBytes)
	byteArray := <-chBytes
	log.Infof(cxt, fmt.Sprintf("%s feeds %s", productQuery.name, string(byteArray)))
	json.Unmarshal(byteArray, p)
	return p
}

func (p *WalmartResult) getStatus() (status int) {
	if p.Errors != nil || p.Items == nil || len(p.Items) < 1 {
		status = StatusRequestUnsuccessfully
	} else {
		status = StatusRequestSuccessfully
	}
	return
}

func (p *WalmartResult) getProduct() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items[0].Name
	}
	return ""
}
func (p *WalmartResult) getDescription() (desc string) {
	if p.getStatus() == StatusRequestSuccessfully {
		if p.Items[0].LongDescription == nil {
			if p.Items[0].ShortDescription != nil {
				return *p.Items[0].ShortDescription
			} else {
				return ""
			}
		} else {
			return *p.Items[0].LongDescription
		}
	}
	return ""
}
func (p *WalmartResult) getPeople() (people string) {
	return ""
}

func (p *WalmartResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return generateBarcodeUrl(p.Items[0].Upc)
	}
	return ""
}

func (p *WalmartResult) getCompany() Company {
	if p.getStatus() == StatusRequestSuccessfully {
		var companyName string = ""
		if p.Items[0].Brand != nil {
			companyName = *p.Items[0].Brand
		}
		var url string = ""
		if p.Items[0].ProductUrl != nil {
			url = *p.Items[0].ProductUrl
		}
		return Company{
			companyName,
			url,
		}
	}
	return Company{"", ""}
}

func (p *WalmartResult) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	if p.getStatus() == StatusRequestSuccessfully {
		if p.Items[0].ImageEntities != nil && len(p.Items[0].ImageEntities) > 0 {
			pi := ProductImage{make([]string, 0), make([]string, 0), make([]string, 0), "", "aws"}
			for _, element := range p.Items[0].ImageEntities {
				if element.Thumbnail != nil {
					pi.Small = append(pi.Small, *element.Thumbnail)
					pi.Thumbnail = *element.Thumbnail
				}
				if element.MediumImage != nil {
					pi.Medium = append(pi.Medium, *element.MediumImage)
				}
				if element.LargeImage != nil {
					pi.Large = append(pi.Large, *element.LargeImage)
				}
				imageList = append(imageList, pi)
			}
		} else {
			pi := ProductImage{make([]string, 0), make([]string, 0), make([]string, 0), "", "aws"}
			if p.Items[0].Thumbnail != nil {
				pi.Thumbnail = *p.Items[0].Thumbnail
				pi.Small = append(pi.Small, *p.Items[0].Thumbnail)
			}
			if p.Items[0].LargeImage != nil {
				pi.Large = append(pi.Large, *p.Items[0].LargeImage)
			}
			if p.Items[0].MediumImage != nil {
				pi.Medium = append(pi.Medium, *p.Items[0].MediumImage)
			}
			imageList = append(imageList, pi)
		}
	}
	return
}
