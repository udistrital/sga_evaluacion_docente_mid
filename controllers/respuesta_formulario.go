package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Respuesta_formularioController operations for Respuesta_formulario
type Respuesta_formularioController struct {
	beego.Controller
}

// URLMapping ...
func (c *Respuesta_formularioController) URLMapping() {
	c.Mapping("PostRespuestaFormulario", c.PostRespuestaFormulario)
}

// PostRespuestaFormulario ...
// @Title PostRespuestaFormulario
// @Description create Respuesta_formulario
// @Param	body		body 	models.Respuesta_formulario	true		"body for Respuesta_formulario content"
// @Success 201 {object} models.Respuesta_formulario
// @Failure 403 body is empty
// @router / [post]
func (c *Respuesta_formularioController) PostRespuestaFormulario() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		respuesta := services.GuardarRespuestas(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
	}
}
