package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"errors"
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// id tipo formulario hace referencia a proceso_id de la tabla plantilla
/*func ConsultaFormulario(id_tipo_formulario string, id_periodo string, id_tercero string, id_espacio string, id_grupo string) (APIResponseDTO requestresponse.APIResponse) {

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
		itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
	}

	// Construir la consulta dinámicamente
	query := fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v,EspacioAcademicoId:%v,PlantillaProcesoId:%v", id_periodo, id_tercero, id_espacio, id_tipo_formulario)
	if id_grupo != "" {
		query += fmt.Sprintf(",Grupos.id:%v", id_grupo)
	}
	query += "&sortby=Id&order=asc&limit=0&Activo=true"

	fmt.Println("Query: ", query)

	var res map[string]interface{}
	errFormulario := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+query, &res)

	if errFormulario == nil {
		if data, ok := res["Data"].([]interface{}); ok && len(data) > 0 {
			if formulario, ok := data[0].(map[string]interface{}); ok && len(formulario) > 0 {
				if id, exists := formulario["Id"].(float64); exists {
					formularioID = int(id)
				}
			}
		}
	}

	// Si el formulario ya existe, retorna un error
	if formularioID > 0 {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Ya existe un formulario para este grupo")
		return APIResponseDTO
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
}*/

const (
	HttpPrefix = "http://"
	EmptyMapString = "[map[]]"
)

// id tipo formulario hace referencia a proceso_id de la tabla plantilla
func ConsultaFormulario(idTipoFormulario, idPeriodo, idTercero, idEspacio, idGrupo string) (APIResponseDTO requestresponse.APIResponse) {
	plantilla, err := obtenerPlantilla(idTipoFormulario)
	if err != nil {
		return helpers.ErrEmiter(err, "")
	}

	itemCampos, err := obtenerItemCampos()
	if err != nil {
		return helpers.ErrEmiter(err, "")
	}

	campos, err := obtenerCampos()
	if err != nil {
		return helpers.ErrEmiter(err, "")
	}

	formularioID, err := verificarFormularioExistente(idPeriodo, idTercero, idEspacio, idTipoFormulario, idGrupo)
	if err != nil {
		return helpers.ErrEmiter(err, "")
	}
	if formularioID > 0 {
		return requestresponse.APIResponseDTO(false, 400, nil, "Ya existe un formulario para este grupo")
	}

	itemCamposMap := mapearCamposPorItem(itemCampos["Data"].([]interface{}), campos["Data"].([]interface{}))

	secciones := procesarSecciones(plantilla["Data"].([]interface{}), itemCamposMap)
	ordenarSeccionesYItems(secciones)

	response := map[string]interface{}{
		"docente":          idTercero,
		"espacioAcademico": idEspacio,
		"seccion":          secciones,
		"tipoEvaluacion":   idTipoFormulario,
		"evaluacion":       idTipoFormulario,
	}
	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func procesarSecciones(data []interface{}, itemCamposMap map[int][]map[string]interface{}) []map[string]interface{} {
	secciones := []map[string]interface{}{}

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

	return secciones
}

func ordenarSeccionesYItems(secciones []map[string]interface{}) {
	for _, seccion := range secciones {
		items := seccion["items"].([]map[string]interface{})
		sort.Slice(items, func(i, j int) bool {
			if items[i]["orden"].(int) == items[j]["orden"].(int) {
				return items[i]["id"].(int) < items[j]["id"].(int)
			}
			return items[i]["orden"].(int) < items[j]["orden"].(int)
		})
		seccion["items"] = items
	}

	sort.Slice(secciones, func(i, j int) bool {
		if secciones[i]["orden"].(int) == secciones[j]["orden"].(int) {
			return secciones[i]["id"].(int) < secciones[j]["id"].(int)
		}
		return secciones[i]["orden"].(int) < secciones[j]["orden"].(int)
	})
}

func mapearCamposPorItem(itemCamposData, camposData []interface{}) map[int][]map[string]interface{} {
	itemCamposMap := make(map[int][]map[string]interface{})
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
		itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
	}
	return itemCamposMap
}


