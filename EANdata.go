package qiaoqiao

import (
	"strconv"
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

type ENAdataStatus struct {
	Code string `json:"code"`
}

type Product struct {
	Attributes EANdataAttributes `json:"attributes"`
	EAN13      string            `json:"EAN13"`
	ISBN10     string            `json:"ISBN10"`
	Barcode    Barcode           `json:"barcode"`
	Image      string            `json:"image"`
}

type EANdataResult struct {
	Status  ENAdataStatus `json:"status"`
	Product Product       `json:"product"`
	Company Company       `json:"company"`
}

func (p *EANdataResult) getStatus() (status int) {
	status, _ = strconv.Atoi(p.Status.Code)
	switch status {
	case 0, 404:
		status = StatusRequestUnsuccessfully
	default:
		status = StatusRequestSuccessfully
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