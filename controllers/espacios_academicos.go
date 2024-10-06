package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// EspacioAcademicoController operations for EspacioAcademico
type EspacioAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *EspacioAcademicoController) URLMapping() {
	c.Mapping("Post", c.PostEspacioAcademico)
}

// PostEspacioAcademico ...
// @Title PostEspacioAcademico
// @Description query EspacioAcademico
// @Param	body		body 	models.Parametros	true		"body for Parametros content"
// @Success 200 {object} models.Parametros
// @Failure 404 body is empty
// @router / [post]
func (c *EspacioAcademicoController) PostEspacioAcademico() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		respuesta := services.ConsultarEspacios(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
	}
}
