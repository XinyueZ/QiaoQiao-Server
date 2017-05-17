package qiaoqiao

import (
	"net/http"
	"google.golang.org/appengine"
	"fmt"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/knowledge/documents/wikipedia", handleDocumentsWikipedia)
	http.HandleFunc("/knowledge/images/wikipdia", handleImagesWikipedia)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	NewStatus(w, "noid", StatusRequestSuccessfully, "This is helloworld from server.").Succ(appengine.NewContext(r))
}

func handleDocumentsWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, "w/api.php?format=json&action=query&prop=extracts|pageimages|langlinks&llprop=autonym&lldir=descending&lllimit=500&piprop=original|name|thumbnail&exlimit=1&redirects=titles&titles=")
}

func handleImagesWikipedia(w http.ResponseWriter, r *http.Request) {
	param := NewParameter(r)
	NewStatus(w, "noid", StatusRequestSuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
}

func handleWikipedia(w http.ResponseWriter, r *http.Request, targetUrl string) {
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
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, string(res))
			} else {
				chBytes = make(chan []byte)
				go wiki.getDoc("en", param.Keyword, chBytes)
				res = <-chBytes
				if res != nil {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintf(w, string(res))
				} else {
					NewStatus(w, "noid", StatusRequestUnsuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
				}
			}
		} else {
			chBytes = make(chan []byte)
			go wiki.getDoc("en", param.Keyword, chBytes)
			res = <-chBytes
			if res != nil {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, string(res))
			} else {
				NewStatus(w, "noid", StatusRequestUnsuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
			}
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(res))
	}
}
