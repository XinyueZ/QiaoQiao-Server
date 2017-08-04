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
	presenter := qEAN.search()

	//searchupc.com
	qSearchUpc := newProductQuery(r, params, searchupcUrl, SEARCH_UPC_KEY, "searchupc")
	presenter.addViewModels(qSearchUpc.search().ProductViewModels)

	//barcodable.com
	qBarcodable := newProductQuery(r, params, barcodableUrl, "", "barcodable")
	presenter.addViewModels(qBarcodable.search().ProductViewModels)

	//upcitemdb.com
	qUpcitemdb := newProductQuery(r, params, upcitemdbUrl, "", "upcitemdb")
	presenter.addViewModels(qUpcitemdb.search().ProductViewModels)

	//aws
	for i := 0; i < len(AWS_ASSOCIATE_LIST); i++ {
		var api AmazonProductAPI
		api.AccessKey = AWS_ACCESS_ID
		api.SecretKey = AWS_SECURITY_KEY
		api.AssociateTag = AWS_ASSOCIATE_LIST[i].Tag
		api.Host = AWS_ASSOCIATE_LIST[i].Host
		api.Client = urlfetch.Client(cxt)
		params := map[string]string{
			"ItemId":        params.Keyword,
			"IdType":        "EAN",
			"SearchIndex":   "All",
			"ResponseGroup": awsSerachResponseGroup,
		}
		result, err := api.ItemLookupWithParams(params)
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
	presenter.show(w)
}
