package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
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
	response := consultarDatos(dataSource, "espacios_estudiante")

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func ConsultarEspacios(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}
	response := consultarDatos(dataSource, "carga_academica_docente")

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func consultarDatos(requestPayload map[string]interface{}, service string) interface{} {

	var respuesta interface{}

	parametros := requestPayload["parametros"].(map[string]interface{})
	identificacion := parametros["identificacion"].(string)

	url := beego.AppConfig.String("ProtocolAdmin") + "://" +
		beego.AppConfig.String("UrlWSO2") +
		beego.AppConfig.String("NsAcademica") + "/" + service + "/" + identificacion + "/0"

	if err := request.GetJsonWSO2(url, &respuesta); err != nil {
		return nil
	}
	return respuesta
}
