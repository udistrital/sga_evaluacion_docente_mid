package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"

	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GuardarRespuestas(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var respuestas []map[string]interface{}
	var nuevaRespuesta map[string]interface{}
	var formularioId int

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Error al parsear el JSON")
		return APIResponseDTO
	}

	formularioId, err := VerificarOCrearFormulario(dataSource)
	if err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al verificar o crear el formulario")
		return APIResponseDTO
	}

	fmt.Println(formularioId)

	if r, ok := dataSource["respuestas"].([]interface{}); ok {
		for _, resp := range r {
			if item, ok := resp.(map[string]interface{}); ok {
				metadata := make(map[string]interface{})

				if itemID, ok := item["item_id"]; ok {
					metadata["item_id"] = itemID
				}
				if valor, ok := item["valor"]; ok {
					metadata["valor"] = valor
				}
				if archivos, ok := item["archivos"]; ok {
					metadata["archivos"] = archivos
				}

				nuevaRespuesta = map[string]interface{}{
					"Activo":            true,
					"FechaCreacion":     time.Now(),
					"FechaModificacion": time.Now(),
					"Id":                0,
					"Metadata":          metadata,
				}

				// Guardar la respuesta
				respuestas = append(respuestas, nuevaRespuesta)
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "No se encontraron respuestas válidas")
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuestas, nil)
	return APIResponseDTO
}

func VerificarOCrearFormulario(dataSource map[string]interface{}) (int, error) {
	var formulario []map[string]interface{}
	var nuevoFormulario map[string]interface{}

	url := fmt.Sprintf(
		"http://%s/formulario?query=PeriodoId:%v,TerceroId:%v,EvaluadoId:%v,EspacioAcademicoId:%v,ProyectoCurricularId:%v&Activo:true&limit=1",
		beego.AppConfig.String("EvaluacionDocenteService"),
		dataSource["id_periodo"],
		dataSource["id_tercero"],
		dataSource["id_evaluado"],
		dataSource["espacio_academico"],
		dataSource["proyecto_curricular"],
	)

	errFormulario := request.GetJson(url, &formulario)
	if errFormulario != nil || len(formulario) == 0 {

		nuevoFormulario = map[string]interface{}{
			"Activo":               true,
			"EspacioAcademicoId":   dataSource["espacio_academico"],
			"EvaluadoId":           dataSource["id_evaluado"],
			"FechaCreacion":        time.Now(),
			"FechaModificacion":    time.Now(),
			"PeriodoId":            dataSource["id_periodo"],
			"ProyectoCurricularId": dataSource["proyecto_curricular"],
			"TerceroId":            dataSource["id_tercero"],
		}

		var createdFormulario map[string]interface{}
		errCrearFormulario := request.SendJson("POST", "http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formulario", nuevoFormulario, &createdFormulario)
		if errCrearFormulario != nil {
			return 0, errors.New("Error al crear el formulario")
		}

		if id, ok := createdFormulario["Id"].(int); ok {
			return id, nil
		}

		return 0, errors.New("No se pudo obtener el ID del formulario creado")
	}

	if id, ok := formulario[0]["Id"].(int); ok {
		return id, nil
	}

	return 0, errors.New("No se pudo obtener el ID del formulario existente")
}
