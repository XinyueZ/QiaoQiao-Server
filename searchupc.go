package qiaoqiao

import "net/http"

type SearchUpc struct {
	r         *http.Request
	targetUrl string
}
