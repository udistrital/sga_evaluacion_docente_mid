package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego"
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

func GetToken() (a *Auth) {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil || instance.isTokenExpired() {
		instance = createToken()
	}
	return instance
}

func createToken() (a *Auth) {
	tokenRequest := LoginPayload{
		Username: beego.AppConfig.String("UsernameOdin"),
		Password: beego.AppConfig.String("PasswordOdin"),
		Version:  beego.AppConfig.String("VersionOdin"),
	}

	jsonData, err := json.Marshal(tokenRequest)
	if err != nil {
		fmt.Println("Error al convertir datos a JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "https://"+beego.AppConfig.String("OdinService")+"odin/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error en la solicitud POST:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: estado de la respuesta %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo de la respuesta:", err)
		return
	}

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		fmt.Println("Error al parsear el JSON de la respuesta:", err)
		return
	}

	instance.Token = loginResp.Token
	instance.ExpiryTime = time.Now().Add(tokenTTL)

	return
}
