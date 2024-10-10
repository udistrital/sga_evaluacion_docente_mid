package helpers

import (
	"sync"
	"time"

	"github.com/udistrital/utils_oas/request"
)

type Auth struct {
	Token      string
	ExpiryTime time.Time
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string `json:"version"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

const tokenTTL = time.Hour * 24

var instance *Auth

var mutex = &sync.Mutex{}

func (a *Auth) isTokenExpired() bool {
	return time.Now().After(a.ExpiryTime)
}

func GetToken(tokenRequest LoginPayload, url string) (a *Auth) {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil || instance.isTokenExpired() {
		instance = createToken(tokenRequest, url)
	}
	return instance
}

func createToken(tokenRequest LoginPayload, url string) (a *Auth) {
	var loginResp LoginResponse
	err := request.SendJson(url, "POST", &loginResp, tokenRequest)

	if err != nil {
		return nil
	}

	instance = &Auth{
		Token:      loginResp.Token,
		ExpiryTime: time.Now().Add(tokenTTL),
	}

	return instance
}
