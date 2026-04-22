package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
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
		log.Println("Error al deserializar el JSON:", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "JSON mal formado")
		c.ServeJSON()
		return
	}

	nombrePeriodo, err := obtenerNombrePeriodo(eval.Periodo)
	if err != nil {
		log.Printf("Error al obtener el nombre del periodo: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error obteniendo periodo: %v", err.Error()))
		c.ServeJSON()
		return
	}

	emailData, err := construirEmailData(eval, nombrePeriodo)
	if err != nil {
		log.Println(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		c.ServeJSON()
		return
	}

	b64PDF, err := helpers.CrearCertificadoPdfEvaluacion(
		emailData.NombreUsuario,
		emailData.DocumentoUsuario,
		eval.DatosEmail.EspacioAcademico,
		emailData.FechaEvaRealizada,
		emailData.NombreEvaluacion,
		emailData.NumeroPeriodo,
		eval.DatosEmail.Docente,
		eval.DatosEmail.Grupo,
	)
	if err != nil {
		log.Printf("Error creando el certificado PDF: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 500, nil, "Error generando el certificado PDF")
		c.ServeJSON()
		return
	}

	adjuntos := []models.EmailAttachment{
		{
			ContentType: "application/pdf",
			FileName:    "Certificado evaluación docente.pdf",
			Base64File:  b64PDF,
		},
	}

	respuesta := services.EnviarEmail(
		[]string{eval.Correo},
		"condor@udistrital.edu.co",
		"PLANTILLA_EVALUACION_DOCENTE",
		emailData,
		adjuntos,
	)

	if respuesta.Success && fmt.Sprintf("%v", respuesta.Status) == "200" {
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, respuesta.Data, fmt.Sprintf("%v", respuesta.Message))
	} else {
		log.Printf("Error al enviar el correo. Status: %v, Mensaje: %v", respuesta.Status, respuesta.Message)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error enviando el correo: %v", respuesta.Message))
	}

	c.ServeJSON()
}

func obtenerNombrePeriodo(periodoID int) (string, error) {
	var resp map[string]interface{}
	url := beego.AppConfig.String("ParametrosService") + "/periodo/" + fmt.Sprintf("%d", periodoID)

	if err := request.GetJson(url, &resp); err != nil {
		return "", fmt.Errorf("no se pudo obtener el periodo desde el servicio: %w", err)
	}

	data, ok := resp["Data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("respuesta del servicio inválida o sin campo 'Data'")
	}

	nombre, ok := data["Nombre"].(string)
	if !ok {
		return "", fmt.Errorf("el campo 'Nombre' no existe o no es string")
	}

	return nombre, nil
}

func construirEmailData(eval models.DatosEmailRequest, nombrePeriodo string) (models.EmailTemplateData, error) {
	data := models.EmailTemplateData{
		NombreEvaluacion:  eval.NombreEvaluacion,
		NumeroPeriodo:     nombrePeriodo,
		FechaEvaRealizada: eval.Fecha,
		HoraEvaRealizada:  eval.Hora,
	}

	switch eval.NombreEvaluacion {
	case "Autoevaluación I", "Heteroevaluación":
		data.NombreUsuario = eval.DatosEmail.Estudiante
		data.DocumentoUsuario = eval.Documento

	case "Coevaluación I", "Autoevaluación II 1", "Autoevaluación II 2", "Autoevaluación II 3":
		data.NombreUsuario = eval.DatosEmail.Docente
		data.DocumentoUsuario = eval.Documento

	case "Coevaluación II":
		data.NombreUsuario = eval.DatosEmail.ConsejoCurricular
		data.DocumentoUsuario = eval.Documento

	default:
		return data, fmt.Errorf("Tipo de evaluación no reconocido: %s", eval.NombreEvaluacion)
	}

	return data, nil
}
