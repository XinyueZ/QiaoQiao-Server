package qiaoqiao

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"net/http"
)

type ProductUpc struct {
	r         *http.Request
	targetUrl string
}

type ProductResult interface {
	getStatus() int
	getProduct() string
	getDescription() string
	getPeople() string
	getBarcodeUrl() string
	getCompany() Company
}

func newProductUpc(r *http.Request, targetUrl string) (p *ProductUpc) {
	p = new(ProductUpc)
	p.r = r
	p.targetUrl = targetUrl
	return
}

func (p *ProductUpc) get(language string, code string, response chan []byte, service string) {
	switch service {
	case "eandata":
		get(p.r, fmt.Sprintf(p.targetUrl, code, EANDATE_KEY), response)
	case "aws":
		for _, assoc := range AWS_ASSOCIATE_LIST {
			get(p.r, getAWSapi(p, &assoc, code), response)
		}
	}
}

type ProductUpcResponse struct {
	r              *http.Request
	ProductUpcItem []*ProductUpcItem `json:"result"`
}

func newProductUpcResponse(r *http.Request) (p *ProductUpcResponse) {
	p = new(ProductUpcResponse)
	p.r = r
	p.ProductUpcItem = []*ProductUpcItem{}
	return
}

func (p *ProductUpcResponse) show(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(p)
	if err == nil {
		fmt.Fprintf(w, "%s", bytes)
	} else {
		NewStatus(w, "noid", StatusRequestUnsuccessfully, "Can't give you UPC information.").show(appengine.NewContext(p.r))
	}
}

type ProductUpcItem struct {
	Status      int     `json:"status"`
	Product     string  `json:"product"`
	Description string  `json:"description"`
	Barcode     string  `json:"barcodeSource"`
	Company     Company `json:"company"`
	People      string  `json:"people"`
	Source      string  `json:"source"`
}

func newProductUpcItem(result ProductResult, source string) (item *ProductUpcItem) {
	item = new(ProductUpcItem)
	item.Source = source
	item.Status = result.getStatus()
	item.Product = result.getProduct()
	item.Description = result.getDescription()
	item.People = result.getPeople()
	item.Barcode = result.getBarcodeUrl()
	item.Company = result.getCompany()
	return
}
