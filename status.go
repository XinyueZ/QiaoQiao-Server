package qiaoqiao

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

const (
	StatusRequestSuccessfully   = 201
	StatusRequestUnsuccessfully = 501
)

//Status gives us request status-information.
type Status struct {
	responseWriter http.ResponseWriter
	ID             string `json:"id"`
	Status         int    `json:"status"`
	Message        string `json:"message"`
}

func NewStatus(w http.ResponseWriter, id string, status int, message string) (r *Status) {
	r = &Status{w, id, status, message}
	return
}


func (p *Status) show(context context.Context) {
	p.setupHeader()
	p.output()
}

func (p *Status) setupHeader() {
	p.responseWriter.Header().Set("Content-Type", "application/json")
}

func (p *Status) output() {
	fmt.Fprintf(p.responseWriter, `{"status":%d, "message": "%s", "id": "%s"}`, p.Status, p.Message, p.ID)
}
