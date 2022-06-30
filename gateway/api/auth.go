package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type authApiHandle struct{}

func RegisterAuthApiHandle(router *mux.Router) *authApiHandle {
	handler := &authApiHandle{}
	router.HandleFunc("/api/auth", handler.home)
	router.HandleFunc("/api/auth/signup", handler.signup)
	return handler
}

func (*authApiHandle) home(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("This home of auth"))
}

func (*authApiHandle) signup(res http.ResponseWriter, req *http.Request) {
	// check signup from
	// validate
	// profile.register
	// profile.getdata
	// authentication.signup

	accessToken := ""
	success := false
	res.Write([]byte(fmt.Sprintf(`{"success": %v ,"access_token":%s}`, success, accessToken)))
}