func verificarFormularioExistente(idPeriodo, idTercero, idEspacio, idTipoFormulario, idGrupo string) (int, error) {
	query := fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v,EspacioAcademicoId:%v,PlantillaProcesoId:%v",
		idPeriodo, idTercero, idEspacio, idTipoFormulario)
	if idGrupo != "" {
		query += fmt.Sprintf(",Grupos.id:%v", idGrupo)
	}
	query += "&sortby=Id&order=asc&limit=0&Activo=true"

	url := fmt.Sprintf("%s%s/%s", HttpPrefix, beego.AppConfig.String("EvaluacionDocenteService"), query)
	var res map[string]interface{}
	err := request.GetJson(url, &res)
	if err != nil {
		return 0, err
	}
	if data, ok := res["Data"].([]interface{}); ok && len(data) > 0 {
		if formulario, ok := data[0].(map[string]interface{}); ok && len(formulario) > 0 {
			if id, exists := formulario["Id"].(float64); exists {
				return int(id), nil
			}
		}
	}
	return 0, nil
}

func obtenerCampos() (map[string]interface{}, error) {
	var campos map[string]interface{}
	url := fmt.Sprintf("%s%s/campo?query=Activo:true&sortby=Id&order=asc&limit=0",
		HttpPrefix,beego.AppConfig.String("EvaluacionDocenteService"))
	err := request.GetJson(url, &campos)
	if err != nil || fmt.Sprintf("%v", campos) == EmptyMapString {
		return nil, fmt.Errorf("error al obtener campos: %v", err)
	}
	return campos, nil
}

func obtenerItemCampos() (map[string]interface{}, error) {
	var itemCampos map[string]interface{}
	url := fmt.Sprintf("%s%s/item_campo?query=Activo:true&sortby=Id&order=asc&limit=0",
		HttpPrefix, beego.AppConfig.String("EvaluacionDocenteService"))
	err := request.GetJson(url, &itemCampos)
	if err != nil || fmt.Sprintf("%v", itemCampos) == EmptyMapString {
		return nil, fmt.Errorf("error al obtener itemCampos: %v", err)
	}
	return itemCampos, nil
}

