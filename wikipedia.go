package qiaoqiao

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"strings"
	"google.golang.org/appengine/log"
)

type Wikipedia struct {
	r         *http.Request
	targetUrl string
}

func newWikipedia(r *http.Request, targetUrl string) (p *Wikipedia) {
	p = new(Wikipedia)
	p.r = r
	p.targetUrl = targetUrl
	return
}

func (p *Wikipedia) getDoc(language string, keyword string, response chan []byte) {
	url := fmt.Sprintf(baseWikiUrl, language, p.targetUrl, strings.Replace(keyword, " ", "_", -1))
	get(p.r, url, response)
}

func (p *Wikipedia) getGeosearchList(language string, keyword string, response chan []byte) {
	url := fmt.Sprintf(baseWikiUrl, language, p.targetUrl, keyword)
	get(p.r, url, response)
}

type WikiResult struct {
	Query Query `json:"query"`
}

type Query struct {
	Pages map[string]Page `json:"pages"`
}

type Page struct {
	Thumbnail Image `json:"thumbnail"`
	Original  Image `json:"original"`
}

type Image struct {
	Source string `json:"source"`
}

func (p *Image) get(r *http.Request, response chan []byte) {
	get(r, p.Source, response)
}

func get(r *http.Request, url string, response chan []byte) {
	cxt := appengine.NewContext(r)

	log.Infof(cxt, fmt.Sprintf("get url %s", url))
	if req, err := http.NewRequest("GET", url, nil); err == nil {
		httpClient := urlfetch.Client(cxt)
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
