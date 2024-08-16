package services

import (
	"fmt"

	"github.com/beego/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

type FormularioService struct {
	// Aquí puedes definir los campos y métodos necesarios para tu servicio
}

func ConsultaFormulario(id_tipo_formulario string, id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {
	fmt.Println("Entró")

	var secciones []map[string]interface{}

	var plantilla []map[string]interface{}
	errPlantilla := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("plantilla?query=EstructuraId:%v", id_tipo_formulario), &plantilla)
	if errPlantilla != nil || fmt.Sprintf("%v", plantilla) == "[map[]]" {
		return helpers.ErrEmiter(errPlantilla, fmt.Sprintf("%v", plantilla))
	}

	for _, p := range plantilla {

		var seccion map[string]interface{}
		errSeccion := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("seccion/%v", p["SeccionId"]), &seccion)
		if errSeccion != nil || fmt.Sprintf("%v", seccion) == "[map[]]" {
			return helpers.ErrEmiter(errSeccion, fmt.Sprintf("%v", seccion))
		}

		var items []map[string]interface{}
		errItems := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item?query=Id:%v", p["ItemId"]), &items)
		if errItems != nil || fmt.Sprintf("%v", items) == "[map[]]" {
			return helpers.ErrEmiter(errItems, fmt.Sprintf("%v", items))
		}

		var itemsStruct map[string]interface{}
		for _, item := range items {

			var itemCampos []map[string]interface{}
			errItemCampos := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("item_campo?query=ItemId:%v", item["Id"]), &itemCampos)
			if errItemCampos != nil || fmt.Sprintf("%v", itemCampos) == "[map[]]" {
				return helpers.ErrEmiter(errItemCampos, fmt.Sprintf("%v", itemCampos))
			}

			var camposStruct []map[string]interface{}
			for _, campo := range itemCampos {

				var campoDetalle map[string]interface{}
				errCampoDetalle := request.GetJson("http://"+beego.AppConfig.String("EvaluacionDocenteService")+fmt.Sprintf("campo/%v", campo["CampoId"]), &campoDetalle)
				if errCampoDetalle != nil || fmt.Sprintf("%v", campoDetalle) == "[map[]]" {
					return helpers.ErrEmiter(errCampoDetalle, fmt.Sprintf("%v", campoDetalle))
				}

				camposStruct = append(camposStruct, map[string]interface{}{
					"tipoCampo": campoDetalle["TipoCampoId"],
					"nombre":    campoDetalle["Nombre"],
					"valor":     campoDetalle["Valor"],
				})
			}

			itemsStruct[item["Nombre"].(string)] = map[string]interface{}{
				"orden":    item["Orden"],
				"campos":   camposStruct,
				"tipoItem": item["TipoItem"],
			}
		}

		secciones = append(secciones, map[string]interface{}{
			"orden":         seccion["Orden"],
			"nombreSeccion": seccion["Nombre"],
			"items":         itemsStruct,
		})
	}

	response := map[string]interface{}{
		"tipoEvaluacion":   "Evaluación de ejemplo",
		"docente":          "Nombre del docente",
		"espacioAcademico": "Espacio académico",
		"seccion":          secciones,
	}

	return requestresponse.APIResponseDTO(true, 200, response, "Consulta exitosa")
}
