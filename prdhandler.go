package qiaoqiao

import (
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"encoding/xml"
)

type ProductHandler func(w http.ResponseWriter, r *http.Request, res []byte)

func handleProductUniversalProductCode(w http.ResponseWriter, r *http.Request, targetUrl string, handler ProductHandler, service string) {
	param := NewParameter(r)

	prodUpc := newProductUpc(r, targetUrl)
	chBytes := make(chan []byte)

	go prodUpc.get(param.Language, param.Keyword, chBytes, service)
	res := <-chBytes

	handler(w, r, res)
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	param := NewParameter(r)
	p := newProductUpcResponse(r)

	chBytes := make(chan []byte, 1+len(AWS_ASSOCIATE_LIST))
	eandataUpc := newProductUpc(r, eandataUrl)
	go eandataUpc.get(param.Language, param.Keyword, chBytes, "eandata")
	eandata := new(EANdataResult)
	json.Unmarshal(<-chBytes, eandata)
	obj := newProductUpcItem(eandata, "eandata")
	if obj.Status == StatusRequestSuccessfully {
		p.ProductUpcItem = append(p.ProductUpcItem, obj)
	}

	awsUpc := newProductUpc(r, awsUrl)
	go awsUpc.get(param.Language, param.Keyword, chBytes, "aws")
	awsLookup := new(ItemLookupResponse)
	xml.Unmarshal(<-chBytes, awsLookup)
	obj = newProductUpcItem(awsLookup, "aws")
	if obj.Status == StatusRequestSuccessfully {
		p.ProductUpcItem = append(p.ProductUpcItem, obj)
	}
	awsLookup = new(ItemLookupResponse)
	xml.Unmarshal(<-chBytes, awsLookup)
	obj = newProductUpcItem(awsLookup, "aws")
	if obj.Status == StatusRequestSuccessfully {
		p.ProductUpcItem = append(p.ProductUpcItem, obj)
	}

	p.show(w)
}

func handleEANdata(w http.ResponseWriter, r *http.Request, res []byte) {
	eandata := new(EANdataResult)
	err := json.Unmarshal(res, eandata)
	handleData(err, r, newProductUpcItem(eandata, "eandata"), w)
}

func handleAWS(w http.ResponseWriter, r *http.Request, res []byte) {
	awsLookup := new(ItemLookupResponse)
	err := xml.Unmarshal(res, awsLookup)
	handleData(err, r, newProductUpcItem(awsLookup, "aws"), w)
}

func handleData(err error, r *http.Request, data *ProductUpcItem, w http.ResponseWriter) {
	if err == nil {
		p := newProductUpcResponse(r)
		p.ProductUpcItem = append(p.ProductUpcItem, data)
		p.show(w)
	} else {
		NewStatus(w, "noid", StatusRequestUnsuccessfully, "Handle on data is fail.").show(appengine.NewContext(r))
	}
}
