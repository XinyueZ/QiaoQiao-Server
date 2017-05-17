package qiaoqiao

import (
	"net/http"
	//For deployment
	//"../in-messager/inmsgr" //For development
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", handleRoot)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	NewStatus(w, "noid", StatusRequestSuccessfully, "This is helloworld from server.").Succ(appengine.NewContext(r))
}
