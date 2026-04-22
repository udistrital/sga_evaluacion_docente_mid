package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// DocumentoController operations for Documentos
type DocumentoController struct {
	beego.Controller
}

// URLMapping ...
func (c *DocumentoController) URLMapping() {
	c.Mapping("GetConcatDocuments", c.GetConcatDocuments)
}

// GetConcatDocuments ...
// @Title GetConcatDocuments
// @Description Obtiene los documentos concatenados en base64 filtrados por evaluado_id y periodo_id
// @Param	periodo		path 	string	true		"Id del periodo (de parametros)"
// @Param	evaluado		path 	string	true		"id del evaluado"
// @Success 200 se obtiene el documento concatenado en base64
// @Success 204 con los parametros solicitados no se encontraron documentos
// @Failure 403 invalid path
// @Failure 400 error en los parametros
// @router /documentos_concatenados/:periodo/:evaluado [get]
func (c *DocumentoController) GetConcatDocuments() {
	defer errorhandler.HandlePanic(&c.Controller)
	perID, err1 := strconv.Atoi(c.Ctx.Input.Param(":periodo"))
	evalID, err2 := strconv.Atoi(c.Ctx.Input.Param(":evaluado"))
	if err1 != nil || err2 != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "Parámetros inválidos: se requieren evaluador_id, periodo_id",
		}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	} else {
		respuesta := services.ConsultarDocumentos(perID, evalID)
		c.Ctx.Output.SetStatus(respuesta.Status)
		c.Data["json"] = respuesta
		c.ServeJSON()
	}
}
