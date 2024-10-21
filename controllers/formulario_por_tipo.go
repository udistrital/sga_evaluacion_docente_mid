package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Formulario_por_tipoController operations for Formulario_por_tipo
type Formulario_por_tipoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Formulario_por_tipoController) URLMapping() {
	c.Mapping("GetFormularioTipo", c.GetFormularioTipo)
	c.Mapping("PostFormularioTipo", c.PostFormularioTipo)
}

// GetFormularioTipo ...
// @Title GetFormularioTipo
// @Description Consultar los formularios por tipo id tercero y periodo
// @Param	id_tipo_formulario	query	string	false	"Id del tipo formulario"
// @Param	id_periodo	query	string	false	"Id del periodo"
// @Param	id_tercero	query	string	false	"Id del tercero"
// @Param	id_espacio	query	string	false	"Id del espacio"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *Formulario_por_tipoController) GetFormularioTipo() {
	defer errorhandler.HandlePanic(&c.Controller)

	id_tipo_formulario := c.GetString("id_tipo_formulario")
	id_periodo := c.GetString("id_periodo")
	id_tercero := c.GetString("id_tercero")
	id_espacio := c.GetString("id_espacio")

	var respuesta requestresponse.APIResponse

	if id_tipo_formulario == "5" {
		respuesta = services.FormularioCoevaluacion(id_periodo, id_tercero, id_espacio)
	} else {
		respuesta = services.ConsultaFormulario(id_tipo_formulario, id_periodo, id_tercero, id_espacio)
	}

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta

	c.ServeJSON()

}

// PostFormularioTipo ...
// @Title PostFormularioTipo
// @Description create PostFormularioTipo
// @Param	body		body 	models.PostFormularioTipo	true		"body for PostFormularioTipo content"
// @Success 201 {object} models.PostFormularioTipo
// @Failure 403 body is empty
// @router / [post]
func (c *Formulario_por_tipoController) PostFormularioTipo() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	if data != nil {
		var payload map[string]interface{}
		err := json.Unmarshal(data, &payload)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Error al procesar datos")
			c.ServeJSON()
			return
		}

		if procesoID, ok := payload["proceso_id"].(float64); ok && int(procesoID) == 5 {
			respuesta := services.CrearFormularioCo(data)
			c.Ctx.Output.SetStatus(respuesta.Status)
			c.Data["json"] = respuesta
			c.ServeJSON()
			return
		}

		respuesta := services.CrearFormulario(data)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()

	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erróneos")
		c.ServeJSON()
	}
}