func obtenerPlantilla(idTipoFormulario string) (map[string]interface{}, error) {
	var plantilla map[string]interface{}
	url := fmt.Sprintf("%s%s/plantilla?query=ProcesoId:%v&Activo:true&sortby=Id&order=asc&limit=0",
		HttpPrefix, beego.AppConfig.String("EvaluacionDocenteService"), idTipoFormulario)
	err := request.GetJson(url, &plantilla)
	if err != nil || fmt.Sprintf("%v", plantilla) == EmptyMapString {
		return nil, fmt.Errorf("error al obtener plantilla: %v", err)
	}
	return plantilla, nil
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


/*func obtenerDescargaArchivos(id_tercero string, id_espacio string, itemId string) map[string]interface{} {
	var formularioIds []string
	var documentos []string

	var response map[string]interface{}
	errFormulario := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=EvaluadoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_tercero), &response)

	if errFormulario == nil {

		if response["Data"] != nil {
			for _, formulario := range response["Data"].([]interface{}) {
				formularioMap := formulario.(map[string]interface{})
				formularioId := fmt.Sprintf("%v", formularioMap["Id"])
				formularioIds = append(formularioIds, formularioId)
			}
		}

		var plantillaResponse map[string]interface{}
		errPlantilla := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=ItemId__Id:%v&sortby=Id&order=asc&limit=0", itemId), &plantillaResponse)
		if errPlantilla == nil && plantillaResponse["Data"] != nil {

			idPlantilla := fmt.Sprintf("%v", plantillaResponse["Data"].([]interface{})[0].(map[string]interface{})["Id"])
			var formrespuestaResponse map[string]interface{}
			errFormrespuesta := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formrespuesta?query=PlantillaId:%v&sortby=Id&order=asc&limit=0", idPlantilla), &formrespuestaResponse)
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
									var respuestaDetalleResponse map[string]interface{}
									errRespuesta := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/respuesta?query=Id:%v&limit=0", respuestaId), &respuestaDetalleResponse)

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
}*/


func obtenerDescargaArchivos(idTercero string, idEspacio string, itemId string) map[string]interface{} {
	formularioIds, err := obtenerFormulariosIds(idTercero)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	idPlantilla, err := obtenerIdPlantillaPorItem(itemId)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	respuestas, err := obtenerRespuestasPorPlantilla(idPlantilla)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	documentos, err := obtenerDocumentosDesdeRespuestas(respuestas, formularioIds)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	return map[string]interface{}{
		"UIDs": documentos,
	}
}

func obtenerFormulariosIds(idTercero string) ([]string, error) {
	var response map[string]interface{}
	url := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + fmt.Sprintf("formulario?query=EvaluadoId:%v&sortby=Id&order=asc&limit=0&Activo=true", idTercero)
	err := request.GetJson(url, &response)
	if err != nil {
		return nil, err
	}

	var formularioIds []string
	if response["Data"] != nil {
		for _, formulario := range response["Data"].([]interface{}) {
			formularioMap := formulario.(map[string]interface{})
			formularioId := fmt.Sprintf("%v", formularioMap["Id"])
			formularioIds = append(formularioIds, formularioId)
		}
	}
	return formularioIds, nil
}

func obtenerIdPlantillaPorItem(itemId string) (string, error) {
	var plantillaResponse map[string]interface{}
	url := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + fmt.Sprintf("plantilla?query=ItemId__Id:%v&sortby=Id&order=asc&limit=0", itemId)
	err := request.GetJson(url, &plantillaResponse)
	if err != nil {
		return "", err
	}

	if plantillaResponse["Data"] != nil && len(plantillaResponse["Data"].([]interface{})) > 0 {
		idPlantilla := fmt.Sprintf("%v", plantillaResponse["Data"].([]interface{})[0].(map[string]interface{})["Id"])
		return idPlantilla, nil
	}
	return "", fmt.Errorf("no se encontró plantilla para item %s", itemId)
}

func obtenerRespuestasPorPlantilla(idPlantilla string) ([]map[string]interface{}, error) {
	var formrespuestaResponse map[string]interface{}
	url := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + fmt.Sprintf("formrespuesta?query=PlantillaId:%v&sortby=Id&order=asc&limit=0", idPlantilla)
	err := request.GetJson(url, &formrespuestaResponse)
	if err != nil {
		return nil, err
	}

	if formrespuestaResponse["Data"] != nil {
		var respuestas []map[string]interface{}
		for _, r := range formrespuestaResponse["Data"].([]interface{}) {
			respuestaMap := r.(map[string]interface{})
			respuestas = append(respuestas, respuestaMap)
		}
		return respuestas, nil
	}
	return nil, nil
}

func obtenerDocumentosDesdeRespuestas(respuestas []map[string]interface{}, formularioIds []string) ([]string, error) {
	var documentos []string

	for _, respuestaMap := range respuestas {
		plantillaRespMap, ok := respuestaMap["PlantillaId"].(map[string]interface{})
		if !ok || plantillaRespMap["Id"] == nil {
			continue
		}

		formularioIdMap, ok := respuestaMap["FormularioId"].(map[string]interface{})
		if !ok || !contains(formularioIds, fmt.Sprintf("%v", formularioIdMap["Id"])) {
			continue
		}

		respuestaIdMap, ok := respuestaMap["RespuestaId"].(map[string]interface{})
		if !ok {
			continue
		}
		respuestaId := fmt.Sprintf("%v", respuestaIdMap["Id"])

		var respuestaDetalleResponse map[string]interface{}
		url := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + fmt.Sprintf("/respuesta?query=Id:%v&limit=0", respuestaId)
		err := request.GetJson(url, &respuestaDetalleResponse)
		if err != nil || respuestaDetalleResponse["Data"] == nil {
			continue
		}

		respuestasDetalle := respuestaDetalleResponse["Data"].([]interface{})
		if len(respuestasDetalle) == 0 {
			continue
		}
		respuestaDetalle := respuestasDetalle[0].(map[string]interface{})

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

	return documentos, nil
}


/*func CrearFormulario(data []byte) (APIResponseDTO requestresponse.APIResponse) {
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
				errResSec := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/seccion/", "POST", &newSec, nuevaSec)
				if errResSec != nil {
					revertir = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una de las secciones")
					return APIResponseDTO
				}
				seccionID := newSec["Data"].(map[string]interface{})["Id"].(float64)

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
							errResItem := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item/", "POST", &newItem, nuevoItem)
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
							errResItemCampo := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item_campo/", "POST", &newItemCampo, nuevoItemCampo)
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
							errResPlantilla := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/plantilla/", "POST", &newPlantilla, nuevaPlantilla)
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

		if len(plantillaIDs) > 0 {
			for _, id := range plantillaIDs {
				var plantilla map[string]interface{}

				errPlantilla := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla?query=Id:%v&Activo:true&sortby=Id&order=asc&limit=0", id), &plantilla)
				if errPlantilla == nil {

					plantillaData := plantilla["Data"].(map[string]interface{})
					plantillaData["Activo"] = false

					var inactivaPlantilla map[string]interface{}
					errPut := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla/%v", id), "PUT", &inactivaPlantilla, plantillaData)
					if errPut != nil {
						return helpers.ErrEmiter(errPut, fmt.Sprintf("Error actualizando plantilla con ID %v: %v", id, errPut))
					}
				}
			}
		}
	}

	return requestresponse.APIResponseDTO(true, 200, dataSource, "Se ha registrado el formulario c:")
}*/

func CrearFormulario(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	var itemIDs, plantillaIDs []float64

	if err := json.Unmarshal(data, &dataSource); err != nil {
		return helpers.ErrEmiter(err, "error al deserializar los datos")
	}

	secciones, ok := dataSource["secciones"].([]interface{})
	if !ok {
		return helpers.ErrEmiter(errors.New("formato inválido"), "Secciones no encontradas")
	}

	for _, seccion := range secciones {
		secMap, ok := seccion.(map[string]interface{})
		if !ok {
			continue
		}

		seccionID, err := crearSeccion(secMap)
		if err != nil {
			rollbackPlantillas(plantillaIDs)
			return helpers.ErrEmiter(err, "Error al guardar una de las secciones")
		}

		items, ok := secMap["items"].([]interface{})
		if ok {
			for _, item := range items {
				itemID, plantillaID, err := procesarItem(item, seccionID, dataSource)
				if err != nil {
					rollbackPlantillas(plantillaIDs)
					return helpers.ErrEmiter(err, "Error al guardar un item o sus relaciones")
				}
				itemIDs = append(itemIDs, itemID)
				plantillaIDs = append(plantillaIDs, plantillaID)
			}
		}
	}

	return requestresponse.APIResponseDTO(true, 200, dataSource, "Se ha registrado el formulario correctamente")
}

func crearSeccion(secMap map[string]interface{}) (float64, error) {
	nuevaSec := map[string]interface{}{
		"Activo": true,
		"Nombre": secMap["nombre"],
		"Orden":  secMap["orden"],
	}
	var newSec map[string]interface{}
	err := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/seccion/", "POST", &newSec, nuevaSec)
	if err != nil {
		return 0, err
	}
	return newSec["Data"].(map[string]interface{})["Id"].(float64), nil
}

func procesarItem(item interface{}, seccionID float64, dataSource map[string]interface{}) (float64, float64, error) {
	itemMap := item.(map[string]interface{})
	nuevoItem := map[string]interface{}{
		"Activo":     true,
		"Nombre":     itemMap["nombre"],
		"Orden":      itemMap["orden"],
		"Porcentaje": itemMap["porcentaje"],
	}

	var newItem map[string]interface{}
	err := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item/", "POST", &newItem, nuevoItem)
	if err != nil {
		return 0, 0, err
	}

	itemID := newItem["Data"].(map[string]interface{})["Id"].(float64)

	// Crear item_campo
	nuevoItemCampo := map[string]interface{}{
		"Activo":     true,
		"CampoId":    map[string]interface{}{"Id": itemMap["campo_id"]},
		"ItemId":     map[string]interface{}{"Id": itemID},
		"Porcentaje": itemMap["porcentaje"],
	}
	var newItemCampo map[string]interface{}
	err = request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item_campo/", "POST", &newItemCampo, nuevoItemCampo)
	if err != nil {
		return itemID, 0, err
	}

	// Crear plantilla
	nuevaPlantilla := map[string]interface{}{
		"Activo":       true,
		"SeccionId":    map[string]interface{}{"Id": seccionID},
		"ItemId":       map[string]interface{}{"Id": itemID},
		"ProcesoId":    dataSource["proceso_id"],
		"EstructuraId": dataSource["estructura"],
	}
	var newPlantilla map[string]interface{}
	err = request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/plantilla/", "POST", &newPlantilla, nuevaPlantilla)
	if err != nil {
		return itemID, 0, err
	}

	plantillaID := newPlantilla["Data"].(map[string]interface{})["Id"].(float64)

	return itemID, plantillaID, nil
}

func rollbackPlantillas(plantillaIDs []float64) {
	for _, id := range plantillaIDs {
		var plantilla map[string]interface{}
		url := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + fmt.Sprintf("/plantilla?query=Id:%v&Activo=true&sortby=Id&order=asc&limit=0", id)
		err := request.GetJson(url, &plantilla)
		if err != nil {
			continue
		}

		plantillaData := plantilla["Data"].(map[string]interface{})
		plantillaData["Activo"] = false

		var inactivaPlantilla map[string]interface{}
		request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla/%v", id), "PUT", &inactivaPlantilla, plantillaData)
	}
}


/*func FormularioCoevaluacion(id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {
	var formularioID int
	var plantilla map[string]interface{}
	errPlantilla := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=ProcesoId:5&Activo:true&sortby=Id&order=asc&limit=0"), &plantilla)
	if errPlantilla != nil || fmt.Sprintf("%v", plantilla) == EmptyMapString {
		return helpers.ErrEmiter(errPlantilla, fmt.Sprintf("%v", plantilla))
	}

	var itemCampos map[string]interface{}
	errItemCampos := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &itemCampos)
	if errItemCampos != nil || fmt.Sprintf("%v", itemCampos) == EmptyMapString {
		return helpers.ErrEmiter(errItemCampos, fmt.Sprintf("%v", itemCampos))
	}

	var campos map[string]interface{}
	errCampos := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &campos)
	if errCampos != nil || fmt.Sprintf("%v", campos) == EmptyMapString {
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
	errFormulario := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v,EspacioAcademicoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_periodo, id_tercero, id_espacio), &res)

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
}*/

func FormularioCoevaluacion(idPeriodo, idTercero, idEspacio string) (APIResponseDTO requestresponse.APIResponse) {
	formularioID := verificarFormularioExistenteDos(idPeriodo, idTercero, idEspacio)

	plantillaData, itemCamposData, camposData, err := obtenerDatosPlantilla()
	if err != nil {
		return helpers.ErrEmiter(err, "Error al obtener datos de plantilla")
	}

	itemCamposMap := mapearItemCampos(itemCamposData, camposData, idTercero, idEspacio)

	secciones := construirSecciones(plantillaData, itemCamposMap, formularioID)
	response := map[string]interface{}{
		"docente":          idTercero,
		"espacioAcademico": idEspacio,
		"seccion":          secciones,
	}
	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}

func obtenerDatosPlantilla() (plantillaData, itemCamposData, camposData []interface{}, err error) {
	var plantilla, itemCampos, campos map[string]interface{}

	queries := []struct {
		url  string
		dest *map[string]interface{}
	}{
		{fmt.Sprintf("plantilla?query=ProcesoId:5&Activo:true&sortby=Id&order=asc&limit=0"), &plantilla},
		{fmt.Sprintf("item_campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &itemCampos},
		{fmt.Sprintf("campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &campos},
	}

	for _, q := range queries {
		fullURL := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + q.url
		if e := request.GetJson(fullURL, q.dest); e != nil || fmt.Sprintf("%v", *q.dest) == EmptyMapString {
			return nil, nil, nil, e
		}
	}

	return plantilla["Data"].([]interface{}), itemCampos["Data"].([]interface{}), campos["Data"].([]interface{}), nil
}

func mapearItemCampos(itemCamposData, camposData []interface{}, idTercero, idEspacio string) map[int][]map[string]interface{} {
	itemCamposMap := make(map[int][]map[string]interface{})

	for _, ic := range itemCamposData {
		icMap := ic.(map[string]interface{})
		campo := icMap["CampoId"].(map[string]interface{})
		campoId := int(campo["Id"].(float64))
		tipoCampo := int(campo["TipoCampoId"].(float64))
		itemId := int(icMap["ItemId"].(map[string]interface{})["Id"].(float64))

		campoInfo := map[string]interface{}{
			"nombre":     campo["Nombre"].(string),
			"campo_id":   campoId,
			"tipo_campo": tipoCampo,
			"valor":      campo["Valor"],
		}

		if tipoCampo == 6 {
			itemRel := int(icMap["Porcentaje"].(float64))
			descarga := obtenerDescargaArchivos(idTercero, idEspacio, strconv.Itoa(itemRel))
			for k, v := range descarga {
				campoInfo[k] = v
			}
			campoInfo["nombre"] = "descarga_archivos"
			campoInfo["tipo_campo"] = 4672
			itemCamposMap[itemId] = append(itemCamposMap[itemRel], campoInfo)
		} else {
			campoInfo["porcentaje"] = icMap["Porcentaje"]
			campoInfo["escala"] = obtenerCamposHijos(campoId, camposData)
			itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
		}
	}
	return itemCamposMap
}

func verificarFormularioExistenteDos(id_periodo, idTercero, idEspacio string) int {
	var res map[string]interface{}
	url := fmt.Sprintf("formulario?query=PeriodoId:%v,EvaluadoId:%v,EspacioAcademicoId:%v&sortby=Id&order=asc&limit=0&Activo=true", id_periodo, idTercero, idEspacio)
	err := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+url, &res)
	if err != nil || res["Data"] == nil {
		return 0
	}
	if data, ok := res["Data"].([]interface{}); ok && len(data) > 0 {
		if formulario, ok := data[0].(map[string]interface{}); ok {
			if id, exists := formulario["Id"].(float64); exists {
				return int(id)
			}
		}
	}
	return 0
}

func construirSecciones(data []interface{}, itemCamposMap map[int][]map[string]interface{}, formularioID int) []map[string]interface{} {
	secciones := []map[string]interface{}{}

	for _, item := range data {
		itemMap := item.(map[string]interface{})
		seccion := itemMap["SeccionId"].(map[string]interface{})
		seccionId := int(seccion["Id"].(float64))

		seccionEncontrada := buscarSeccion(secciones, seccionId)
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

		if formularioID > 0 && VerificarRespuesta(formularioID, itemId).Status == 200 {
			continue // ya existe respuesta
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

	return secciones
}

func buscarSeccion(secciones []map[string]interface{}, id int) map[string]interface{} {
	for _, sec := range secciones {
		if sec["id"].(int) == id {
			return sec
		}
	}
	return nil
}


/*func CrearFormularioCo(data []byte) (APIResponseDTO requestresponse.APIResponse) {
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
				errResSec := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/seccion/", "POST", &newSec, nuevaSec)
				if errResSec != nil {
					revertir = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar una de las secciones")
					return APIResponseDTO
				}
				seccionID := newSec["Data"].(map[string]interface{})["Id"].(float64)

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
								itemID = idExistente
							} else {
								nuevoItem := map[string]interface{}{
									"Activo": true,
									"Nombre": nombreItem,
									"Orden":  ordenItem,
								}
								var newItem map[string]interface{}
								errResItem := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item/", "POST", &newItem, nuevoItem)
								if errResItem != nil {
									revertir = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 500, nil, "Error al guardar uno de los items")
									return APIResponseDTO
								}
								itemID = newItem["Data"].(map[string]interface{})["Id"].(float64)
								itemIDs = append(itemIDs, itemID)
								nombreItemMap[nombreItem] = itemID
							}
							nuevoItemCampo := map[string]interface{}{
								"Activo":     true,
								"CampoId":    map[string]interface{}{"Id": campoID},
								"ItemId":     map[string]interface{}{"Id": itemID},
								"Porcentaje": porcentaje,
							}
							var newItemCampo map[string]interface{}
							errResItemCampo := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item_campo/", "POST", &newItemCampo, nuevoItemCampo)
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
							errResPlantilla := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/plantilla/", "POST", &newPlantilla, nuevaPlantilla)
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

		if len(plantillaIDs) > 0 {
			for _, id := range plantillaIDs {
				var plantilla map[string]interface{}

				errPlantilla := request.GetJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla?query=Id:%v&Activo:true&sortby=Id&order=asc&limit=0", id), &plantilla)
				if errPlantilla == nil {

					plantillaData := plantilla["Data"].(map[string]interface{})
					plantillaData["Activo"] = false

					var inactivaPlantilla map[string]interface{}
					errPut := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla/%v", id), "PUT", &inactivaPlantilla, plantillaData)
					if errPut != nil {
						return helpers.ErrEmiter(errPut, fmt.Sprintf("Error actualizando plantilla con ID %v: %v", id, errPut))
					}
				}
			}
		}
	}

	return requestresponse.APIResponseDTO(true, 200, dataSource, "Se ha registrado el formulario c:")
}*/

func CrearFormularioCo(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}
	if err := json.Unmarshal(data, &dataSource); err != nil {
		return helpers.ErrEmiter(err, "error al deserializar los datos")
	}

	secciones, _ := dataSource["secciones"].([]interface{})
	nombreItemMap := make(map[string]float64)
	var plantillaIDs []float64

	plantillaIDs, err := crearSeccionesConItems(secciones, dataSource, nombreItemMap)
	if err != nil {
		revertirPlantillas(plantillaIDs)
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, dataSource, "Se ha registrado el formulario c:")
}

func crearSeccionesConItems(secciones []interface{}, dataSource map[string]interface{}, nombreItemMap map[string]float64) ([]float64, error) {
	var plantillaIDs []float64

	for _, seccion := range secciones {
		secMap := seccion.(map[string]interface{})
		nuevaSec := map[string]interface{}{
			"Activo": true,
			"Nombre": secMap["nombre"],
			"Orden":  secMap["orden"],
		}
		var newSec map[string]interface{}
		err := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/seccion/", "POST", &newSec, nuevaSec)
		if err != nil {
			return plantillaIDs, fmt.Errorf("Error al guardar una de las secciones")
		}
		seccionID := newSec["Data"].(map[string]interface{})["Id"].(float64)

		items, _ := secMap["items"].([]interface{})
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			nombreItem := itemMap["nombre"].(string)
			ordenItem := itemMap["orden"]
			campoID := itemMap["campo_id"].(float64)
			porcentaje := itemMap["porcentaje"]
			if campoID == 4672 {
				porcentaje = itemMap["item_relacion_id"].(float64)
			}

			itemID, err := crearItemSiNoExiste(nombreItem, ordenItem, nombreItemMap)
			if err != nil {
				return plantillaIDs, err
			}

			err = crearItemCampo(campoID, itemID, porcentaje)
			if err != nil {
				return plantillaIDs, err
			}

			plantillaID, err := crearPlantilla(seccionID, itemID, dataSource)
			if err != nil {
				return plantillaIDs, err
			}
			plantillaIDs = append(plantillaIDs, plantillaID)
		}
	}

	return plantillaIDs, nil
}

func crearItemSiNoExiste(nombre string, orden interface{}, nombreItemMap map[string]float64) (float64, error) {
	if id, existe := nombreItemMap[nombre]; existe {
		return id, nil
	}
	nuevoItem := map[string]interface{}{
		"Activo": true,
		"Nombre": nombre,
		"Orden":  orden,
	}
	var newItem map[string]interface{}
	err := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item/", "POST", &newItem, nuevoItem)
	if err != nil {
		return 0, fmt.Errorf("Error al guardar uno de los items")
	}
	id := newItem["Data"].(map[string]interface{})["Id"].(float64)
	nombreItemMap[nombre] = id
	return id, nil
}

func crearItemCampo(campoID float64, itemID float64, porcentaje interface{}) error {
	nuevoItemCampo := map[string]interface{}{
		"Activo":     true,
		"CampoId":    map[string]interface{}{"Id": campoID},
		"ItemId":     map[string]interface{}{"Id": itemID},
		"Porcentaje": porcentaje,
	}
	var newItemCampo map[string]interface{}
	err := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/item_campo/", "POST", &newItemCampo, nuevoItemCampo)
	if err != nil {
		return fmt.Errorf("Error al guardar uno de los items_campo")
	}
	return nil
}

func crearPlantilla(seccionID, itemID float64, dataSource map[string]interface{}) (float64, error) {
	nuevaPlantilla := map[string]interface{}{
		"Activo":       true,
		"SeccionId":    map[string]interface{}{"Id": seccionID},
		"ItemId":       map[string]interface{}{"Id": itemID},
		"ProcesoId":    dataSource["proceso_id"],
		"EstructuraId": dataSource["estructura"],
	}
	var newPlantilla map[string]interface{}
	err := request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+"/plantilla/", "POST", &newPlantilla, nuevaPlantilla)
	if err != nil {
		return 0, fmt.Errorf("Error al guardar una plantilla")
	}
	return newPlantilla["Data"].(map[string]interface{})["Id"].(float64), nil
}

func revertirPlantillas(plantillaIDs []float64) {
	for _, id := range plantillaIDs {
		var plantilla map[string]interface{}
		url := HttpPrefix + beego.AppConfig.String("EvaluacionDocenteService") + fmt.Sprintf("/plantilla?query=Id:%v&Activo:true&sortby=Id&order=asc&limit=0", id)
		err := request.GetJson(url, &plantilla)
		if err != nil {
			continue
		}
		plantillaData := plantilla["Data"].(map[string]interface{})
		plantillaData["Activo"] = false

		var inactiva map[string]interface{}
		_ = request.SendJson(HttpPrefix+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("/plantilla/%v", id), "PUT", &inactiva, plantillaData)
	}
}

