package services

import (
	"fmt"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func EnviarEmail(to []string, source string, template string, templateData models.EmailTemplateData, attachments []models.EmailAttachment) requestresponse.APIResponse {
	payload := models.EnvioEmailPayload{
		Source:   source,
		Template: template,
		Destinations: []models.EmailDestination{
			{
				Destination:             models.Destination{ToAddresses: to},
				ReplacementTemplateData: templateData,
				Attachments:             attachments,
			},
		},
		DefaultTemplateData: templateData,
	}

	emailService := beego.AppConfig.String("NotificacionService")
	url := emailService + "email/enviar_templated_email"

	var response map[string]interface{}

	err := request.SendJsonEscapeUnicode(url, "POST", &response, payload)
	if err != nil {
		beego.Error("Error en request.SendJsonEscapeUnicode:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al enviar el correo: %v", err))
	}

	success, ok := response["Success"].(bool)
	if ok && success {
		status := fmt.Sprintf("%v", response["Status"])
		if status == "200" {
			data := response["Data"]
			return requestresponse.APIResponseDTO(true, 200, data, "Correo enviado correctamente.")
		}
		return requestresponse.APIResponseDTO(false, 400, response, "El servidor respondió pero no fue exitoso.")
	}

	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	beego.Error("Respuesta inválida del servidor:", string(responseJSON))

	return requestresponse.APIResponseDTO(false, 500, response, "Respuesta inválida del servidor al intentar enviar el correo.")
}
