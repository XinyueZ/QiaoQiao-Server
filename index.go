package qiaoqiao

import (
	"net/http"
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/knowledge/documents/wikipedia", handleDocumentsWikipedia)
	http.HandleFunc("/knowledge/images/wikipedia", handleImagesWikipedia)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	NewStatus(w, "noid", StatusRequestSuccessfully, "This is helloworld from server.").Succ(appengine.NewContext(r))
}

func handleDocumentsWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiDocuments, outputWikipediaDocument)
}

func handleImagesWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiImages, outputWikipediaImage)

}
