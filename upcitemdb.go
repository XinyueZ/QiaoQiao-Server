package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine/log"
)

type UpcItemDbResult struct {
	Code   string `json:"code"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Items  []UpcItemDbItem `json:"items"`
}

type UpcItemDbItem struct {
	Ean         string `json:"ean"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Isbn        string `json:"isbn"`
	Publisher   string `json:"publisher"`
	Images      []string `json:"images"`
}

func (p *UpcItemDbResult) parse(productQuery *ProductQuery) IProductResult {
	cxt := appengine.NewContext(productQuery.r)
	chBytes := make(chan []byte)
	go get(productQuery.r, fmt.Sprintf(productQuery.targetUrl, productQuery.params.Keyword), chBytes)
	byteArray := <-chBytes
	log.Infof(cxt, fmt.Sprintf("%s feeds %s", productQuery.name, string(byteArray)))
	json.Unmarshal(byteArray, p)
	return p
}

func (p *UpcItemDbResult) getStatus() (status int) {
	if p.Code == "INVALID_UPC" || p.Total == 0 || len(p.Items) < 1 {
		status = StatusRequestUnsuccessfully
	} else {
		status = StatusRequestSuccessfully
	}
	return
}

func (p *UpcItemDbResult) getProduct() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Items[0].Title
}
func (p *UpcItemDbResult) getDescription() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Items[0].Description
}
func (p *UpcItemDbResult) getPeople() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return ""
}
func (p *UpcItemDbResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return generateBarcodeUrl(p.Items[0].Ean)
}
func (p *UpcItemDbResult) getCompany() Company {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return Company{"", ""}
	}
	return Company{p.Items[0].Publisher, ""}
}

func (p *UpcItemDbResult) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	if p.getStatus() == StatusRequestSuccessfully {
		if p.Items[0].Images != nil && len(p.Items[0].Images) > 0 {
			pi := ProductImage{make([]string, 0), "", "aws"}
			for _, element := range p.Items[0].Images {
				pi.Url = append(pi.Url, element)
				pi.Thumbnail = element
				imageList = append(imageList, pi)
			}
		}
	}
	return
}
