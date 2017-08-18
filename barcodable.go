package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine/log"
)

type BarcodableResult struct {
	codeType string
	Status   int         `json:"status"`
	Message  string         `json:"message"`
	Item     BarcodableItem `json:"item"`
}

type BarcodableItem struct {
	Ean   string           `json:"ean"`
	Isbn  string           `json:"isbn"`
	Asins []BarcodableAsin `json:"asins"`
}

type BarcodableAsin struct {
	Asin         string   `json:"asin"`
	Title        string   `json:"title"`
	Mpn          string   `json:"mpn"`
	PartNumber   string   `json:"part_number"`
	Brand        string   `json:"brand"`
	Manufacturer string   `json:"manufacturer"`
	Url          string   `json:"url"`
	Images       []string `json:"images"`
	Categories   []string `json:"categories"`
}

func (p *BarcodableResult) parse(productQuery *ProductQuery) IProductResult {
	res := p.parseInternal(productQuery)

	if res.getStatus() == StatusRequestUnsuccessfully {
		newCodeType := p.codeType
		if newCodeType == "upc" {
			newCodeType = "ean"
		} else {
			newCodeType = "upc"
		}
		p.setCodeType(newCodeType)
		res = p.parseInternal(productQuery)
	}
	return res
}

func (p *BarcodableResult) parseInternal(productQuery *ProductQuery) IProductResult {
	cxt := appengine.NewContext(productQuery.r)
	chBytes := make(chan []byte)
	go get(productQuery.r, fmt.Sprintf(productQuery.targetUrl, p.codeType, productQuery.params.Keyword), chBytes)
	byteArray := <-chBytes
	log.Infof(cxt, fmt.Sprintf("%s feeds %s", productQuery.name, string(byteArray)))
	json.Unmarshal(byteArray, p)
	return p
}

func (p *BarcodableResult) setCodeType(codeType string) (*BarcodableResult) {
	p.codeType = codeType
	return p
}

func (p *BarcodableResult) getStatus() (status int) {
	switch p.Status {
	case 200:
		if len(p.Item.Asins) > 0 {
			status = StatusRequestSuccessfully
		} else {
			status = StatusRequestUnsuccessfully
		}
	default:
		status = StatusRequestUnsuccessfully
	}
	return
}

func (p *BarcodableResult) getProduct() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Item.Asins[0].Title
}
func (p *BarcodableResult) getDescription() (desc string) {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	desc = p.getProduct()
	for _, s := range p.Item.Asins[0].Categories {
		desc = "\n" + desc + s + "\n"
	}
	return desc
}
func (p *BarcodableResult) getPeople() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return ""
}
func (p *BarcodableResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return generateBarcodeUrl(p.Item.Ean)
}
func (p *BarcodableResult) getCompany() Company {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return Company{"", ""}
	}
	return Company{
		p.Item.Asins[0].Manufacturer + "\n" + p.Item.Asins[0].Brand,
		"",
	}
}

func (p *BarcodableResult) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	if p.getStatus() == StatusRequestSuccessfully {
		if p.Item.Asins[0].Images != nil && len(p.Item.Asins[0].Images) > 0 {
			pi := ProductImage{make([]string, 0),  make([]string, 0), make([]string, 0), "", "barcodable"}
			for _, element := range p.Item.Asins[0].Images {
				pi.Medium = append(pi.Medium, element)
				pi.Thumbnail = element
				imageList = append(imageList, pi)
			}
		}
	}
	return
}
