package qiaoqiao

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func handleProduct(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	param := NewParameter(r)
	p := newProductUpcResponse(r)

	chBytes := make(chan []byte)
	eandataUpc := newProductUpc(r, eandataUrl)
	go eandataUpc.get(param.Language, param.Keyword, chBytes, "eandata")
	byteArray := <-chBytes
	eandata := new(EANdataResult)
	json.Unmarshal(byteArray, eandata)
	log.Infof(cxt, fmt.Sprintf("eandata response %s", string(byteArray)))
	obj := newProductUpcItem(eandata, "eandata")
	if obj.Status == StatusRequestSuccessfully {
		p.ProductUpcItem = append(p.ProductUpcItem, obj)
	}

	for i := 0; i < len(AWS_ASSOCIATE_LIST); i++ {
		var api AmazonProductAPI
		api.AccessKey = AWS_ACCESS_ID
		api.SecretKey = AWS_SECURITY_KEY
		api.AssociateTag = AWS_ASSOCIATE_LIST[i].Tag
		api.Host = AWS_ASSOCIATE_LIST[i].Host
		api.Client = urlfetch.Client(cxt)
		params := map[string]string{
			"ItemId":        param.Keyword,
			"IdType":        "EAN",
			"SearchIndex":   "All",
			"ResponseGroup": awsSerachResponseGroup,
		}
		result, err := api.ItemLookupWithParams(params)
		if err == nil {
			aws := new(ItemLookupResponse)
			xml.Unmarshal([]byte(result), aws)
			log.Infof(cxt, fmt.Sprintf("aws response %s", result))
			obj = newProductUpcItem(aws, "aws")
			p.ProductUpcItem = append(p.ProductUpcItem, obj)
		}
	}
	p.show(w)
}

//type ProductHandler func(w http.ResponseWriter, r *http.Request, res []byte)
//
//func handleProductUniversalProductCode(w http.ResponseWriter, r *http.Request, targetUrl string, handler ProductHandler, service string) {
//	param := NewParameter(r)
//
//	prodUpc := newProductUpc(r, targetUrl)
//	chBytes := make(chan []byte)
//
//	go prodUpc.get(param.Language, param.Keyword, chBytes, service)
//	res := <-chBytes
//
//	handler(w, r, res)
//}

//func handleEANdata(w http.ResponseWriter, r *http.Request, res []byte) {
//	eandata := new(EANdataResult)
//	err := json.Unmarshal(res, eandata)
//	handleData(err, r, newProductUpcItem(eandata, "eandata"), w)
//}
//
//func handleData(err error, r *http.Request, data *ProductUpcItem, w http.ResponseWriter) {
//	if err == nil {
//		p := newProductUpcResponse(r)
//		p.ProductUpcItem = append(p.ProductUpcItem, data)
//		p.show(w)
//	} else {
//		NewStatus(w, "noid", StatusRequestUnsuccessfully, "Handle on data is fail.").show(appengine.NewContext(r))
//	}
//}
