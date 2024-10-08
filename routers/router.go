// @APIVersion 1.0.0
// @Title SGA MID - Evaluación Docente
// @Description Microservicio MID del SGA MID que complementa evaluación docente
package routers

import (
	"github.com/udistrital/sga_evaluacion_docente_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"

	"github.com/astaxie/beego"
)

func init() {
	beego.ErrorController(&errorhandler.ErrorHandlerController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/formulario_por_tipo",
			beego.NSInclude(
				&controllers.Formulario_por_tipoController{},
			),
		),
		beego.NSNamespace("/respuesta_formulario",
			beego.NSInclude(
				&controllers.Respuesta_formularioController{},
			),
		),
		beego.NSNamespace("/metricas",
			beego.NSInclude(
				&controllers.MetricasController{},
			),
		),
		beego.NSNamespace("/carga_academica",
			beego.NSInclude(
				&controllers.CargaAcademicaController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
