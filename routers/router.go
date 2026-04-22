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
		beego.NSNamespace("/espacios_academicos",
			beego.NSInclude(
				&controllers.EspacioAcademicoController{},
			),
		),
		beego.NSNamespace("/reporte_autoevaluacion_ii_tres_consejo",
			beego.NSInclude(
				&controllers.ReporteAutoevaluacionIITresConsejoController{},
			),
		),
		beego.NSNamespace("/reporte_heteroevaluacion_consejo",
			beego.NSInclude(
				&controllers.ReporteHeteroevaluacionConsejoController{},
			),
		),
		beego.NSNamespace("/reporte_coevaluacion_i_consejo",
			beego.NSInclude(
				&controllers.ReporteCoevaluacionIConsejoController{},
			),
		),
		beego.NSNamespace("/documento",
			beego.NSInclude(
				&controllers.DocumentoController{},
			),
		),
		beego.NSNamespace("/enviar_notificacion",
			beego.NSInclude(
				&controllers.EnviarEmailController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
