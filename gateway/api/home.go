package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type homeApiHandle struct{}

func RegisterHomeApiHandle(router *mux.Router) *homeApiHandle {
	handler := &homeApiHandle{}
	router.HandleFunc("/api", handler.home)
	return handler
}

func (*homeApiHandle) home(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("This home of api"))
}
