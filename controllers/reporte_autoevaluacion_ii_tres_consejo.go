package controllers

import (
	"fmt"
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
// @Success 200 {object} interface{} "Resultado del servicio de reporte"
// @Failure 400 el evaluador_id o ids inválidos
// @Failure 500 error interno al ejecutar la consulta
func (c *ReporteAutoevaluacionIITresConsejoController) Get() {
	fmt.Println("ingrsa al controlador")
	evaluadorId := c.GetString("evaluador_id")
	periodoId, err1 := c.GetInt("periodo_id")
	procesoId, err2 := c.GetInt("proceso_id")

	if err1 != nil || err2 != nil || evaluadorId == "" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "Parámetros inválidos: se requieren evaluador_id, periodo_id y proceso_id",
		}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}
	fmt.Println("aaaaaaaaaaaaaaaaaaaa")
	resultado := services.GetReporteAutoevaluacionIITresConsejo(evaluadorId, periodoId, procesoId)
	c.Data["json"] = resultado
	c.ServeJSON()
	fmt.Println("bbbbbbbbbbb")
	/*resultado, err := services.GetReporteAutoevaluacionIITresConsejo(evaluadorId, periodoId, procesoId)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = resultado
	c.ServeJSON()*/
}