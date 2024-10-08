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

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	formulario, err := VerificarOCrearFormulario(data)
	if err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al verificar o crear el formulario")
		return APIResponseDTO
	}

	if r, ok := dataSource["respuestas"].([]interface{}); ok {
		for _, resp := range r {
			if item, ok := resp.(map[string]interface{}); ok {
				metadata := make(map[string]interface{})

				if itemID, ok := item["item_id"]; ok {
					metadata["item_id"] = itemID

					if campo, ok := item["campo_id"]; ok { //id del campo hijo
						metadata["campo_id"] = campo
					}
					if valor, ok := item["valor"]; ok {
						metadata["valor"] = valor
					}
					if archivos, ok := item["archivos"]; ok {
						metadata["archivos"] = archivos
					}

					plantilla, err := ObtenerPlantillaPorItemID(itemID)
					if err != nil {
						InactivarFormulario(formulario["Id"].(int))
						APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al obtener la plantilla: %v", err))
						return APIResponseDTO
					}

					metadataJSON, err := json.Marshal(metadata)
					if err != nil {
						InactivarFormulario(formulario["Id"].(int))
						APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al serializar metadata: %v", err))
						return APIResponseDTO
					}

					nuevaRespuesta := map[string]interface{}{
						"Activo":            true,
						"FechaCreacion":     time.Now(),
						"FechaModificacion": time.Now(),
						"Metadata":          string(metadataJSON),
					}

					var nuevaRes map[string]interface{}
					errRespuestas := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/respuesta/", "POST", &nuevaRes, nuevaRespuesta)
					if errRespuestas != nil {
						InactivarFormulario(formulario["Id"].(int))
						APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una de las respuestas")
						return APIResponseDTO
					}
					formularioID := formulario["Id"].(float64)
					plantillaID := plantilla["Id"].(float64)
					respuestaID := nuevaRes["Data"].(map[string]interface{})["Id"].(float64)
					existe := VerificarRespuesta(int(formularioID), int(plantillaID))
					if existe.Status == 200 {
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, formulario, fmt.Sprintf("Ya se han registrado respuestas para este formulario"))
						return APIResponseDTO
					} else {
						fmt.Println("Respuesta no existe")

						relacion := map[string]interface{}{
							"Activo":            true,
							"FechaCreacion":     time.Now(),
							"FechaModificacion": time.Now(),
							"FormularioId":      map[string]interface{}{"Id": int(formularioID)},
							"PlantillaId":       map[string]interface{}{"Id": int(plantillaID)},
							"RespuestaId":       map[string]interface{}{"Id": int(respuestaID)},
						}

						var response map[string]interface{}
						errRelacion := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formrespuesta/", "POST", &response, relacion)
						if errRelacion != nil {
							InactivarFormulario(formulario["Id"].(int))
							APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al crear la relación: %v", errRelacion))
							return APIResponseDTO
						}

						respuestas = append(respuestas, nuevaRes)
					}
				}
			}
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuestas, nil)
		return APIResponseDTO
	}

	InactivarFormulario(formulario["Id"].(int))
	APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "No se encontraron respuestas válidas")
	return APIResponseDTO
}

func VerificarOCrearFormulario(data []byte) (map[string]interface{}, error) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		return nil, fmt.Errorf("error al deserializar los datos: %w", err)
	}

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=Activo:true&sortby=Id&order=asc&limit=0"), &response)
	if errFormulario != nil {
		return nil, fmt.Errorf("error en la petición GET: %w", errFormulario)
	}

	dataList, ok := response["Data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error al convertir los datos a la estructura esperada")
	}

	for _, item := range dataList {
		if formulario, ok := item.(map[string]interface{}); ok {

			if formulario["PeriodoId"] == dataSource["id_periodo"] &&
				formulario["TerceroId"] == dataSource["id_tercero"] &&
				formulario["EvaluadoId"] == dataSource["id_evaluado"] &&
				formulario["ProyectoCurricularId"] == dataSource["proyecto_curricular"] &&
				formulario["EspacioAcademicoId"] == dataSource["espacio_academico"] {

				return formulario, nil
			}
		}
	}

	nuevoFormulario := map[string]interface{}{
		"Activo":               true,
		"EspacioAcademicoId":   dataSource["espacio_academico"],
		"EvaluadoId":           dataSource["id_evaluado"],
		"FechaCreacion":        time.Now(),
		"FechaModificacion":    time.Now(),
		"PeriodoId":            dataSource["id_periodo"],
		"ProyectoCurricularId": dataSource["proyecto_curricular"],
		"TerceroId":            dataSource["id_tercero"],
	}

	errNuevoForm := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formulario/", "POST", &response, nuevoFormulario)
	if errNuevoForm != nil {
		return nil, fmt.Errorf("no se pudo obtener el ID del formulario creado, datos: %v")
	} else {
		// TODO: ajuste temporal
		if response["Success"] == true {
			return response["Data"].(map[string]interface{}), nil
		} else {
			var resp map[string]interface{}
			fmt.Println("goes by here")
			errCheck := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formulario?sortby=Id&order=desc&limit=1&fields=Id", &resp)
			if errCheck == nil && fmt.Sprintf("%v", resp["Data"]) != "[map[]]" {
				nuevoFormulario["TerceroId"] = resp["Data"].([]interface{})[0].(map[string]interface{})["Id"]
				errNuevoForm = request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/formulario/", "POST", &response, nuevoFormulario)
				if errNuevoForm != nil {
					return nil, fmt.Errorf("no se pudo obtener el ID del formulario creado, datos: %v")
				} else {
					return response["Data"].(map[string]interface{}), nil
				}
			} else {
				return nil, fmt.Errorf("no se pudo obtener el ID del formulario creado, datos: %v")
			}
		}
		// END TODO
	}

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

func ObtenerPlantillaPorItemID(itemID interface{}) (map[string]interface{}, error) {
	url := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "/plantilla/?limit=0"
	var response map[string]interface{}
	err := request.SendJson(url, "GET", &response, nil)
	if err != nil {
		return nil, err
	}

	if data, ok := response["Data"].([]interface{}); ok {
		for _, item := range data {
			if plantilla, ok := item.(map[string]interface{}); ok {
				if itemIDFromPlantilla, ok := plantilla["ItemId"].(map[string]interface{}); ok {
					if itemIDFromPlantilla["Id"] == itemID {
						return plantilla, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("Plantilla no encontrada para el item_id: %v", itemID)
}

func VerificarRespuesta(formularioID int, plantillaID int) (APIResponseDTO requestresponse.APIResponse) {
	url := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "formrespuesta?query=Activo:true,FormularioId.Id:" + fmt.Sprint(formularioID) + ",PlantillaId.Id:" + fmt.Sprint(plantillaID)
	var response map[string]interface{}
	err := request.GetJson(url, &response)
	if err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al verificar la respuesta: %v", err))
		return APIResponseDTO
	}

	if data, ok := response["Data"].([]interface{}); ok {
		if len(data) == 0 || (len(data) == 1 && len(data[0].(map[string]interface{})) == 0) {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No se encontraron respuestas previas.")
			return APIResponseDTO
		}
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil, "Respuestas previas encontradas.")
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No se encontraron respuestas previas.")
	return APIResponseDTO
}
