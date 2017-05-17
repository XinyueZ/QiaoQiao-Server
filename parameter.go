package qiaoqiao

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Parameter struct {
	Language string `json:"language"`
	Keyword  string `json:"keyword"`
}

func NewParameter(r *http.Request) (p *Parameter) {
	p = new(Parameter)
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		if e := json.Unmarshal(bytes, p); e == nil {
			return
		}
	}
	p = nil
	return
}

