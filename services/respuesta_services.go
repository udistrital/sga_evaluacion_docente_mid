package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"

	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GuardarRespuestas(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var respuestas []map[string]interface{}
	var nuevaRes map[string]interface{}

	fmt.Println("JSON Recibido:", string(data))

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	formularioId, err := VerificarOCrearFormulario(data)
	if err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al verificar o crear el formulario")
		return APIResponseDTO
	}

	fmt.Println("Formulario ID:", formularioId)

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

				metadataJSON, err := json.Marshal(metadata)
				if err != nil {
					APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al serializar metadata: %v", err))
					return APIResponseDTO
				}

				nuevaRespuesta := map[string]interface{}{
					"Activo":            true,
					"FechaCreacion":     time.Now(),
					"FechaModificacion": time.Now(),
					"Id":                24,
					"Metadata":          string(metadataJSON),
				}
				respuestas = append(respuestas, nuevaRespuesta)

				errRespuestas := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/respuesta/", "POST", &nuevaRes, nuevaRespuesta)
				if errRespuestas == nil {
					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nuevaRes)
					return APIResponseDTO
				} else {
					InactivarRespuesta(nuevaRes["Id"].(int))
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "No se encontraron respuestas válidas")
					return APIResponseDTO
				}
			}
		}
	} else {
		InactivarFormulario(formularioId)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "No se encontraron respuestas válidas")
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuestas, nil)
	return APIResponseDTO
}

func VerificarOCrearFormulario(data []byte) (int, error) {
	var nuevoFormulario map[string]interface{}
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		return 0, fmt.Errorf("error al deserializar los datos: %w", err)
	}

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=Activo:true&sortby=Id&order=asc&limit=0"), &response)
	fmt.Println("Respuesta de la petición GET:", response)
	if errFormulario != nil {
		return 0, fmt.Errorf("error en la petición GET: %w", errFormulario)
	}

	dataList, ok := response["Data"].([]interface{})
	if !ok {
		return 0, fmt.Errorf("formato de respuesta inesperado: %v", response)
	}
	if len(dataList) > 0 {
		if firstFormulario, ok := dataList[0].(map[string]interface{}); ok {
			if id, ok := firstFormulario["Id"].(float64); ok {
				return int(id), nil
			}
			return 0, fmt.Errorf("el formulario existente no tiene un ID válido: %v", firstFormulario)
		}
		return 0, fmt.Errorf("estructura inesperada en los datos del formulario: %v", dataList[0])
	}

	nuevoFormulario = map[string]interface{}{
		"Activo":               true,
		"EspacioAcademicoId":   dataSource["espacio_academico"],
		"EvaluadoId":           dataSource["id_evaluado"],
		"FechaCreacion":        time.Now().Format(time.RFC3339),
		"FechaModificacion":    time.Now().Format(time.RFC3339),
		"PeriodoId":            dataSource["id_periodo"],
		"ProyectoCurricularId": dataSource["proyecto_curricular"],
		"TerceroId":            dataSource["id_tercero"],
	}

	fmt.Println("Datos del nuevo formulario:", nuevoFormulario)

	urlPost := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "/formulario"
	fmt.Println("URL de creación del formulario:", urlPost)

	var createdFormulario map[string]interface{}
	errCrearFormulario := request.SendJson("POST", urlPost, nuevoFormulario, &createdFormulario)
	if errCrearFormulario != nil {
		return 0, fmt.Errorf("error al crear el formulario: %w", errCrearFormulario)
	}

	if id, ok := createdFormulario["Id"].(float64); ok {
		return int(id), nil
	}

	return 0, fmt.Errorf("no se pudo obtener el ID del formulario creado, datos: %v", createdFormulario)
}
func InactivarFormulario(id int) error {
	var formulario map[string]interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formulario/"+fmt.Sprint(id), &formulario)
	if err != nil {
		return fmt.Errorf("error al obtener el formulario con ID %d: %v", id, err)
	}
	formulario["Activo"] = false
	err = request.SendJson("PUT", "http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formulario/"+fmt.Sprint(id), formulario, nil)
	if err != nil {
		return fmt.Errorf("error al inactivar el formulario con ID %d: %v", id, err)
	}

	return nil
}

func InactivarRespuesta(id int) error {
	var respuesta map[string]interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/respuesta/"+fmt.Sprint(id), &respuesta)
	if err != nil {
		return fmt.Errorf("error al obtener la respuesta con ID %d: %v", id, err)
	}

	respuesta["Activo"] = false

	err = request.SendJson("PUT", "http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/respuesta/"+fmt.Sprint(id), respuesta, nil)
	if err != nil {
		return fmt.Errorf("error al inactivar la respuesta con ID %d: %v", id, err)
	}

	return nil
}
