package qiaoqiao

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

type WikipediaHandler func(w http.ResponseWriter, r *http.Request, res []byte)

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
				handler(w, r, res)
			} else {
				chBytes = make(chan []byte)
				go wiki.getDoc("en", param.Keyword, chBytes)
				res = <-chBytes
				if res != nil {
					handler(w, r, res)
				} else {
					NewStatus(w, "noid", StatusRequestUnsuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
				}
			}
		} else {
			chBytes = make(chan []byte)
			go wiki.getDoc("en", param.Keyword, chBytes)
			res = <-chBytes
			if res != nil {
				handler(w, r, res)
			} else {
				NewStatus(w, "noid", StatusRequestUnsuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
			}
		}
	} else {
		handler(w, r, res)
	}
}

func outputWikipediaDocument(w http.ResponseWriter, r *http.Request, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", res)
}

func outputWikipediaImage(w http.ResponseWriter, r *http.Request, res []byte) {
	w.Header().Set("Content-Type", "image/*")
	repo := new(WikiResult)
	json.Unmarshal(res, repo)
	for _, v := range repo.Query.Pages {
		chBytes := make(chan []byte)
		go v.Original.get(r, chBytes)
		bys := <-chBytes
		if bys != nil {
			fmt.Fprintf(w, "%s", bys)
		} else {
			chBytes = make(chan []byte)
			go v.Thumbnail.get(r, chBytes)
			bys = <-chBytes
			if bys != nil {
				fmt.Fprintf(w, "%s", bys)
			}
		}
		return
	}
}
