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

func checkPostParams(req *http.Request, params []string) bool {
	ps := mux.Vars(req)
	found := 0
	for i := range params {
		for j := range ps {
			if params[i] == j && ps[j] != "" {
				found++
				break
			}
		}
	}

	return found == len(params)
}
