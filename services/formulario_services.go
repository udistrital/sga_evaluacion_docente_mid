package services

import (
	"fmt"
	"sort"

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
		campoInfo := map[string]interface{}{
			"nombre":     campo["Nombre"].(string),
			"tipo_campo": int(campo["TipoCampoId"].(float64)),
			"valor":      campo["Valor"],
			"porcentaje": itemCampoMap["Porcentaje"],
			"escala":     obtenerCamposHijos(campoId, camposData),
		}
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
				}
				hijos = append(hijos, hijoInfo)
			}
		}
	}

	return hijos
}
