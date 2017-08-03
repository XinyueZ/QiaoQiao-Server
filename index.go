package qiaoqiao

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/knowledge/product/upc", handleProductUpc)
	http.HandleFunc("/knowledge/images/daily", handleImageDaily)
	http.HandleFunc("/knowledge/id/wikipedia", handleIdWikipedia)
	http.HandleFunc("/knowledge/images/wikipedia", handleImagesWikipedia)
	http.HandleFunc("/knowledge/documents/wikipedia", handleDocumentsWikipedia)
	http.HandleFunc("/knowledge/geosearch/wikipedia", handleGeosearchWikipedia)
	http.HandleFunc("/knowledge/thumbnails/wikipedia", handleThumbnailsWikipedia)
}

func handleIdWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipediaId(w, r, urlWikiId, outputWikipediaDocument)
}

func handleImagesWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiImages, outputWikipediaImage)
}

func handleThumbnailsWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiThumbnails, outputWikipediaThumbnail)
}

func handleDocumentsWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipedia(w, r, urlWikiDocuments, outputWikipediaDocument)
}

func handleGeosearchWikipedia(w http.ResponseWriter, r *http.Request) {
	handleWikipediaGeosearch(w, r, urlWikiGeosearch, outputWikipediaDocument)
}

func handleImageDaily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"creatives\" : [{\"url\" : \"%s\"}]}", "https://source.unsplash.com/random")
}

func handleProductUpc(w http.ResponseWriter, r *http.Request) {
	handleProductUniversalProductCode(w, r, eandataUrl, handleEANdata, "eandata")
	//handleProductUniversalProductCode(w, r, awsUrl, handleAWS, "aws")
}
