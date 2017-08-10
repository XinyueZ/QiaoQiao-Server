package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine/log"
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

func (p *SearchUpcResult) setCode(code string) IProductResult {
	p.code = code
	return p
}

func (p *SearchUpcResult) parse(productQuery *ProductQuery) IProductResult {
	cxt := appengine.NewContext(productQuery.r)
	chBytes := make(chan []byte)
	go get(productQuery.r, fmt.Sprintf(productQuery.targetUrl, productQuery.params.Keyword, productQuery.key), chBytes)
	byteArray := <-chBytes
	log.Infof(cxt, fmt.Sprintf("%s feeds %s", productQuery.name, string(byteArray)))
	json.Unmarshal(byteArray, p)
	return p
}

func (p *SearchUpcResult) getStatus() (status int) {
	if p.Result == nil || p.Result.ProductName == " " {
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
	return generateBarcodeUrl(p.code)
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

func (p *SearchUpcResult) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	return
}
