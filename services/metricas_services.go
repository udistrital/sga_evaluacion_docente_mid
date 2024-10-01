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
	var itemIds []string
	var itemIdsCampos []string
	var plantillaIds []string
	var respuestasIds []string
	var formularioIds []string

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["periodo_id"]), &response)

	type ItemPlantillaRespuesta struct {
		ItemId             string                   `json:"item_id"`
		PlantillaId        string                   `json:"plantilla_id"`
		CantidadRespuestas int                      `json:"cantidad_respuestas"`
		RespuestasDetalle  []map[string]interface{} `json:"respuestas_detalle"` // Metadata, aquí puede asignarse el valor, o el UID de los archivos según se requiera
	}

	var itemsPlantillasRespuestas []ItemPlantillaRespuesta

	if errFormulario == nil {
		if dataSource["campos"].(map[string]interface{})["componente"] != nil {
			var camposResponse map[string]interface{}
			errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=TipoCampoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["campos"].(map[string]interface{})["componente"]), &camposResponse)
			if errCampos == nil {
				if camposResponse["Data"] != nil {
					for _, campo := range camposResponse["Data"].([]interface{}) {
						campoId := fmt.Sprintf("%v", campo.(map[string]interface{})["Id"])
						campoIds = append(campoIds, campoId)

						var itemCampoResponse map[string]interface{}
						errItemCampo := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=CampoId:%s&Activo=true&limit=0", campoId), &itemCampoResponse)

						if errItemCampo == nil {
							if itemCampoResponse["Data"] != nil {
								for _, itemCampo := range itemCampoResponse["Data"].([]interface{}) {
									if itemCampo != nil && itemCampo.(map[string]interface{})["ItemId"] != nil {
										itemObj := itemCampo.(map[string]interface{})["ItemId"].(map[string]interface{})
										if itemObj["Id"] != nil {
											itemId := fmt.Sprintf("%v", itemObj["Id"])
											itemIdsCampos = append(itemIdsCampos, itemId)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if dataSource["campos"].(map[string]interface{})["vinculacion"] != nil {
			var vinculacionResponse map[string]interface{}
			errVinculacion := request.GetJson("http://"+beego.AppConfig.String("PlanDocenteService")+fmt.Sprintf("plan_docente?query=tipo_vinculacion_id:%v&sortby=Id&order=asc&limit=0", dataSource["campos"].(map[string]interface{})["vinculacion"]), &vinculacionResponse)

			if errVinculacion == nil {
				if data, ok := vinculacionResponse["Data"].(map[string]interface{}); ok {
					var docenteIds []string

					for _, item := range data {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if docenteId, ok := itemMap["docente_id"].(string); ok {
								docenteIds = append(docenteIds, docenteId)
							}
						}
					}

					if response["Data"] != nil {
						var filteredFormularios []map[string]interface{}
						for _, formulario := range response["Data"].([]interface{}) {
							formularioMap := formulario.(map[string]interface{})
							evaluadoId := fmt.Sprintf("%v", formularioMap["EvaluadoId"])
							if contains(docenteIds, evaluadoId) {
								filteredFormularios = append(filteredFormularios, formularioMap)
							}
						}
						response["Data"] = filteredFormularios
					}
				}
			}
		}

		itemIds = itemIdsCampos

		if response["Data"] != nil {
			for _, formulario := range response["Data"].([]interface{}) {
				formularioMap := formulario.(map[string]interface{})
				formularioId := fmt.Sprintf("%v", formularioMap["Id"])
				formularioIds = append(formularioIds, formularioId)
			}
		}

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

							var formrespuestaResponse map[string]interface{}
							errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?sortby=Id&order=asc&limit=0"), &formrespuestaResponse)
							if errFormrespuesta == nil {
								var respuestasDetalle []map[string]interface{}
								if formrespuestaResponse["Data"] != nil {
									for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
										respuestaMap := respuesta.(map[string]interface{})
										plantillaRespId := fmt.Sprintf("%v", respuestaMap["PlantillaId"].(map[string]interface{})["Id"])

										if plantillaRespId == plantillaId && contains(formularioIds, fmt.Sprintf("%v", respuestaMap["FormularioId"].(map[string]interface{})["Id"])) {
											respuestaId := fmt.Sprintf("%v", respuestaMap["RespuestaId"].(map[string]interface{})["Id"])
											respuestasIds = append(respuestasIds, respuestaId)

											var respuestaDetalleResponse map[string]interface{}
											errRespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("respuesta/%s&order=asc&limit=0", respuestaId), &respuestaDetalleResponse)
											if errRespuesta == nil {
												if respuestaDetalleResponse["Data"] != nil {
													respuestaDetalle := respuestaDetalleResponse["Data"].(map[string]interface{})
													metadataStr := respuestaDetalle["Metadata"]
													if metadataStr != nil {
														var metadataMap map[string]interface{}
														err := json.Unmarshal([]byte(metadataStr.(string)), &metadataMap)
														if err == nil {
															valor, ok := metadataMap["valor"]
															if !ok {
																valor = nil
															}
															campoId, ok := metadataMap["campo_id"]
															if !ok {
																campoId = nil
															}
															respuestasDetalle = append(respuestasDetalle, map[string]interface{}{
																"Metadata": metadataMap,
																"Valor":    valor,
																"CampoId":  campoId,
															})
														}
													}
												}
											}
										}
									}
								}

								itemsPlantillasRespuestas = append(itemsPlantillasRespuestas, ItemPlantillaRespuesta{
									ItemId:             itemId,
									PlantillaId:        plantillaId,
									CantidadRespuestas: len(respuestasIds),
									RespuestasDetalle:  respuestasDetalle,
								})
							}
						}
					}
				}
			}
		}
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, itemsPlantillasRespuestas, "")
	return APIResponseDTO
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ReporteFacultad(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var campoIds []string
	var itemIds []string
	var itemIdsCampos []string
	var plantillaIds []string
	var respuestasIds []string
	var idsProyectos []string
	var formularioIds []string

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}
	var resFacultad map[string]interface{}
	errFacultad := request.GetJson("http://"+beego.AppConfig.String("ProyectoService")+fmt.Sprintf("proyecto-academico?sortby=Id&order=asc&limit=0&Activo=true"), &resFacultad)

	if errFacultad == nil {

		facultadID := dataSource["facultad_id"].(float64)

		if data, ok := resFacultad["Data"].([]interface{}); ok {
			for _, proyecto := range data {
				if proyectoMap, ok := proyecto.(map[string]interface{}); ok {
					if proyectoAcademico, ok := proyectoMap["ProyectoAcademico"].(map[string]interface{}); ok {
						if id, ok := proyectoAcademico["FacultadId"].(float64); ok {

							if id == facultadID {
								idsProyectos = append(idsProyectos, fmt.Sprintf("%.0f", proyectoAcademico["Id"].(float64)))
							}
						}
					}
				}
			}
		}
	}

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["periodo_id"]), &response)

	type ItemPlantillaRespuesta struct {
		ItemId             string                   `json:"item_id"`
		PlantillaId        string                   `json:"plantilla_id"`
		CantidadRespuestas int                      `json:"cantidad_respuestas"`
		RespuestasDetalle  []map[string]interface{} `json:"respuestas_detalle"` //Metadata, aqui puede asiganrse el valor, o el UID de los archivos según se requeira
	}

	var itemsPlantillasRespuestas []ItemPlantillaRespuesta

	if errFormulario == nil {

		if data, ok := response["Data"].([]interface{}); ok {
			var filteredFormularios []interface{}
			for _, formulario := range data {
				if formularioMap, ok := formulario.(map[string]interface{}); ok {
					if proyectoCurricularID, ok := formularioMap["ProyectoCurricularId"].(float64); ok {
						for _, idProyecto := range idsProyectos {
							if fmt.Sprintf("%.0f", proyectoCurricularID) == idProyecto {
								filteredFormularios = append(filteredFormularios, formularioMap)
								break
							}
						}
					}
				}
			}
			response["Data"] = filteredFormularios
		}

		if dataSource["campos"].(map[string]interface{})["componente"] != nil {

			var camposResponse map[string]interface{}
			errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=TipoCampoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["campos"].(map[string]interface{})["componente"]), &camposResponse)
			if errCampos == nil {

				if camposResponse["Data"] != nil {
					for _, campo := range camposResponse["Data"].([]interface{}) {
						campoId := fmt.Sprintf("%v", campo.(map[string]interface{})["Id"])
						campoIds = append(campoIds, campoId)

						var itemCampoResponse map[string]interface{}
						errItemCampo := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=CampoId:%s&Activo=true&limit=0", campoId), &itemCampoResponse)

						if errItemCampo == nil {
							if itemCampoResponse["Data"] != nil {
								for _, itemCampo := range itemCampoResponse["Data"].([]interface{}) {
									if itemCampo != nil && itemCampo.(map[string]interface{})["ItemId"] != nil {
										itemObj := itemCampo.(map[string]interface{})["ItemId"].(map[string]interface{})
										if itemObj["Id"] != nil {
											itemId := fmt.Sprintf("%v", itemObj["Id"])
											itemIdsCampos = append(itemIdsCampos, itemId)
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if dataSource["campos"].(map[string]interface{})["vinculacion"] != nil {
			var vinculacionResponse map[string]interface{}
			errVinculacion := request.GetJson("http://"+beego.AppConfig.String("PlanDocenteService")+fmt.Sprintf("plan_docente?query=tipo_vinculacion_id:%v&sortby=Id&order=asc&limit=0", dataSource["campos"].(map[string]interface{})["vinculacion"]), &vinculacionResponse)

			if errVinculacion == nil {
				if data, ok := vinculacionResponse["Data"].(map[string]interface{}); ok {
					var docenteIds []string

					for _, item := range data {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if docenteId, ok := itemMap["docente_id"].(string); ok {
								docenteIds = append(docenteIds, docenteId)
							}
						}
					}

					if response["Data"] != nil {
						var filteredFormularios []map[string]interface{}
						for _, formulario := range response["Data"].([]interface{}) {
							formularioMap := formulario.(map[string]interface{})
							evaluadoId := fmt.Sprintf("%v", formularioMap["EvaluadoId"])
							if contains(docenteIds, evaluadoId) {
								filteredFormularios = append(filteredFormularios, formularioMap)
							}
						}
						response["Data"] = filteredFormularios
					}
				}
			}
		}

		itemIds = itemIdsCampos

		if response["Data"] != nil {
			for _, formulario := range response["Data"].([]interface{}) {
				formularioMap := formulario.(map[string]interface{})
				formularioId := fmt.Sprintf("%v", formularioMap["Id"])
				formularioIds = append(formularioIds, formularioId)
			}
		}

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

							var formrespuestaResponse map[string]interface{}
							errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?sortby=Id&order=asc&limit=0"), &formrespuestaResponse)
							if errFormrespuesta == nil {

								var respuestasDetalle []map[string]interface{}
								if formrespuestaResponse["Data"] != nil {
									for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
										respuestaMap := respuesta.(map[string]interface{})
										plantillaRespId := fmt.Sprintf("%v", respuestaMap["PlantillaId"].(map[string]interface{})["Id"])

										if plantillaRespId == plantillaId && contains(formularioIds, fmt.Sprintf("%v", respuestaMap["FormularioId"].(map[string]interface{})["Id"])) {
											respuestaId := fmt.Sprintf("%v", respuestaMap["RespuestaId"].(map[string]interface{})["Id"])
											respuestasIds = append(respuestasIds, respuestaId)

											var respuestaDetalleResponse map[string]interface{}
											errRespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("respuesta/%s&limit=0", respuestaId), &respuestaDetalleResponse)
											if errRespuesta == nil {
												if respuestaDetalleResponse["Data"] != nil {
													respuestaDetalle := respuestaDetalleResponse["Data"].(map[string]interface{})
													metadataStr := respuestaDetalle["Metadata"]
													if metadataStr != nil {
														var metadataMap map[string]interface{}
														err := json.Unmarshal([]byte(metadataStr.(string)), &metadataMap)
														if err == nil {
															valor, ok := metadataMap["valor"]
															if !ok {
																valor = nil
															}
															campoId, ok := metadataMap["campo_id"]
															if !ok {
																campoId = nil
															}
															respuestasDetalle = append(respuestasDetalle, map[string]interface{}{
																"Metadata": metadataMap,
																"Valor":    valor,
																"CampoId":  campoId,
															})
														}
													}
												}
											}
										}
									}
								}

								itemsPlantillasRespuestas = append(itemsPlantillasRespuestas, ItemPlantillaRespuesta{
									ItemId:             itemId,
									PlantillaId:        plantillaId,
									CantidadRespuestas: len(respuestasIds),
									RespuestasDetalle:  respuestasDetalle,
								})
							}
						}
					}
				}
			}
		}
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, itemsPlantillasRespuestas, "")
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
	var campoIds []string
	var itemIds []string
	var itemIdsCampos []string
	var plantillaIds []string
	var respuestasIds []string
	var formularioIds []string

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["periodo_id"], dataSource["estudiante_id"]), &response)

	type ItemPlantillaRespuesta struct {
		ItemId             string                   `json:"item_id"`
		PlantillaId        string                   `json:"plantilla_id"`
		CantidadRespuestas int                      `json:"cantidad_respuestas"`
		RespuestasDetalle  []map[string]interface{} `json:"respuestas_detalle"` //Metadata, aqui puede asiganrse el valor, o el UID de los archivos según se requeira
	}

	var itemsPlantillasRespuestas []ItemPlantillaRespuesta

	if errFormulario == nil {

		if dataSource["campos"].(map[string]interface{})["componente"] != nil {

			var camposResponse map[string]interface{}
			errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=TipoCampoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["campos"].(map[string]interface{})["componente"]), &camposResponse)
			if errCampos == nil {

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
											itemIdsCampos = append(itemIdsCampos, itemId)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		itemIds = itemIdsCampos

		if response["Data"] != nil {
			for _, formulario := range response["Data"].([]interface{}) {
				formularioMap := formulario.(map[string]interface{})
				formularioId := fmt.Sprintf("%v", formularioMap["Id"])
				formularioIds = append(formularioIds, formularioId)
			}
		}

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

							var formrespuestaResponse map[string]interface{}
							errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?sortby=Id&order=asc&limit=0"), &formrespuestaResponse)
							if errFormrespuesta == nil {

								var respuestasDetalle []map[string]interface{}
								if formrespuestaResponse["Data"] != nil {
									for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
										respuestaMap := respuesta.(map[string]interface{})
										plantillaRespId := fmt.Sprintf("%v", respuestaMap["PlantillaId"].(map[string]interface{})["Id"])

										if plantillaRespId == plantillaId && contains(formularioIds, fmt.Sprintf("%v", respuestaMap["FormularioId"].(map[string]interface{})["Id"])) {
											respuestaId := fmt.Sprintf("%v", respuestaMap["RespuestaId"].(map[string]interface{})["Id"])
											respuestasIds = append(respuestasIds, respuestaId)

											var respuestaDetalleResponse map[string]interface{}
											errRespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("respuesta/%s&limit=0", respuestaId), &respuestaDetalleResponse)
											if errRespuesta == nil {
												if respuestaDetalleResponse["Data"] != nil {
													respuestaDetalle := respuestaDetalleResponse["Data"].(map[string]interface{})
													metadataStr := respuestaDetalle["Metadata"]
													if metadataStr != nil {
														var metadataMap map[string]interface{}
														err := json.Unmarshal([]byte(metadataStr.(string)), &metadataMap)
														if err == nil {
															valor, ok := metadataMap["valor"]
															if !ok {
																valor = nil
															}
															campoId, ok := metadataMap["campo_id"]
															if !ok {
																campoId = nil
															}
															respuestasDetalle = append(respuestasDetalle, map[string]interface{}{
																"Metadata": metadataMap,
																"Valor":    valor,
																"CampoId":  campoId,
															})
														}
													}
												}
											}
										}
									}
								}

								itemsPlantillasRespuestas = append(itemsPlantillasRespuestas, ItemPlantillaRespuesta{
									ItemId:             itemId,
									PlantillaId:        plantillaId,
									CantidadRespuestas: len(respuestasIds),
									RespuestasDetalle:  respuestasDetalle,
								})
							}
						}
					}
				}
			}
		}
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, itemsPlantillasRespuestas, "")
	return APIResponseDTO
}

