package qiaoqiao

import (
	"net/http"
)

func init() {
	http.HandleFunc("/knowledge/images/wikipedia", handleImagesWikipedia)
	http.HandleFunc("/knowledge/documents/wikipedia", handleDocumentsWikipedia)
	http.HandleFunc("/knowledge/geosearch/wikipedia", handleGeosearchWikipedia)
}

func handleImagesWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiImages, outputWikipediaImage)

}

func handleDocumentsWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiDocuments, outputWikipediaDocument)
}

func handleGeosearchWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipediaGeosearch(w, r, urlWikiGeosearch, outputWikipediaDocument)
}
