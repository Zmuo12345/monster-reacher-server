package authentication

import (
	"context"
	"testing"
)

var accessToken = ""

func TestSignUp(t *testing.T) {
	server := NewAuthenticationServer()

	res, err := server.SignUp(context.Background(), &SignUpRequest{
		Id:       "test",
		Ip:       "test",
		Platform: "test",
	})

	if err != nil {
		t.Error(err)
	}

	accessToken = res.GetAccessToken()
}

func TestSignIn(t *testing.T) {
	server := NewAuthenticationServer()

	res, err := server.SignIn(context.Background(), &SignInRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		t.Error(err)
	}

	if !res.GetIsValid() {
		t.Error("access token must valid")
	}
}

func TestSignOut(t *testing.T) {
	server := NewAuthenticationServer()

	_, err := server.SignOut(context.Background(), &SignOutRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		t.Error(err)
	}

	res, err := server.SignIn(context.Background(), &SignInRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		t.Error(err)
	}

	if res.GetIsValid() {
		t.Error("access token must not valid")
	}
}
