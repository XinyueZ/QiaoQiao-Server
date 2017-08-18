package qiaoqiao

import (
	"strconv"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type EANdataAttributes struct {
	Product         string `json:"product"`
	Description     string `json:"description"`
	LongDescription string `json:"long_desc"`
	Language        string `json:"language_text"`
	LanguageText    string `json:"language_text_long"`
	Price           string `json:"price_new"`
	PriceUnit       string `json:"price_new_extra"`
	Author          string `json:"author"`
}

type Barcode struct {
	Url string `json:"EAN13"`
}

type Company struct {
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type EANdataStatus struct {
	Code string `json:"code"`
}

type EANProduct struct {
	Attributes *EANdataAttributes `json:"attributes"`
	EAN13      string            `json:"EAN13"`
	ISBN10     string            `json:"ISBN10"`
	Barcode    Barcode           `json:"barcode"`
	Image      string            `json:"image"`
}

type EANdataResult struct {
	Status  EANdataStatus `json:"status"`
	Product EANProduct    `json:"product"`
	Company Company       `json:"company"`
}

func (p *EANdataResult) parse(productQuery *ProductQuery) IProductResult {
	cxt := appengine.NewContext(productQuery.r)
	chBytes := make(chan []byte)
	go get(productQuery.r, fmt.Sprintf(productQuery.targetUrl, productQuery.params.Keyword, productQuery.key), chBytes)
	byteArray := <-chBytes
	log.Infof(cxt, fmt.Sprintf("%s feeds %s", productQuery.name, string(byteArray)))
	json.Unmarshal(byteArray, p)
	return p
}

func (p *EANdataResult) getStatus() (status int) {
	status, _ = strconv.Atoi(p.Status.Code)
	switch status {
	case 0, 404:
		status = StatusRequestUnsuccessfully
	default:
		if p.Product.Attributes == nil {
			status = StatusRequestUnsuccessfully
		} else {
			status = StatusRequestSuccessfully
		}
	}
	return
}

func (p *EANdataResult) getProduct() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Product.Attributes.Product
}
func (p *EANdataResult) getDescription() (desc string) {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	if p.Product.Attributes.LongDescription != "" {
		desc = p.Product.Attributes.LongDescription
	} else {
		desc = p.Product.Attributes.Description
	}
	return desc
}
func (p *EANdataResult) getPeople() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Product.Attributes.Author
}
func (p *EANdataResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Product.Barcode.Url
}
func (p *EANdataResult) getCompany() Company {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return Company{"", ""}
	}
	return p.Company
}

func (p *EANdataResult) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	if p.getStatus() == StatusRequestSuccessfully {
		pi := ProductImage{make([]string, 0), make([]string, 0), make([]string, 0), "", "eandata"}
		pi.Medium = append(pi.Medium, p.Product.Image)
		pi.Thumbnail = p.Product.Image
		imageList = append(imageList, pi)
	}
	return
}
