package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// CargaAcademicaController operations for CargaAcademica
type CargaAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *CargaAcademicaController) URLMapping() {
	c.Mapping("Post", c.PostCargaAcademica)
}

// PostCargaAcademica ...
// @Title PostCargaAcademica
// @Description create CargaAcademica
// @Param	body		body 	models.CargaAcademica	true		"body for CargaAcademica content"
// @Success 201 {object} models.CargaAcademica
// @Failure 403 body is empty
// @router / [post]
func (c *CargaAcademicaController) PostCargaAcademica() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		respuesta := services.ConsultarCarga(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
	}
}
