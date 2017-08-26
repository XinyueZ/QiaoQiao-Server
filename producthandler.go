package qiaoqiao

import (
	"encoding/xml"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func handleProductByUpc(w http.ResponseWriter, r *http.Request) {
	presenter := buildProductResponse(r)
	presenter.show(w)
}

func handleProductDetailByUpc(w http.ResponseWriter, r *http.Request) {
	presenter := buildProductResponse(r)
	presenter.toDetail().show(w)
}

func buildProductResponse(r *http.Request) *ProductResponse {
	cxt := appengine.NewContext(r)
	params := NewParameter(r)

	ch := make(chan []*ProductViewModel, 5) // don't care the first one, eandata.com, 5 for others exclude aws
	sign := make(chan int, 1)

	// eandata.com
	qEAN := newProductQuery(r, params, eandataUrl, EANDATE_KEY, "eandata")
	presenter := qEAN.search(new(EANdataResult))

	//1. searchupc.com
	qSearchUpc := newProductQuery(r, params, searchupcUrl, SEARCH_UPC_KEY, "searchupc")
	ch <- qSearchUpc.search(new(SearchUpcResult).setCode(params.Keyword)).ProductViewModels

	//2. barcodable.com
	qBarcodable := newProductQuery(r, params, barcodableUrl, "", "barcodable")
	ch <- qBarcodable.search(new(BarcodableResult).setCodeType("upc")).ProductViewModels

	//3. upcitemdb.com
	qUpcitemdb := newProductQuery(r, params, upcitemdbUrl, "", "upcitemdb")
	ch <- qUpcitemdb.search(new(UpcItemDbResult)).ProductViewModels

	//4. Walmart
	qWalmart := newProductQuery(r, params, walmartUrl, WALMART_KEY, "walmart")
	ch <- qWalmart.search(new(WalmartResult)).ProductViewModels

	//5. tesco
	qTesco := newProductQuery(r, params, tescoUrl, TESCO_KEY, "tesco")
	ch <- qTesco.search(new(TescoResult)).ProductViewModels

	// aws
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

	close(ch)
	go func() {
		var out []*ProductViewModel
		ok := true
		for {
			select {
			case out, ok = <-ch:
				if !ok {
					break
				} else {
					presenter.addViewModels(out)
				}
			}

			if !ok {
				sign <- 0
				break
			}
		}

	}()
	<-sign
	return presenter
}
