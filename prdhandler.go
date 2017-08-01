package qiaoqiao

import "net/http"

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
	//TODO Output feed of UPC to client
}
