package qiaoqiao

import (
	"net/http"
	"fmt"
)

type ProductUpc struct {
	r         *http.Request
	targetUrl string
}

func newProductUpc(r *http.Request, targetUrl string) (p *ProductUpc) {
	p = new(ProductUpc)
	p.r = r
	p.targetUrl = targetUrl
	return
}


func (p *ProductUpc) get(language string, code string, response chan []byte) {
	url := fmt.Sprintf(p.targetUrl, code, eandateKey)
	get(p.r, url, response)
}

