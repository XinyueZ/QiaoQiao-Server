package qiaoqiao

import (
	"encoding/xml"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func handleProduct(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	params := NewParameter(r)

	//eandata.com
	qEAN := newProductQuery(r, params, eandataUrl, EANDATE_KEY, "eandata")
	presenter := qEAN.search(new(EANdataResult))

	//searchupc.com
	qSearchUpc := newProductQuery(r, params, searchupcUrl, SEARCH_UPC_KEY, "searchupc")
	presenter.addViewModels(qSearchUpc.search(new(SearchUpcResult).setCode(params.Keyword)).ProductViewModels)

	//barcodable.com
	qBarcodable := newProductQuery(r, params, barcodableUrl, "", "barcodable")
	presenter.addViewModels(qBarcodable.search(new(BarcodableResult).setCodeType("upc")).ProductViewModels)

	//upcitemdb.com
	qUpcitemdb := newProductQuery(r, params, upcitemdbUrl, "", "upcitemdb")
	presenter.addViewModels(qUpcitemdb.search(new(UpcItemDbResult)).ProductViewModels)

	//Walmart
	qWalmart:= newProductQuery(r, params, walmartUrl, WALMART_KEY, "walmart")
	presenter.addViewModels(qWalmart.search(new(WalmartResult)).ProductViewModels)

	//tesco
	qTesco := newProductQuery(r, params, tescoUrl, TESCO_KEY, "tesco")
	presenter.addViewModels(qTesco.search(new(TescoResult)).ProductViewModels)


	//aws
	for i := 0; i < len(AWS_ASSOCIATE_LIST); i++ {
		var api AmazonProductAPI
		api.AccessKey = AWS_ACCESS_ID
		api.SecretKey = AWS_SECURITY_KEY
		api.AssociateTag = AWS_ASSOCIATE_LIST[i].Tag
		api.Host = AWS_ASSOCIATE_LIST[i].Host
		api.Client = urlfetch.Client(cxt)
		awsparams := map[string]string{
			"ItemId":        params.Keyword,
			"IdType":        "EAN",
			"SearchIndex":   "All",
			"ResponseGroup": awsSerachResponseGroup,
		}
		result, err := api.ItemLookupWithParams(awsparams)
		if err == nil {
			aws := new(ItemLookupResponse)
			xml.Unmarshal([]byte(result), aws)
			if aws.getStatus() == StatusRequestSuccessfully {
				log.Infof(cxt, fmt.Sprintf("aws feeds %s", result))
				obj := newProductViewModel(aws, "aws")
				presenter.addViewModel(obj)
			} else {
				awsparams := map[string]string{
					"ItemId":        params.Keyword,
					"IdType":        "UPC",
					"SearchIndex":   "All",
					"ResponseGroup": awsSerachResponseGroup,
				}
				result, err := api.ItemLookupWithParams(awsparams)
				if err == nil {
					aws := new(ItemLookupResponse)
					xml.Unmarshal([]byte(result), aws)
					if aws.getStatus() == StatusRequestSuccessfully {
						log.Infof(cxt, fmt.Sprintf("aws feeds %s", result))
						obj := newProductViewModel(aws, "aws")
						presenter.addViewModel(obj)
					}
				}
			}
		}
	}
	presenter.show(w)
}
