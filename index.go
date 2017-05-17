package qiaoqiao

import (
	"net/http"
	"google.golang.org/appengine"
	"fmt"
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
	param := NewParameter(r)
	NewStatus(w, "noid", StatusRequestSuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
}

func handleImagesWikipedia(w http.ResponseWriter, r *http.Request) {
	param := NewParameter(r)
	NewStatus(w, "noid", StatusRequestSuccessfully, fmt.Sprintf("language: %s, keyword: %s", param.Language, param.Keyword)).Succ(appengine.NewContext(r))
}
