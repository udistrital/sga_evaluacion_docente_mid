package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/requestresponse"
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

type Parametros struct {
	Rol            string `json:"rol"`
	Identificacion string `json:"identificacion"`
	Facultad       string `json:"facultad"`
	Proyecto       string `json:"proyecto"`
	Anio           int    `json:"anio"`
	Periodo        int    `json:"periodo"`
}

type RequestPayload struct {
	Parametros Parametros `json:"parametros"`
}

const tokenTTL = time.Hour * 24

func ConsultarCarga(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	auth := &Auth{}
	auth.getToken()
	response := auth.consultarCarga(data)

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func (a *Auth) getToken() {

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

	a.Token = loginResp.Token
	a.ExpiryTime = time.Now().Add(tokenTTL)
}

func (a *Auth) isTokenExpired() bool {
	return time.Now().After(a.ExpiryTime)
}

func (a *Auth) consultarCarga(requestPayload []byte) []map[string]interface{} {
	if a.isTokenExpired() {
		fmt.Println("El token ha expirado, obteniendo uno nuevo...")
		a.getToken()
	}

	req, err := http.NewRequest("POST", "https://"+beego.AppConfig.String("OdinService")+"odin/gen/apis?api=api_carga_academica&proc=carga_academica", bytes.NewBuffer(requestPayload))
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error en la solicitud POST:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo de la respuesta:", err)
		return nil
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return nil
	}

	return result
}
