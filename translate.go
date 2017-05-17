package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"net/http"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
)

type Translator struct {
	r *http.Request
}

func newTranslator(r *http.Request) (p *Translator) {
	p = new(Translator)
	p.r = r
	return
}

func (p *Translator) get(q string, target string, format string, response chan *Response) {
	cxt := appengine.NewContext(p.r)
	url := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2/?q=%s&target=%s&format=%s&key=%s", q, target, format, googleApiKey)
	if req, err := http.NewRequest("GET", url, nil); err == nil {
		httpClient := urlfetch.Client(cxt)
		r, err := httpClient.Do(req)
		if r != nil {
			defer r.Body.Close()
		}
		if err == nil {
			if bytes, err := ioutil.ReadAll(r.Body); err == nil {
				resp := new(Response)
				json.Unmarshal(bytes, resp)
				response <- resp
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

type Response struct {
	Data Data `json:"data"`
}

type TranslateTextResponseTranslation struct {
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
	Model                  string `json:"model"`
	TranslatedText         string `json:"translatedText"`
}

type Data struct {
	Translations []TranslateTextResponseTranslation  `json:"translations"`
}
