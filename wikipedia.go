package qiaoqiao

import (
	"net/http"
	"google.golang.org/appengine"
	"fmt"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
	"strings"
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
	cxt := appengine.NewContext(p.r)
	url := fmt.Sprintf("https://%s.wikipedia.org/%s%s", language, p.targetUrl, keyword)
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
				if strings.Contains(string(bytes), "{\"pages\":{\"-1\"") {
					response <- nil
					return
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

type WikiResult struct {
	Query Query `json:"query"`
}

type Query struct {
	Pages Pages `json:"pages"`
}

type Page struct {
	Thumbnail Image `json:"thumbnail"`
	Original  Image `json:"original"`
}

type Pages struct {
	Page map[string]Page
}

type Image struct {
	Source string `json:"source"`
}
