package qiaoqiao

import (
	"fmt"
	"net/http"
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
