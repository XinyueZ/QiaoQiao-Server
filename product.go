package qiaoqiao

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

type ProductQuery struct {
	r               *http.Request
	params          *Parameter
	targetUrl       string
	key             string
	name            string
	productResponse *ProductResponse
}

type IProductResult interface {
	getStatus() int
	getProduct() string
	getDescription() string
	getPeople() string
	getBarcodeUrl() string
	getCompany() Company
}

func newProductQuery(r *http.Request, params *Parameter, targetUrl string, key string, name string) (p *ProductQuery) {
	p = new(ProductQuery)
	p.r = r
	p.params = params
	p.targetUrl = targetUrl
	p.key = key
	p.name = name
	p.productResponse = newProductResponse(r)
	return
}

func (p *ProductQuery) get(language string, code string, key string, service string) (result IProductResult) {
	cxt := appengine.NewContext(p.r)
	result = nil
	switch service {
	case "eandata":
		result = new(EANdataResult)
		chBytes := make(chan []byte)
		go get(p.r, fmt.Sprintf(p.targetUrl, code, key), chBytes)
		byteArray := <-chBytes
		log.Infof(cxt, fmt.Sprintf("eandata feeds %s", string(byteArray)))
		json.Unmarshal(byteArray, result)
	case "searchupc":
		searchUpcResult := new(SearchUpcResult)
		searchUpcResult.code = p.params.Keyword
		result = searchUpcResult
		chBytes := make(chan []byte)
		go get(p.r, fmt.Sprintf(p.targetUrl, code, key), chBytes)
		byteArray := <-chBytes
		log.Infof(cxt, fmt.Sprintf("searchupc feeds %s", string(byteArray)))
		json.Unmarshal(byteArray, result)
	case "barcodable":
		result = new(BarcodableResult)
		chBytes := make(chan []byte)
		go get(p.r, fmt.Sprintf(p.targetUrl, code), chBytes)
		byteArray := <-chBytes
		log.Infof(cxt, fmt.Sprintf("barcodable feeds %s", string(byteArray)))
		json.Unmarshal(byteArray, result)
	case "upcitemdb":
		result = new(UpcItemDbResult)
		chBytes := make(chan []byte)
		go get(p.r, fmt.Sprintf(p.targetUrl, code), chBytes)
		byteArray := <-chBytes
		log.Infof(cxt, fmt.Sprintf("upcitemdb feeds %s", string(byteArray)))
		json.Unmarshal(byteArray, result)
	}
	return
}

func (p *ProductQuery) search() (prdResp *ProductResponse) {
	productViewModel := newProductViewModel(p.get(p.params.Language, p.params.Keyword, p.key, p.name), p.name)
	if productViewModel.Status == StatusRequestSuccessfully {
		p.productResponse.ProductViewModels = append(p.productResponse.ProductViewModels, productViewModel)
	}
	prdResp = p.productResponse
	return
}

type ProductResponse struct {
	r                 *http.Request
	ProductViewModels []*ProductViewModel `json:"result"`
}

func newProductResponse(r *http.Request) (p *ProductResponse) {
	p = new(ProductResponse)
	p.r = r
	p.ProductViewModels = []*ProductViewModel{}
	return
}
func (p *ProductResponse) addViewModel(viewModel *ProductViewModel) (res *ProductResponse) {
	res = p
	p.ProductViewModels = append(p.ProductViewModels, viewModel)
	return
}

func (p *ProductResponse) addViewModels(viewModels []*ProductViewModel) (res *ProductResponse) {
	res = p
	for _, element := range viewModels {
		p.addViewModel(element)
	}
	return
}

func (p *ProductResponse) show(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(p)
	if err == nil {
		fmt.Fprintf(w, "%s", bytes)
	} else {
		NewStatus(w, "noid", StatusRequestUnsuccessfully, "Can't give you UPC information.").show(appengine.NewContext(p.r))
	}
}

type ProductViewModel struct {
	Status      int     `json:"status"`
	Product     string  `json:"product"`
	Description string  `json:"description"`
	Barcode     string  `json:"barcodeSource"`
	Company     Company `json:"company"`
	People      string  `json:"people"`
	Source      string  `json:"source"`
}

func newProductViewModel(result IProductResult, source string) (item *ProductViewModel) {
	item = new(ProductViewModel)
	item.Source = source
	item.Status = result.getStatus()
	item.Product = result.getProduct()
	item.Description = result.getDescription()
	item.People = result.getPeople()
	item.Barcode = result.getBarcodeUrl()
	item.Company = result.getCompany()
	return
}
