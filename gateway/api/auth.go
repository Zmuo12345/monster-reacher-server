package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"wartech-studio.com/monster-reacher/gateway/services/authentication"
	"wartech-studio.com/monster-reacher/gateway/services/profile"
)

type authApiHandle struct{}

func RegisterAuthApiHandle(router *mux.Router) *authApiHandle {
	handler := &authApiHandle{}
	router.HandleFunc("/api/auth", handler.home)
	router.HandleFunc("/api/auth/signup", handler.register)
	return handler
}

func (*authApiHandle) home(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("This home of auth"))
}

func (*authApiHandle) register(res http.ResponseWriter, req *http.Request) {

	if strings.ToLower(req.Method) != "post" {
		res.Write([]byte("use POST for register by user,password,email or service_name,service_id,service_token"))
		return
	}

	if !checkPostParams(req, []string{"user", "password", "email"}) {
		res.Write([]byte(`{"success": false ,"message":"please check params user,password,email"}`))
		return
	}

	serivces, ok := ServicesDiscoveryCache.CheckRequireServices([]string{"authentication", "profile"})

	if !ok {
		res.Write([]byte(`{"success": false ,"message":"serivces is offline"}`))
		return
	}

	params := mux.Vars(req)

	// check signup from

	ccProfile, err := grpc.Dial(serivces["profile"].GetHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"serivces is error %s"}`, err.Error())))
		return
	}
	defer ccProfile.Close()

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	cProfile := profile.NewProfileClient(ccProfile)

	result, err := cProfile.UserIsValid(ctx, &profile.UserIsValidRequest{User: params["user"]})

	if err != nil {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"serivces is error %s"}`, err.Error())))
		return
	}

	if !result.GetSuccess() {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"user %s is exist"}`, params["user"])))
		return
	}

	resultRegister, err := cProfile.Register(ctx, &profile.RegisterRequest{
		User:     params["user"],
		Password: params["password"],
		Email:    params["email"],
	})

	if err != nil {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"serivces is error %s"}`, err.Error())))
		return
	}

	if resultRegister.GetId() == "" {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"user %s register fail"}`, params["user"])))
		return
	}

	ccAuth, err := grpc.Dial(serivces["authentication"].GetHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"serivces is error %s"}`, err.Error())))
		return
	}
	defer ccAuth.Close()

	cAuth := authentication.NewAuthenticationClient(ccAuth)

	resSignUp, err := cAuth.SignUp(ctx, &authentication.SignUpRequest{Id: resultRegister.GetId()})
	if err != nil {
		res.Write([]byte(fmt.Sprintf(`{"success": false ,"message":"serivces is error %s"}`, err.Error())))
		return
	}
	res.Write([]byte(fmt.Sprintf(`{"success": true ,"access_token":%s , "id":%s}`, resSignUp.GetAccessToken(), resultRegister.GetId())))
}
