package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// MetricasController operations for Metricas
type MetricasController struct {
	beego.Controller
}

// URLMapping ...
func (c *MetricasController) URLMapping() {
	c.Mapping("MetricasHeteroevaluacion", c.MetricasHeteroevaluacion)
	c.Mapping("MetricasAutoevaluacion", c.MetricasAutoevaluacion)
	c.Mapping("MetricasCoevaluacion", c.MetricasCoevaluacion)
}

// MetricasHeteroevaluacion ...
// @Title MetricasHeteroevaluacion
// @Description query MetricasHeteroevaluacion
// @Param	body		body 	models.MetricasHeteroevaluacion	true		"body for MetricasHeteroevaluacion content"
// @Success 201 {object} models.MetricasHeteroevaluacion
// @Failure 403 body is empty
// @router /Heteroevaluacion [post]
func (c *MetricasController) MetricasHeteroevaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		respuesta := services.MetricasHeteroevaluacion(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
	}
}

// MetricasAutoevaluacion ...
// @Title MetricasAutoevaluacion
// @Description query MetricasAutoevaluacion
// @Param	body		body 	models.MetricasAutoevaluacion	true		"body for MetricasAutoevaluacion content"
// @Success 201 {object} models.MetricasAutoevaluacion
// @Failure 403 body is empty
// @router /Autoevaluacion [post]
func (c *MetricasController) MetricasAutoevaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		respuesta := services.MetricasAutoevaluacion(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
	}
}

// MetricasCoevaluacion ...
// @Title MetricasCoevaluacion
// @Description query MetricasCoevaluacion
// @Param	body		body 	models.MetricasCoevaluacion	true		"body for MetricasCoevaluacion content"
// @Success 201 {object} models.MetricasCoevaluacion
// @Failure 403 body is empty
// @router /Coevaluacion [post]
func (c *MetricasController) MetricasCoevaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		respuesta := services.MetricasCoevaluacion(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
	}
}
