package services

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// id tipo formulario hace referencia a proceso_id de la tabla plantilla
func ConsultaFormulario(id_tipo_formulario string, id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {

	var plantilla map[string]interface{}
	errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=ProcesoId:%v&Activo:true&sortby=Id&order=asc&limit=0", id_tipo_formulario), &plantilla)
	if errPlantilla != nil || fmt.Sprintf("%v", plantilla) == "[map[]]" {
		return helpers.ErrEmiter(errPlantilla, fmt.Sprintf("%v", plantilla))
	}

	var itemCampos map[string]interface{}
	errItemCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &itemCampos)
	if errItemCampos != nil || fmt.Sprintf("%v", itemCampos) == "[map[]]" {
		return helpers.ErrEmiter(errItemCampos, fmt.Sprintf("%v", itemCampos))
	}

	var campos map[string]interface{}
	errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &campos)
	if errCampos != nil || fmt.Sprintf("%v", campos) == "[map[]]" {
		return helpers.ErrEmiter(errCampos, fmt.Sprintf("%v", campos))
	}

	secciones := []map[string]interface{}{}
	data := plantilla["Data"].([]interface{})
	itemCamposMap := make(map[int][]map[string]interface{})
	itemCamposData := itemCampos["Data"].([]interface{})
	camposData := campos["Data"].([]interface{})
	for _, itemCampo := range itemCamposData {
		itemCampoMap := itemCampo.(map[string]interface{})
		itemId := int(itemCampoMap["ItemId"].(map[string]interface{})["Id"].(float64))
		campo := itemCampoMap["CampoId"].(map[string]interface{})
		campoId := int(campo["Id"].(float64))
		tipoCampo := int(campo["TipoCampoId"].(float64))
		campoInfo := map[string]interface{}{
			"nombre":     campo["Nombre"].(string),
			"campo_id":   campoId,
			"tipo_campo": tipoCampo,
			"valor":      campo["Valor"],
			"porcentaje": itemCampoMap["Porcentaje"],
			"escala":     obtenerCamposHijos(campoId, camposData),
		}
		if tipoCampo == 6686 && id_tipo_formulario == "5" { // 6 es carga de archivos
			descargaArchivos := obtenerDescargaArchivos(id_tercero, id_espacio)
			for key, value := range descargaArchivos {
				campoInfo[key] = value
				campoInfo["nombre"] = "descarga_archivos"
				campoInfo["tipo_campo"] = 4672 // 5 es descarga de archivos
			}
		}
		itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
	}

	for _, item := range data {
		itemMap := item.(map[string]interface{})
		seccion := itemMap["SeccionId"].(map[string]interface{})
		seccionId := int(seccion["Id"].(float64))

		var seccionEncontrada map[string]interface{}
		for _, sec := range secciones {
			if sec["id"] == seccionId {
				seccionEncontrada = sec
				break
			}
		}

		if seccionEncontrada == nil {
			seccionNueva := map[string]interface{}{
				"id":     seccionId,
				"nombre": seccion["Nombre"].(string),
				"orden":  int(seccion["Orden"].(float64)),
				"items":  []map[string]interface{}{},
			}
			secciones = append(secciones, seccionNueva)
			seccionEncontrada = seccionNueva
		}

		itemId := int(itemMap["ItemId"].(map[string]interface{})["Id"].(float64))
		itemOrden := int(itemMap["ItemId"].(map[string]interface{})["Orden"].(float64))
		itemInfo := map[string]interface{}{
			"id":     itemId,
			"nombre": itemMap["ItemId"].(map[string]interface{})["Nombre"].(string),
			"orden":  itemOrden,
			"campos": itemCamposMap[itemId],
		}
		seccionEncontrada["items"] = append(seccionEncontrada["items"].([]map[string]interface{}), itemInfo)
	}

	for _, seccion := range secciones {
		items := seccion["items"].([]map[string]interface{})
		sort.Slice(items, func(i, j int) bool {
			if items[i]["orden"].(int) == items[j]["orden"].(int) {
				return items[i]["id"].(int) < items[j]["id"].(int)
			}
			return items[i]["orden"].(int) < items[j]["orden"].(int)
		})
	}

	sort.Slice(secciones, func(i, j int) bool {
		if secciones[i]["orden"].(int) == secciones[j]["orden"].(int) {
			return secciones[i]["id"].(int) < secciones[j]["id"].(int)
		}
		return secciones[i]["orden"].(int) < secciones[j]["orden"].(int)
	})

	response := map[string]interface{}{
		"docente":          id_tercero,
		"espacioAcademico": id_espacio,
		"seccion":          secciones,
		"tipoEvaluacion":   id_tipo_formulario,
		"evaluacion":       id_tipo_formulario,
	}

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func obtenerCamposHijos(campoId int, camposData []interface{}) []map[string]interface{} {
	var hijos []map[string]interface{}

	var campoPadre map[string]interface{}
	for _, campo := range camposData {
		campoMap := campo.(map[string]interface{})
		if int(campoMap["Id"].(float64)) == campoId {
			campoPadre = campoMap
			break
		}
	}
	if campoPadre == nil {
		hijoInfo := map[string]interface{}{
			"nombre":     campoPadre["Nombre"].(string),
			"tipo_campo": int(campoPadre["TipoCampoId"].(float64)),
			"valor":      campoPadre["Valor"],
		}
		hijos = append(hijos, hijoInfo)
	}
	if campoPadre != nil {
		for _, campo := range camposData {
			campoMap := campo.(map[string]interface{})
			if int(campoMap["CampoPadreId"].(float64)) == campoId {
				hijoInfo := map[string]interface{}{
					"nombre":     campoMap["Nombre"].(string),
					"tipo_campo": int(campoMap["TipoCampoId"].(float64)),
					"valor":      campoMap["Valor"],
					"campo_id":   campoMap["Id"],
				}
				hijos = append(hijos, hijoInfo)
			}
		}
	}

	return hijos
}
func obtenerDescargaArchivos(id_tercero string, id_espacio string) map[string]interface{} {

	var campoIds []string
	var itemIds []string
	var itemIdsCampos []string
	var plantillaIds []string
	var respuestasIds []string
	var formularioIds []string
	var documentos []string

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=EvaluadoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_tercero), &response)

	if errFormulario == nil {

		var camposResponse map[string]interface{}
		errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=TipoCampoId:6&sortby=Id&order=asc&limit=0&Activo=true"), &camposResponse)
		if errCampos == nil {

			for _, campo := range camposResponse["Data"].([]interface{}) {
				campoId := fmt.Sprintf("%v", campo.(map[string]interface{})["Id"])
				campoIds = append(campoIds, campoId)

				var itemCampoResponse map[string]interface{}
				errItemCampo := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=CampoId:%s&Activo=true&limit=0&order=asc", campoId), &itemCampoResponse)
				if errItemCampo == nil && itemCampoResponse["Data"] != nil {
					for _, itemCampo := range itemCampoResponse["Data"].([]interface{}) {
						if itemCampoMap, ok := itemCampo.(map[string]interface{}); ok {
							if itemObj, ok := itemCampoMap["ItemId"].(map[string]interface{}); ok && itemObj["Id"] != nil {
								itemId := fmt.Sprintf("%v", itemObj["Id"])
								itemIdsCampos = append(itemIdsCampos, itemId)
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
		if errPlantilla == nil && plantillaResponse["Data"] != nil {
			for _, plantilla := range plantillaResponse["Data"].([]interface{}) {

				if itemPlantilla, ok := plantilla.(map[string]interface{})["ItemId"].(map[string]interface{}); ok {
					itemId := fmt.Sprintf("%v", itemPlantilla["Id"])
					for _, id := range itemIds {
						if itemId == id {
							plantillaId := fmt.Sprintf("%v", plantilla.(map[string]interface{})["Id"])
							plantillaIds = append(plantillaIds, plantillaId)

							var formrespuestaResponse map[string]interface{}
							errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?sortby=Id&order=asc&limit=0"), &formrespuestaResponse)
							if errFormrespuesta == nil && formrespuestaResponse["Data"] != nil {

								for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
									respuestaMap := respuesta.(map[string]interface{})
									if plantillaRespMap, ok := respuestaMap["PlantillaId"].(map[string]interface{}); ok {
										if _, ok := plantillaRespMap["Id"]; ok {
											formularioIdMap, ok := respuestaMap["FormularioId"].(map[string]interface{})
											if ok && contains(formularioIds, fmt.Sprintf("%v", formularioIdMap["Id"])) {
												respuestaIdMap, ok := respuestaMap["RespuestaId"].(map[string]interface{})
												if ok {
													respuestaId := fmt.Sprintf("%v", respuestaIdMap["Id"])
													respuestasIds = append(respuestasIds, respuestaId)

													var respuestaDetalleResponse map[string]interface{}
													errRespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("respuesta/%s&order=asc&limit=0", respuestaId), &respuestaDetalleResponse)

													if errRespuesta == nil && respuestaDetalleResponse["Data"] != nil {
														respuestaDetalle := respuestaDetalleResponse["Data"].(map[string]interface{})

														if metadataStr, ok := respuestaDetalle["Metadata"].(string); ok && metadataStr != "" {
															var metadataMap map[string]interface{}
															err := json.Unmarshal([]byte(metadataStr), &metadataMap)
															if err == nil {
																archivos, ok := metadataMap["archivos"].([]interface{})
																if ok {
																	for _, archivo := range archivos {
																		if archivoStr, ok := archivo.(string); ok {
																			documentos = append(documentos, archivoStr)
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}

								}
							}
						}
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"UIDs": documentos,
	}
}
