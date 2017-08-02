package qiaoqiao

import (
	"net/http"
	"encoding/json"
	"google.golang.org/appengine"
)

type ProductHandler func(w http.ResponseWriter, r *http.Request, res []byte)

func handleProductUniversalProductCode(w http.ResponseWriter, r *http.Request, targetUrl string, handler ProductHandler) {
	param := NewParameter(r)
	prodUpc := newProductUpc(r, targetUrl)
	chBytes := make(chan []byte)
	go prodUpc.get(param.Language, param.Keyword, chBytes)
	res := <-chBytes
	handler(w, r, res)
}

func handleEANdata(w http.ResponseWriter, r *http.Request, res []byte) {
	eandata := new(EANdataResult)
	err := json.Unmarshal(res, eandata)
	if err == nil {
		newProductUpcResponse(r, eandata).show(w)
	} else {
		NewStatus(w, "noid", StatusRequestUnsuccessfully, "EANdata call fail.").show(appengine.NewContext(r))
	}
}
