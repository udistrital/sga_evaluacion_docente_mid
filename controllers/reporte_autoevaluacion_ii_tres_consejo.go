package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
)

// ReporteAutoevaluacionIITresConsejoController operations for ReporteAutoevaluacionIITresConsejo
type ReporteAutoevaluacionIITresConsejoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteAutoevaluacionIITresConsejoController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// @router / [get]
// @Title GetReporteAutoevaluacionIITresConsejo
// @Description Obtiene promedio y respuestas abiertas de autoevaluación II (tipo consejo) por evaluador
// @Param   evaluador_id	query	string	true	"ID del evaluador (documento)"
// @Param   periodo_id  	query	int	true	"ID del periodo académico"
// @Param   proceso_id  	query	int	true	"ID del proceso de evaluación"
// @Param   nombre_evaluador	query	string	true	"Nombre del evaluador"
// @Success 200 {object} interface{} "Resultado del servicio de reporte"
// @Failure 400 el evaluador_id o ids inválidos
// @Failure 500 error interno al ejecutar la consulta
func (c *ReporteAutoevaluacionIITresConsejoController) Get() {
	evaluadorId := c.GetString("evaluador_id")
	periodoId, err1 := c.GetInt("periodo_id")
	procesoId, err2 := c.GetInt("proceso_id")
	nombreEvaluador := c.GetString("nombre_evaluador")

	if err1 != nil || err2 != nil || evaluadorId == "" || nombreEvaluador == "" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "Parámetros inválidos: se requieren evaluador_id, periodo_id, proceso_id y nombre_evaluador",
		}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	resultado := services.GetReporteAutoevaluacionIITresConsejo(evaluadorId, periodoId, procesoId, nombreEvaluador)
	c.Data["json"] = resultado
	c.ServeJSON()
}
