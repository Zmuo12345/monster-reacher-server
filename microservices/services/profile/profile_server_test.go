package profile

import (
	"context"
	"testing"
)

func TestGetData(t *testing.T) {

}

func TestAuthentication(t *testing.T) {

}

func TestAuthenticationByService(t *testing.T) {

}

func TestRegister(t *testing.T) {
	server := NewProfileServer()
	res, err := server.Register(context.Background(), &RegisterRequest{
		User:     "test",
		Password: "testpass",
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(res.GetSuccess())
}

func TestRegisterByService(t *testing.T) {

}
