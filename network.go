package qiaoqiao

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
)

func get(r *http.Request, url string, response chan []byte) {
	getMETHOD(r, url, response, nil)
}

func getWithHeader(r *http.Request, url string, response chan []byte, header *http.Header) {
	getMETHOD(r, url, response, header)
}

func getMETHOD(r *http.Request, url string, response chan []byte, header *http.Header) {
	cxt := appengine.NewContext(r)

	log.Infof(cxt, fmt.Sprintf("get url %s", url))
	if req, err := http.NewRequest("GET", url, nil); err == nil {
		httpClient := urlfetch.Client(cxt)
		if header != nil {
			req.Header = *header
		}
		r, err := httpClient.Do(req)
		if r != nil {
			defer r.Body.Close()
		}
		if err == nil {
			if bytes, err := ioutil.ReadAll(r.Body); err == nil {
				wikiRes := new(WikiResult)
				json.Unmarshal(bytes, wikiRes)

				for k, _ := range wikiRes.Query.Pages {
					if k == "-1" {
						response <- nil
						return
					}
				}
				response <- bytes
			} else {
				response <- nil
			}
		} else {
			response <- nil
		}
	} else {
		response <- nil
	}
}