func ReporteDocente(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var campoIds []string
	var itemIds []string
	var itemIdsCampos []string
	var plantillaIds []string
	var respuestasIds []string
	var formularioIds []string

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["periodo_id"], dataSource["docente_id"]), &response)

	type ItemPlantillaRespuesta struct {
		ItemId             string                   `json:"item_id"`
		PlantillaId        string                   `json:"plantilla_id"`
		CantidadRespuestas int                      `json:"cantidad_respuestas"`
		RespuestasDetalle  []map[string]interface{} `json:"respuestas_detalle"` //Metadata, aqui puede asiganrse el valor, o el UID de los archivos según se requeira
	}

	var itemsPlantillasRespuestas []ItemPlantillaRespuesta

	if errFormulario == nil {

		if dataSource["campos"].(map[string]interface{})["componente"] != nil {

			var camposResponse map[string]interface{}
			errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=TipoCampoId:%v&sortby=Id&order=asc&limit=0&Activo=true", dataSource["campos"].(map[string]interface{})["componente"]), &camposResponse)
			if errCampos == nil {

				if camposResponse["Data"] != nil {
					for _, campo := range camposResponse["Data"].([]interface{}) {
						campoId := fmt.Sprintf("%v", campo.(map[string]interface{})["Id"])
						campoIds = append(campoIds, campoId)

						var itemCampoResponse map[string]interface{}
						errItemCampo := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=CampoId:%s&Activo=true&limit=0", campoId), &itemCampoResponse)

						if errItemCampo == nil {
							if itemCampoResponse["Data"] != nil {
								for _, itemCampo := range itemCampoResponse["Data"].([]interface{}) {
									if itemCampo != nil && itemCampo.(map[string]interface{})["ItemId"] != nil {
										itemObj := itemCampo.(map[string]interface{})["ItemId"].(map[string]interface{})
										if itemObj["Id"] != nil {
											itemId := fmt.Sprintf("%v", itemObj["Id"])
											itemIdsCampos = append(itemIdsCampos, itemId)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		itemIds = itemIdsCampos

		if response["Data"] != nil {
			for _, formulario := range response["Data"].([]interface{}) {
				formularioMap := formulario.(map[string]interface{})
				formularioId := fmt.Sprintf("%v", formularioMap["Id"])
				formularioIds = append(formularioIds, formularioId)
			}
		}

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

							var formrespuestaResponse map[string]interface{}
							errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?sortby=Id&order=asc&limit=0"), &formrespuestaResponse)
							if errFormrespuesta == nil {

								var respuestasDetalle []map[string]interface{}
								if formrespuestaResponse["Data"] != nil {
									for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
										respuestaMap := respuesta.(map[string]interface{})
										plantillaRespId := fmt.Sprintf("%v", respuestaMap["PlantillaId"].(map[string]interface{})["Id"])

										if plantillaRespId == plantillaId && contains(formularioIds, fmt.Sprintf("%v", respuestaMap["FormularioId"].(map[string]interface{})["Id"])) {
											respuestaId := fmt.Sprintf("%v", respuestaMap["RespuestaId"].(map[string]interface{})["Id"])
											respuestasIds = append(respuestasIds, respuestaId)

											var respuestaDetalleResponse map[string]interface{}
											errRespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("respuesta/%s&limit=0", respuestaId), &respuestaDetalleResponse)
											if errRespuesta == nil {
												if respuestaDetalleResponse["Data"] != nil {
													respuestaDetalle := respuestaDetalleResponse["Data"].(map[string]interface{})
													metadataStr := respuestaDetalle["Metadata"]
													if metadataStr != nil {
														var metadataMap map[string]interface{}
														err := json.Unmarshal([]byte(metadataStr.(string)), &metadataMap)
														if err == nil {
															valor, ok := metadataMap["valor"]
															if !ok {
																valor = nil
															}
															campoId, ok := metadataMap["campo_id"]
															if !ok {
																campoId = nil
															}
															respuestasDetalle = append(respuestasDetalle, map[string]interface{}{
																"Metadata": metadataMap,
																"Valor":    valor,
																"CampoId":  campoId,
															})
														}
													}
												}
											}
										}
									}
								}

								itemsPlantillasRespuestas = append(itemsPlantillasRespuestas, ItemPlantillaRespuesta{
									ItemId:             itemId,
									PlantillaId:        plantillaId,
									CantidadRespuestas: len(respuestasIds),
									RespuestasDetalle:  respuestasDetalle,
								})
							}
						}
					}
				}
			}
		}
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, itemsPlantillasRespuestas, "")
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
