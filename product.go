package qiaoqiao

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"net/http"
)

func showJson(w http.ResponseWriter, p IRequest) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(p)
	if err == nil {
		fmt.Fprintf(w, "%s", bytes)
	} else {
		NewStatus(w, "noid", StatusRequestUnsuccessfully, "Can't give you UPC information.").show(appengine.NewContext(p.request()))
	}
}

type IRequest interface {
	request() *http.Request
}

type IParsable interface {
	parse(productQuery *ProductQuery) IProductResult
}

type ProductQuery struct {
	r               *http.Request
	params          *Parameter
	targetUrl       string
	key             string
	name            string
	productResponse *ProductResponse
}

type ProductImage struct {
	Large     []string `json:"large"`
	Medium    []string `json:"medium"`
	Small     []string `json:"small"`
	Thumbnail string   `json:"thumbnail"`
	Brand     string   `json:"brand"`
}

type IProductResult interface {
	IParsable
	getStatus() int
	getProduct() string
	getDescription() string
	getPeople() string
	getBarcodeUrl() string
	getCompany() Company
	getProductImage() []ProductImage
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

func (p *ProductQuery) search(productResult IProductResult) *ProductResponse {
	productViewModel := newProductViewModel(productResult.parse(p), p.name)
	if productViewModel.Status == StatusRequestSuccessfully {
		p.productResponse.ProductViewModels = append(p.productResponse.ProductViewModels, productViewModel)
	}
	return p.productResponse
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

func (p *ProductResponse) request() (r *http.Request) {
	r = p.r
	return
}

func (p *ProductResponse) show(w http.ResponseWriter) {
	showJson(w, p)
}

func (p *ProductResponse) toDetail() (ret *ProductViewModel) {
	ret = new(ProductViewModel)
	ret.r = p.r
	for _, v := range p.ProductViewModels {
		ret.Status = v.Status
		if len(ret.Product) < len(v.Product) {
			ret.Product = v.Product
		}
		if len(ret.Description) < len(v.Description) {
			ret.Description = v.Description
		}
		if len(ret.Barcode) < len(v.Barcode) {
			ret.Barcode = v.Barcode
		}
		if len(ret.People) < len(v.People) {
			ret.People = v.People
		}
		ret.Source = ret.Source + " " + v.Source

		if len(ret.ProductImageList) < len(v.ProductImageList) {
			ret.ProductImageList = v.ProductImageList
		}
		if len(ret.Company.Name) < len(v.Company.Name) {
			ret.Company.Name = v.Company.Name
		}
		if len(ret.Company.Logo) < len(v.Company.Logo) {
			ret.Company.Logo = v.Company.Logo
		}
	}
	return
}

type ProductViewModel struct {
	r                *http.Request
	Status           int            `json:"status"`
	Product          string         `json:"product"`
	Description      string         `json:"description"`
	Barcode          string         `json:"barcodeSource"`
	Company          Company        `json:"company"`
	People           string         `json:"people"`
	Source           string         `json:"source"`
	ProductImageList []ProductImage `json:"product_image_list"`
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
	item.ProductImageList = result.getProductImage()
	return
}
func (p *ProductViewModel) request() (r *http.Request) {
	r = p.r
	return
}

func (p *ProductViewModel) show(w http.ResponseWriter) {
	showJson(w, p)
}
