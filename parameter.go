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
	if r.Method == "POST" {
		p = new(Parameter)
		if bytes, e := ioutil.ReadAll(r.Body); e == nil {
			if e := json.Unmarshal(bytes, p); e == nil {
				return
			}
		}
		p = nil
		return
	}

	if r.Method == "GET" {
		p = new(Parameter)
		args := r.URL.Query()
		p.Language = args["language"][0]
		p.Keyword = args["keyword"][0]
		return
	}
	return
}
