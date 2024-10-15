package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

type Parametros struct {
	Identificacion   string `json:"identificacion"`
	Facultad         string `json:"facultad"`
	CodigoProyecto   string `json:"codigo_proyecto"`
	CodigoEspacio    string `json:"codigo_espacio"`
	CodigoEstudiante string `json:"codigo_estudiante"`
}

type RequestPayload struct {
	Parametros Parametros `json:"parametros"`
}

func ConsultarCarga(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}
	response := consultarDatos(dataSource, "api_carga_academica", "carga_academica")

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func ConsultarEspacios(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}
	response := consultarDatos(dataSource, "api_espacio_curso", "espacios_academicos")

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func consultarDatos(requestPayload map[string]interface{}, api string, proc string) []map[string]interface{} {
	tokenRequest := helpers.LoginPayload{
		Username: beego.AppConfig.String("UsernameOdin"),
		Password: beego.AppConfig.String("PasswordOdin"),
		Version:  beego.AppConfig.String("VersionOdin"),
	}
	odinService := beego.AppConfig.String("OdinService")

	auth := helpers.GetToken(tokenRequest, "https://"+odinService+"odin/auth/login")
	if auth == nil {
		return nil
	}

	var response []map[string]interface{}
	request.SetHeader("Bearer " + auth.Token)
	err := request.SendJson("https://"+odinService+"odin/gen/apis?api="+api+"&proc="+proc, "POST", &response, requestPayload)

	if err != nil {
		return nil
	}

	return response

}
