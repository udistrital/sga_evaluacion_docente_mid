package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func ConsultaFormulario(id_tipo_formulario string, id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {
	fmt.Println("Entró")

	var plantilla map[string]interface{}
	errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=ProcesoId:%v&Activo:true&sortby=Id&order=asc&limit=0", id_tipo_formulario), &plantilla)
	if errPlantilla != nil || fmt.Sprintf("%v", plantilla) == "[map[]]" {
		fmt.Println(plantilla)
		return helpers.ErrEmiter(errPlantilla, fmt.Sprintf("%v", plantilla))
	}

	var itemCampos map[string]interface{}
	errItemCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &itemCampos)
	if errItemCampos != nil || fmt.Sprintf("%v", itemCampos) == "[map[]]" {
		fmt.Println(itemCampos)
		return helpers.ErrEmiter(errItemCampos, fmt.Sprintf("%v", itemCampos))
	}

	var campos map[string]interface{}
	errCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo?query=Activo:true&sortby=Id&order=asc&limit=0"), &campos)
	if errCampos != nil || fmt.Sprintf("%v", campos) == "[map[]]" {
		fmt.Println(campos)
		return helpers.ErrEmiter(errCampos, fmt.Sprintf("%v", campos))
	}

	secciones := map[int]map[string]interface{}{}
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
			"escala":     obtenerCamposHijos(campoId, camposData), //funcion para encontrar los campos hijos si los hay si no regresa lo mismo
		}
		itemCamposMap[itemId] = append(itemCamposMap[itemId], campoInfo)
	}

	for _, item := range data {
		itemMap := item.(map[string]interface{})
		seccion := itemMap["SeccionId"].(map[string]interface{})
		seccionId := int(seccion["Id"].(float64))

		if _, exists := secciones[seccionId]; !exists {
			secciones[seccionId] = map[string]interface{}{
				"nombre": seccion["Nombre"].(string),
				"orden":  int(seccion["Orden"].(float64)),
				"items":  []map[string]interface{}{},
			}
		}

		itemId := int(itemMap["ItemId"].(map[string]interface{})["Id"].(float64))
		itemInfo := map[string]interface{}{
			"id":     itemId,
			"nombre": itemMap["ItemId"].(map[string]interface{})["Nombre"].(string),
			"orden":  int(itemMap["ItemId"].(map[string]interface{})["Orden"].(float64)),
			"campos": itemCamposMap[itemId],
		}
		secciones[seccionId]["items"] = append(secciones[seccionId]["items"].([]map[string]interface{}), itemInfo)
	}

	ordenSecciones := map[string]interface{}{}
	ordenNombres := map[int]string{
		0: "descripcion",
		1: "cuantitativa",
		2: "cualitativa",
		3: "carga de archivos",
		4: "descarga de archivos",
	}

	for _, seccion := range secciones {
		seccionMap := seccion
		orden := seccionMap["orden"].(int)
		nombre := ordenNombres[orden]
		if nombre == "" {
			nombre = fmt.Sprintf("seccion_%d", orden)
		}
		ordenSecciones[nombre] = seccionMap
	}

	//Aqui queadaría organizada por secciones
	response := map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Consulta exitosa",
		"Data": map[string]interface{}{
			"docente":          id_tercero, //consultar cuando exista la data del tercero evaluado
			"espacioAcademico": id_espacio, //consultar cuando exista la data del espacio academico
			"seccion":          ordenSecciones,
			"tipoEvaluacion":   id_tipo_formulario, //consultar a parametro  y se le pasa el id del tipo de evaluacion
		},
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
