package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	//"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"
	"github.com/udistrital/utils_oas/requestresponse"
	"fmt"
	//"strings"
	"log"
	"encoding/json"
)

// EnviarEmailController operations for EnviarEmail
type EnviarEmailController struct {
	beego.Controller
}

// URLMapping ...
func (c *EnviarEmailController) URLMapping() {
	c.Mapping("PostEnviarEmail", c.PostEnviarEmail)
}

// PostEnviarEmail ...
// @Title PostEnviarEmail
// @Description Recibe los datos para enviar el email
// @Param	body		body 	models.DatosEmailRequest	true		"body for PostEnviarEmail content"
// @Success 201 {object} models.DatosEmailRequest
// @Failure 400 body is empty or bad request
// @router / [post]
func (c *EnviarEmailController) PostEnviarEmail() {
	defer errorhandler.HandlePanic(&c.Controller)

	var eval models.DatosEmailRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &eval); err != nil {
		log.Println("Error unmarshalling body:", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "JSON mal formado")
		c.ServeJSON()
		return
	}

	data := models.EmailTemplateData{
		NombreEvaluacion:  eval.NombreEvaluacion,
		NumeroPeriodo:     fmt.Sprintf("%d", eval.Periodo),
		FechaEvaRealizada: eval.Fecha,
		HoraEvaRealizada:  eval.Hora,
	}

	switch eval.NombreEvaluacion {
		case "Autoevaluación I", "Heteroevaluación":
			data.NombreUsuario = eval.DatosEmail.Estudiante
			data.DocumentoUsuario = eval.Documento

		case "Coevaluación I", "Autoevaluación II 3", "Autoevaluación II 2", "Autoevaluación II 1":
			data.NombreUsuario = eval.DatosEmail.Docente
			data.DocumentoUsuario = eval.Documento
			fmt.Println("ingreso a autoevaluación ii:")

		case "Coevaluación II":
			data.NombreUsuario = eval.DatosEmail.ConsejoCurricular
			data.DocumentoUsuario = eval.Documento
			fmt.Println("ingreso a coevaluación ii")

		default:
			log.Printf("Tipo de evaluación no reconocido: %s\n", eval.NombreEvaluacion)
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Tipo de evaluación no reconocido")
			c.ServeJSON()
			return
	}

	// Envío del correo
	respuesta := services.EnviarEmail(
		[]string{eval.Correo},
		"condor@udistrital.edu.co",
		"PLANTILLA_EVALUACION_DOCENTE",
		data,
		nil, // sin adjuntos
	)

	log.Println("Resultado del envío:", respuesta)

	c.Data["json"] = requestresponse.APIResponseDTO(true, 200, nil, "Correo enviado exitosamente")
	c.ServeJSON()
}
