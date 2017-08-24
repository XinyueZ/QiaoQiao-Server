package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine/log"
	"net/http"
)

type TescoResult struct {
	TescoProduct []TescoProduct  `json:"products"`
	DataSource string
}

type TescoProduct struct {
	Upc              string `json:"gtin"`
	TescoId          string `json:"catId"`
	ShortDescription *string `json:"description"`
	LongDescription  *string `json:"marketingText"`
	Brand            *string `json:"brand"`
}


func (p *TescoResult) parse(productQuery *ProductQuery) IProductResult {
	cxt := appengine.NewContext(productQuery.r)
	chBytes := make(chan []byte)
	header:=make(http.Header)
	header.Add("Ocp-Apim-Subscription-Key", productQuery.key)
	go getWithHeader(productQuery.r, fmt.Sprintf(productQuery.targetUrl, productQuery.params.Keyword), chBytes, &header)
	byteArray := <-chBytes
	log.Infof(cxt, fmt.Sprintf("%s feeds %s", productQuery.name, string(byteArray)))
	json.Unmarshal(byteArray, p)
	p.DataSource = productQuery.name
	return p
}

func (p *TescoResult) getStatus() (status int) {
	if p.TescoProduct == nil || len(p.TescoProduct) == 0 {
		status = StatusRequestUnsuccessfully
	} else {
		status = StatusRequestSuccessfully
	}
	return
}

func (p *TescoResult) getProduct() string {
	if p.getStatus() == StatusRequestSuccessfully {
		if p.TescoProduct[0].ShortDescription != nil {
			return *p.TescoProduct[0].ShortDescription
		}
	}
	return ""
}
func (p *TescoResult) getDescription() (desc string) {
	if p.getStatus() == StatusRequestSuccessfully {
		if p.TescoProduct[0].LongDescription != nil {
			return *p.TescoProduct[0].LongDescription
		}
	}
	return ""
}
func (p *TescoResult) getPeople() (people string) {
	return ""
}

func (p *TescoResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return generateBarcodeUrl(p.TescoProduct[0].Upc)
	}
	return ""
}

func (p *TescoResult) getCompany() Company {
	if p.getStatus() == StatusRequestSuccessfully {
		var companyName string = ""
		if p.TescoProduct[0].Brand != nil {
			companyName = *p.TescoProduct[0].Brand
		}

		return Company{
			companyName,
			"",
		}
	}
	return Company{"", ""}
}

func (p *TescoResult) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	return
}

