package qiaoqiao

import (
	"encoding/xml"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func concatAppend(slices [][]byte) []byte {
	var tmp []byte
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	params := NewParameter(r)

	//eandata.com
	eandataUpc := newProductRequest(r, params, eandataUrl, EANDATE_KEY, "eandata")
	presenter := eandataUpc.search()

	//searchupc.com
	searchupc := newProductRequest(r, params, searchupcUrl, SEARCH_UPC_KEY, "searchupc")
	presenter.addViewModels(searchupc.search().ProductViewModels)

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
				log.Infof(cxt, fmt.Sprintf("aws response %s", result))
				obj := newProductViewModel(aws, "aws")
				presenter.addViewModel(obj)
			}
		}
	}
	presenter.show(w)
}
