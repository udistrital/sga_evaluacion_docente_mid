package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// id tipo formulario hace referencia a proceso_id de la tabla plantilla
func ConsultaFormulario(id_tipo_formulario string, id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {

	var formularioID int
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
		/*if tipoCampo == 6686 && id_tipo_formulario == "5" { // 6 es carga de archivos
			descargaArchivos := obtenerDescargaArchivos(id_tercero, id_espacio)
			for key, value := range descargaArchivos {
				campoInfo[key] = value
				campoInfo["nombre"] = "descarga_archivos"
				campoInfo["tipo_campo"] = 4672 // 5 es descarga de archivos
			}
		}*/
		itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
	}

	var res map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v,EspacioAcademicoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_periodo, id_tercero, id_espacio), &res)

	if errFormulario == nil {
		if data, ok := res["Data"].([]interface{}); ok && len(data) > 0 {
			if formulario, ok := data[0].(map[string]interface{}); ok && len(formulario) > 0 {
				if id, exists := formulario["Id"].(float64); exists {
					formularioID = int(id)
				}
			}
		}
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
		if formularioID > 0 {
			existe := VerificarRespuesta(formularioID, itemId)
			if existe.Status == 200 {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Ya se han registrado respuestas para este formulario"))
				return APIResponseDTO
			}
		}
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
func obtenerDescargaArchivos(id_tercero string, id_espacio string, itemId string) map[string]interface{} {
	var formularioIds []string
	var documentos []string
	fmt.Println("aaaaaaaaaa")
	fmt.Println(itemId)
	fmt.Println("aaaaaaaaaa")

	var response map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=EvaluadoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_tercero), &response)

	if errFormulario == nil {

		if response["Data"] != nil {
			for _, formulario := range response["Data"].([]interface{}) {
				formularioMap := formulario.(map[string]interface{})
				formularioId := fmt.Sprintf("%v", formularioMap["Id"])
				formularioIds = append(formularioIds, formularioId)
			}
		}

		var plantillaResponse map[string]interface{}
		errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=ItemId__Id:%v&sortby=Id&order=asc&limit=0", itemId), &plantillaResponse)
		if errPlantilla == nil && plantillaResponse["Data"] != nil {

			idPlantilla := fmt.Sprintf("%v", plantillaResponse["Data"].([]interface{})[0].(map[string]interface{})["Id"])
			var formrespuestaResponse map[string]interface{}
			errFormrespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?query=PlantillaId:%v&sortby=Id&order=asc&limit=0", idPlantilla), &formrespuestaResponse)
			if errFormrespuesta == nil && formrespuestaResponse["Data"] != nil {

				for _, respuesta := range formrespuestaResponse["Data"].([]interface{}) {
					respuestaMap := respuesta.(map[string]interface{})
					if plantillaRespMap, ok := respuestaMap["PlantillaId"].(map[string]interface{}); ok {

						if _, ok := plantillaRespMap["Id"]; ok {

							formularioIdMap, ok := respuestaMap["FormularioId"].(map[string]interface{})
							if ok && contains(formularioIds, fmt.Sprintf("%v", formularioIdMap["Id"])) {

								respuestaIdMap, ok := respuestaMap["RespuestaId"].(map[string]interface{})
								if ok {
									fmt.Println(formularioIdMap["Id"])
									respuestaId := fmt.Sprintf("%v", respuestaIdMap["Id"])
									fmt.Println(respuestaId)
									var respuestaDetalleResponse map[string]interface{}
									errRespuesta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/respuesta?query=Id:%v&limit=0", respuestaId), &respuestaDetalleResponse)

									if errRespuesta == nil && respuestaDetalleResponse["Data"] != nil {
										if respuestas, ok := respuestaDetalleResponse["Data"].([]interface{}); ok && len(respuestas) > 0 {
											respuestaDetalle := respuestas[0].(map[string]interface{})

											if metadataStr, ok := respuestaDetalle["Metadata"].(string); ok && metadataStr != "" {
												var metadataMap map[string]interface{}
												err := json.Unmarshal([]byte(metadataStr), &metadataMap)
												if err == nil {
													if archivos, ok := metadataMap["archivos"].([]interface{}); ok {
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

	return map[string]interface{}{
		"UIDs": documentos,
	}
}

func CrearFormulario(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var revertir bool = false
	var itemIDs []float64
	var plantillaIDs []float64

	if err := json.Unmarshal(data, &dataSource); err != nil {
		return helpers.ErrEmiter(err, "error al deserializar los datos")
	}

	secciones, ok := dataSource["secciones"].([]interface{})
	if ok {
		for _, seccion := range secciones {
			secMap, ok := seccion.(map[string]interface{})
			if ok {
				nombreSeccion := secMap["nombre"]
				ordenSeccion := secMap["orden"]

				nuevaSec := map[string]interface{}{
					"Activo": true,
					"Nombre": nombreSeccion,
					"Orden":  ordenSeccion,
				}
				var newSec map[string]interface{}
				errResSec := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/seccion/", "POST", &newSec, nuevaSec)
				if errResSec != nil {
					revertir = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una de las secciones")
					return APIResponseDTO
				}
				seccionID := newSec["Data"].(map[string]interface{})["Id"].(float64)

				fmt.Printf("Sección: %v, Orden: %v, SecID: %v\n", nombreSeccion, ordenSeccion, seccionID)

				items, ok := secMap["items"].([]interface{})
				if ok {

					for _, item := range items {
						itemMap, ok := item.(map[string]interface{})
						if ok {
							nombreItem := itemMap["nombre"]
							ordenItem := itemMap["orden"]
							campoID := itemMap["campo_id"]
							porcentaje := itemMap["porcentaje"]
							nuevoItem := map[string]interface{}{
								"Activo":     true,
								"Nombre":     nombreItem,
								"Orden":      ordenItem,
								"Porcentaje": porcentaje,
							}
							var newItem map[string]interface{}
							errResItem := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/item/", "POST", &newItem, nuevoItem)
							if errResItem != nil {
								revertir = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar uno de los items")
								return APIResponseDTO
							}
							itemID := newItem["Data"].(map[string]interface{})["Id"].(float64)
							itemIDs = append(itemIDs, itemID)

							nuevoItemCampo := map[string]interface{}{
								"Activo":     true,
								"CampoId":    map[string]interface{}{"Id": campoID},
								"ItemId":     map[string]interface{}{"Id": itemID},
								"Porcentaje": porcentaje,
							}
							var newItemCampo map[string]interface{}
							errResItemCampo := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/item_campo/", "POST", &newItemCampo, nuevoItemCampo)
							if errResItemCampo != nil {
								revertir = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar uno de los items_campo")
								return APIResponseDTO
							}
							fmt.Printf("  Ítem: %v, Orden: %v, Campo ID: %v,Item ID: %v, Porcentaje: %v\n", nombreItem, ordenItem, campoID, itemID, porcentaje)
							nuevaPlantilla := map[string]interface{}{
								"Activo":       true,
								"SeccionId":    map[string]interface{}{"Id": seccionID},
								"ItemId":       map[string]interface{}{"Id": itemID},
								"ProcesoId":    dataSource["proceso_id"],
								"EstructuraId": dataSource["estructura"],
							}
							var newPlantilla map[string]interface{}
							errResPlantilla := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/plantilla/", "POST", &newPlantilla, nuevaPlantilla)
							if errResPlantilla != nil {
								revertir = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una plantilla")
								return APIResponseDTO
							}
							plantillaID := newPlantilla["Data"].(map[string]interface{})["Id"].(float64)
							plantillaIDs = append(plantillaIDs, plantillaID)
						}
					}
				}
			}
		}
	}

	if revertir {
		fmt.Println("IDs de las Plantillas:", plantillaIDs)

		if len(plantillaIDs) > 0 {
			for _, id := range plantillaIDs {
				var plantilla map[string]interface{}

				errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla?query=Id:%v&Activo:true&sortby=Id&order=asc&limit=0", id), &plantilla)
				if errPlantilla == nil {

					plantillaData := plantilla["Data"].(map[string]interface{})
					plantillaData["Activo"] = false

					var inactivaPlantilla map[string]interface{}
					errPut := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla/%v", id), "PUT", &inactivaPlantilla, plantillaData)
					if errPut != nil {
						return helpers.ErrEmiter(errPut, fmt.Sprintf("Error actualizando plantilla con ID %v: %v", id, errPut))
					}
				}
			}
		}
	}

	return requestresponse.APIResponseDTO(true, 200, dataSource, "Se ha registrado el formulario c:")
}

func FormularioCoevaluacion(id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {
	var formularioID int
	var plantilla map[string]interface{}
	errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=ProcesoId:5&Activo:true&sortby=Id&order=asc&limit=0"), &plantilla)
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
		campo := itemCampoMap["CampoId"].(map[string]interface{})
		campoId := int(campo["Id"].(float64))
		tipoCampo := int(campo["TipoCampoId"].(float64))
		itemId := int(itemCampoMap["ItemId"].(map[string]interface{})["Id"].(float64))
		campoInfo := map[string]interface{}{
			"nombre":     campo["Nombre"].(string),
			"campo_id":   campoId,
			"tipo_campo": tipoCampo,
			"valor":      campo["Valor"],
		}
		if tipoCampo == 6 {
			itemRel := int(itemCampoMap["Porcentaje"].(float64))
			descargaArchivos := obtenerDescargaArchivos(id_tercero, id_espacio, strconv.Itoa(itemRel))
			for key, value := range descargaArchivos {
				campoInfo[key] = value
				fmt.Println("//////////////////")
				fmt.Println(key)
				fmt.Println(value)
				fmt.Println("//////////////////")
				campoInfo["nombre"] = "descarga_archivos"
				campoInfo["tipo_campo"] = 4672
			}
			itemCamposMap[itemId] = append(itemCamposMap[itemRel], campoInfo)
		} else {
			campoInfo["porcentaje"] = itemCampoMap["Porcentaje"]
			campoInfo["escala"] = obtenerCamposHijos(campoId, camposData)
			itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
		}

	}

	var res map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v,EspacioAcademicoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_periodo, id_tercero, id_espacio), &res)

	if errFormulario == nil {
		if data, ok := res["Data"].([]interface{}); ok && len(data) > 0 {
			if formulario, ok := data[0].(map[string]interface{}); ok && len(formulario) > 0 {
				if id, exists := formulario["Id"].(float64); exists {
					formularioID = int(id)
				}
			}
		}
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
		if formularioID > 0 {
			existe := VerificarRespuesta(formularioID, itemId)
			if existe.Status == 200 {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Ya se han registrado respuestas para este formulario"))
				return APIResponseDTO
			}
		}
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
	}

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func CrearFormularioCo(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var revertir bool = false
	var itemIDs []float64
	var plantillaIDs []float64
	nombreItemMap := make(map[string]float64)

	if err := json.Unmarshal(data, &dataSource); err != nil {
		return helpers.ErrEmiter(err, "error al deserializar los datos")
	}

	secciones, ok := dataSource["secciones"].([]interface{})
	if ok {
		for _, seccion := range secciones {
			secMap, ok := seccion.(map[string]interface{})
			if ok {
				nombreSeccion := secMap["nombre"]
				ordenSeccion := secMap["orden"]

				nuevaSec := map[string]interface{}{
					"Activo": true,
					"Nombre": nombreSeccion,
					"Orden":  ordenSeccion,
				}
				var newSec map[string]interface{}
				errResSec := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/seccion/", "POST", &newSec, nuevaSec)
				if errResSec != nil {
					revertir = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una de las secciones")
					return APIResponseDTO
				}
				seccionID := newSec["Data"].(map[string]interface{})["Id"].(float64)

				fmt.Printf("Sección: %v, Orden: %v, SecID: %v\n", nombreSeccion, ordenSeccion, seccionID)

				items, ok := secMap["items"].([]interface{})
				if ok {

					for _, item := range items {
						itemMap, ok := item.(map[string]interface{})
						if ok {
							nombreItem := itemMap["nombre"].(string)
							ordenItem := itemMap["orden"]
							campoID := itemMap["campo_id"].(float64)
							porcentaje := itemMap["porcentaje"]
							if campoID == 4672 {
								porcentaje = itemMap["item_relacion_id"].(float64)
							}

							var itemID float64
							if idExistente, existe := nombreItemMap[nombreItem]; existe {
								fmt.Printf("Item ya existe: %s, usando itemID: %v\n", nombreItem, idExistente)
								itemID = idExistente
							} else {
								nuevoItem := map[string]interface{}{
									"Activo": true,
									"Nombre": nombreItem,
									"Orden":  ordenItem,
								}
								var newItem map[string]interface{}
								errResItem := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/item/", "POST", &newItem, nuevoItem)
								if errResItem != nil {
									revertir = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar uno de los items")
									return APIResponseDTO
								}
								itemID = newItem["Data"].(map[string]interface{})["Id"].(float64)
								itemIDs = append(itemIDs, itemID)
								nombreItemMap[nombreItem] = itemID
								fmt.Printf("Nuevo item creado: %s, itemID: %v\n", nombreItem, itemID)
							}
							nuevoItemCampo := map[string]interface{}{
								"Activo":     true,
								"CampoId":    map[string]interface{}{"Id": campoID},
								"ItemId":     map[string]interface{}{"Id": itemID},
								"Porcentaje": porcentaje,
							}
							var newItemCampo map[string]interface{}
							errResItemCampo := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/item_campo/", "POST", &newItemCampo, nuevoItemCampo)
							if errResItemCampo != nil {
								revertir = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar uno de los items_campo")
								return APIResponseDTO
							}
							nuevaPlantilla := map[string]interface{}{
								"Activo":       true,
								"SeccionId":    map[string]interface{}{"Id": seccionID},
								"ItemId":       map[string]interface{}{"Id": itemID},
								"ProcesoId":    dataSource["proceso_id"],
								"EstructuraId": dataSource["estructura"],
							}
							var newPlantilla map[string]interface{}
							errResPlantilla := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+"/plantilla/", "POST", &newPlantilla, nuevaPlantilla)
							if errResPlantilla != nil {
								revertir = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una plantilla")
								return APIResponseDTO
							}
							plantillaID := newPlantilla["Data"].(map[string]interface{})["Id"].(float64)
							plantillaIDs = append(plantillaIDs, plantillaID)

						}
					}
				}
			}
		}
	}

	if revertir {
		fmt.Println("IDs de las Plantillas:", plantillaIDs)

		if len(plantillaIDs) > 0 {
			for _, id := range plantillaIDs {
				var plantilla map[string]interface{}

				errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla?query=Id:%v&Activo:true&sortby=Id&order=asc&limit=0", id), &plantilla)
				if errPlantilla == nil {

					plantillaData := plantilla["Data"].(map[string]interface{})
					plantillaData["Activo"] = false

					var inactivaPlantilla map[string]interface{}
					errPut := request.SendJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla/%v", id), "PUT", &inactivaPlantilla, plantillaData)
					if errPut != nil {
						return helpers.ErrEmiter(errPut, fmt.Sprintf("Error actualizando plantilla con ID %v: %v", id, errPut))
					}
				}
			}
		}
	}

	return requestresponse.APIResponseDTO(true, 200, dataSource, "Se ha registrado el formulario c:")
}
