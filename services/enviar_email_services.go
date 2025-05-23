package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"

	"encoding/json"
)

func EnviarEmail(to []string, source string, template string, templateData models.EmailTemplateData, attachments []models.EmailAttachment ) requestresponse.APIResponse {
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
/*package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"
	"github.com/udistrital/utils_oas/requestresponse"
)

func EnviarEmail() requestresponse.APIResponse {
	// Construir payload
	payload := models.EnvioEmailPayload{
		Source:   "condor@udistrital.edu.co",
		Template: "PLANTILLA_EVALUACION_DOCENTE",
		Destinations: []models.EmailDestination{
			{
				Destination: models.Destination{
					ToAddresses: []string{"jerodrigueza@udistrital.edu.co"},
				},
				ReplacementTemplateData: models.EmailTemplateData{
					NombreUsuario:     "BARON CAMACHO LUZ AMPARO",
					DocumentoUsuario:  "12.345.678",
					NombreEvaluacion:  "Heteroevaluación",
					NumeroPeriodo:     "61",
					FechaEvaRealizada: "07 de mayo de 2025",
					HoraEvaRealizada:  "09:40 AM",
				},
				Attachments: []models.EmailAttachment{
					{
						ContentType: "application/pdf",
						FileName:    "prueba.pdf",
						Base64File:  "", // base64 si es necesario
					},
				},
			},
		},
		DefaultTemplateData: models.EmailTemplateData{
			NombreUsuario:     "BARON CAMACHO LUZ AMPARO",
			DocumentoUsuario:  "12.345.678",
			NombreEvaluacion:  "Heteroevaluación",
			NumeroPeriodo:     "61",
			FechaEvaRealizada: "07 de mayo de 2025",
			HoraEvaRealizada:  "09:40 AM",
		},
	}

	// Serializar a JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error serializando JSON: %v", err))
	}

	// URL de destino
	emailService := beego.AppConfig.String("NotificacionService")
	url := emailService + "email/enviar_templated_email"
	fmt.Println("URL destino:", url)

	// Crear cliente HTTP
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error creando request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar request
	resp, err := client.Do(req)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error al enviar request: %v", err))
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error leyendo respuesta: %v", err))
	}

	// Decodificar respuesta
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return requestresponse.APIResponseDTO(false, 500, string(body), fmt.Sprintf("Respuesta no es JSON válido: %v", err))
	}

	fmt.Printf("Respuesta del servicio:\n%+v\n", result)

	// Verificar éxito
	if success, ok := result["Success"].(bool); ok && success {
		return requestresponse.APIResponseDTO(true, 200, result["Data"], "Correo enviado correctamente.")
	}

	return requestresponse.APIResponseDTO(false, 500, result, "Respuesta inválida del servidor.")
}*/
