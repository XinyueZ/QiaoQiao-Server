package qiaoqiao

import (
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"google.golang.org/appengine"
)

type ProductUpc struct {
	r         *http.Request
	targetUrl string
}

func newProductUpc(r *http.Request, targetUrl string) (p *ProductUpc) {
	p = new(ProductUpc)
	p.r = r
	p.targetUrl = targetUrl
	return
}

func (p *ProductUpc) get(language string, code string, response chan []byte) {
	url := fmt.Sprintf(p.targetUrl, code, eandateKey)
	get(p.r, url, response)
}

type ProductUpcResponse struct {
	r           *http.Request
	Status      int  `json:"status"`
	Product     string  `json:"product"`
	Description string `json:"description"`
	Barcode     string `json:"description"`
	Company     Company `json:"company"`
	People      string `json:"people"`
	Source      string `json:"source"`
}

func newProductUpcResponse(r *http.Request, eandata *EANdataResult) (p *ProductUpcResponse) {
	p = new(ProductUpcResponse)

	p.Source = "eandata"
	p.r = r
	p.Status, _ = strconv.Atoi(eandata.Status.Code)
	if p.Status  == 404 {
		p.Status = StatusRequestUnsuccessfully
	 	return
	}

	p.Product = eandata.Product.Attributes.Product
	if eandata.Product.Attributes.LongDescription != "" {
		p.Description = eandata.Product.Attributes.LongDescription
	} else {
		p.Description = eandata.Product.Attributes.Description
	}
	p.People = eandata.Product.Attributes.Author
	p.Barcode = eandata.Product.Barcode.Url
	p.Company = eandata.Company
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
