package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
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
	response := consultarDatos(data, "api_carga_academica", "carga_academica")

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func ConsultarEspacios(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	response := consultarDatos(data, "api_espacio_curso", "espacios_academicos")

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func consultarDatos(requestPayload []byte, api string, proc string) []map[string]interface{} {
	auth := helpers.GetToken()
	req, err := http.NewRequest("POST", "https://"+beego.AppConfig.String("OdinService")+"odin/gen/apis?api="+api+"&proc="+proc, bytes.NewBuffer(requestPayload))
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.Token)

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
