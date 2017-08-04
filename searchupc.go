package qiaoqiao

import (
	"strings"
)

type SearchUpcResult struct {
	code   string
	Result *SearchUpcResultItem `json:"0"`
}

type SearchUpcResultItem struct {
	ProductName string  `json:"productname"`
	ImageUrl    string  `json:"imageurl"`
	ProductUrl  string  `json:"producturl"`
	Price       string  `json:"price"`
	Currency    string  `json:"currency"`
	SalePrice   string  `json:"saleprice"`
	StoreName   string  `json:"storename"`
}

func (p *SearchUpcResult) getStatus() (status int) {
	if p.Result == nil || strings.Trim(p.Result.ProductName, "") == "" {
		status = StatusRequestUnsuccessfully
	} else {
		status = StatusRequestSuccessfully
	}
	return
}

func (p *SearchUpcResult) getProduct() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Result.ProductName
}
func (p *SearchUpcResult) getDescription() (desc string) {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	desc = p.getProduct() + "\n" + p.Result.ProductUrl + "\n" + p.Result.Price + " " + p.Result.Currency
	return desc
}
func (p *SearchUpcResult) getPeople() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Result.StoreName
}
func (p *SearchUpcResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return "http://www.searchupc.com/drawupc.aspx?q=" + p.code
}
func (p *SearchUpcResult) getCompany() Company {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return Company{"", ""}
	}
	return Company{
		p.Result.StoreName,
		"",
	}
}
