package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func MetricasHeteroevaluacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	tipoReporte, ok := dataSource["tipo_reporte"].(string)
	if !ok || (tipoReporte != "global" && tipoReporte != "facultad") {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Tipo de reporte no válido o no especificado")
		return APIResponseDTO
	}

	if tipoReporte == "global" {
		return ReporteGlobal(data)
	} else if tipoReporte == "facultad" {
		return ReporteFacultad(data)
	}

	return requestresponse.APIResponseDTO(false, 500, nil, "Error inesperado al procesar el reporte")
}

func ReporteGlobal(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var campoIds []string
	var docenteIds []string
	var itemIds []string
	var plantillaIds []string
	var respuestasIds []string

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}
	fmt.Println(campoIds)
	fmt.Println(docenteIds)

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["periodo_id"]), &response)
	if errFormulario == nil {
		fmt.Println("--------------------")
		if dataSource["campos"].(map[string]interface{})["componente"] != nil {
			fmt.Println("true")

			var camposResponse map[string]interface{}
			errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=TipoCampoId:2&sortby=Id&order=asc&limit=0&Activo=true"), &camposResponse)
			if errCampos == nil {
				fmt.Println(camposResponse)

				if camposResponse["Data"] != nil {
					for _, campo := range camposResponse["Data"].([]interface{}) {
						campoId := fmt.Sprintf("%v", campo.(map[string]interface{})["Id"])
						campoIds = append(campoIds, campoId)

						var itemCampoResponse map[string]interface{}
						errItemCampo := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=CampoId:%s&Activo=true", campoId), &itemCampoResponse)

						if errItemCampo == nil {
							if itemCampoResponse["Data"] != nil {
								for _, itemCampo := range itemCampoResponse["Data"].([]interface{}) {

									if itemCampo != nil && itemCampo.(map[string]interface{})["ItemId"] != nil {
										itemObj := itemCampo.(map[string]interface{})["ItemId"].(map[string]interface{})

										if itemObj["Id"] != nil {
											itemId := fmt.Sprintf("%v", itemObj["Id"])
											itemIds = append(itemIds, itemId)
										}
									}
								}
							}
						}
					}

					fmt.Println("ItemIds:", itemIds)
					var plantillaResponse map[string]interface{}
					errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?sortby=Id&order=asc&limit=0"), &plantillaResponse)
					if errPlantilla == nil {
						if plantillaResponse["Data"] != nil {
							for _, plantilla := range plantillaResponse["Data"].([]interface{}) {
								itemPlantilla := plantilla.(map[string]interface{})["ItemId"].(map[string]interface{})
								itemId := fmt.Sprintf("%v", itemPlantilla["Id"])

								for _, id := range itemIds {
									if itemId == id {
										plantillaId := fmt.Sprintf("%v", plantilla.(map[string]interface{})["Id"])
										plantillaIds = append(plantillaIds, plantillaId)
										break
									}
								}
							}
							fmt.Println("PlantillaIds:", plantillaIds)

							var formrespuestaResponse map[string]interface{}
							errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?sortby=Id&order=asc&limit=0"), &formrespuestaResponse)
							if errFormrespuesta == nil {
								if formrespuestaResponse["Data"] != nil {
									for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
										respuestaMap := respuesta.(map[string]interface{})
										plantillaRespId := fmt.Sprintf("%v", respuestaMap["PlantillaId"].(map[string]interface{})["Id"])

										for _, plantillaId := range plantillaIds {
											if plantillaRespId == plantillaId {
												respuestaId := fmt.Sprintf("%v", respuestaMap["RespuestaId"].(map[string]interface{})["Id"])
												respuestasIds = append(respuestasIds, respuestaId)
											}
										}
									}
									fmt.Println("Respuestas coincidentes:", respuestasIds)
								}
							}
						}
					}
				}
			}
		}

		if dataSource["campos"].(map[string]interface{})["vinculacion"] != nil {
			fmt.Println("true")

			var resVinculacion map[string]interface{}
			errVinculacion := request.GetJson("http://"+beego.AppConfig.String("PlanDocenteService")+fmt.Sprintf("plan_docente?query=tipo_vinculacion_id:293&sortby=Id&order=asc&limit=0"), &resVinculacion)
			if errVinculacion == nil {
				fmt.Println(resVinculacion["Data"].([]interface{}))
				for _, item := range resVinculacion["Data"].([]interface{}) {
					docenteId := fmt.Sprintf("%v", item.(map[string]interface{})["docente_id"])
					docenteIds = append(docenteIds, docenteId)
				}
			}

			var formVinc []interface{}
			var formIds []string
			for _, formulario := range response["Data"].([]interface{}) {
				terceroId := fmt.Sprintf("%v", formulario.(map[string]interface{})["EvaluadoId"])
				for _, docenteId := range docenteIds {
					if terceroId == docenteId {
						formVinc = append(formVinc, formulario)
						formId := fmt.Sprintf("%v", formulario.(map[string]interface{})["Id"])
						formIds = append(formIds, formId)
						break
					}
				}
			}

			fmt.Println("Formularios", formIds)
		}
		if docenteIds != nil {
			fmt.Println("Docentes:", docenteIds)
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, "formVinc", "Reporte global procesado exitosamente")
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 500, dataSource, "Error al consultar el formulario")
	return APIResponseDTO
}

func ReporteFacultad(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte de facultad procesado exitosamente")
	return APIResponseDTO
}

func MetricasAutoevaluacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	tipoReporte, ok := dataSource["tipo_reporte"].(string)
	if !ok || (tipoReporte != "global" && tipoReporte != "facultad") {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Tipo de reporte no válido o no especificado")
		return APIResponseDTO
	}

	if tipoReporte == "estudiante" {
		return ReporteEstudiante(data)
	} else if tipoReporte == "docente" {
		return ReporteDocente(data)
	}

	return requestresponse.APIResponseDTO(false, 500, nil, "Error inesperado al procesar el reporte")
}

func ReporteEstudiante(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte global procesado exitosamente")
	return APIResponseDTO
}

func ReporteDocente(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte de facultad procesado exitosamente")
	return APIResponseDTO
}
func MetricasCoevaluacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte de facultad procesado exitosamente")
	return APIResponseDTO
}
