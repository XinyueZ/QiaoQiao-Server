package qiaoqiao

import (
	"google.golang.org/appengine"
	"fmt"
	"net/http"
	"google.golang.org/appengine/log"
)

type WikipediaHandler func(w http.ResponseWriter, res []byte)

func handleWikipedia(w http.ResponseWriter, r *http.Request, targetUrl string, handler WikipediaHandler) {
	param := NewParameter(r)
	wiki := newWikipedia(r, targetUrl)
	chBytes := make(chan []byte)
	go wiki.getDoc(param.Language, param.Keyword, chBytes)
	res := <-chBytes
	if res == nil {
		chResponse := make(chan *Response)
		go newTranslator(r).get(param.Keyword, param.Language, "text", chResponse)
		transResponse := <-chResponse
		if transResponse != nil && len(transResponse.Data.Translations) > 0 {
			translation := transResponse.Data.Translations[0]
			translatedText := translation.TranslatedText

			log.Infof(appengine.NewContext(r), fmt.Sprintf("translated q=%s -> %s", param.Keyword, translatedText))
			go wiki.getDoc(param.Language, translatedText, chBytes)
			res = <-chBytes
			if res != nil {
				handler(w, res)
			} else {
				chBytes = make(chan []byte)
				go wiki.getDoc("en", param.Keyword, chBytes)
				res = <-chBytes
				if res != nil {
					handler(w, res)
				} else {
					NewStatus(w, "noid", StatusRequestUnsuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
				}
			}
		} else {
			chBytes = make(chan []byte)
			go wiki.getDoc("en", param.Keyword, chBytes)
			res = <-chBytes
			if res != nil {
				handler(w, res)
			} else {
				NewStatus(w, "noid", StatusRequestUnsuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
			}
		}
	} else {
		handler(w, res)
	}
}

func outputWikipediaDocument(w http.ResponseWriter, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", res)
}
