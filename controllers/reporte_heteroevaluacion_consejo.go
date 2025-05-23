package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/services"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"

	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"log"
	"fmt"
)

// ReporteHeteroevaluacionConsejoController operations for ReporteHeteroevaluacionConsejo
type ReporteHeteroevaluacionConsejoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteHeteroevaluacionConsejoController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// @router / [get]
// @Title Get reporte de heteroevaluación - Consejo
// @Description Obtiene los promedios por ámbito y promedio general de heteroevaluación para un evaluado
// @Param   evaluado_id	query	string	true	"ID del evaluado (ej. documento)"
// @Param   periodo_id  	query	int	true	"ID del periodo académico"
// @Param   proceso_id  	query	int	true	"ID del proceso"
// @Success 200 {object} interface{} "Resultado del servicio de reporte"
// @Failure 400 parámetros inválidos
// @Failure 500 error al ejecutar la consulta
func (c *ReporteHeteroevaluacionConsejoController) Get() {

	b64, err := helpers.CrearCertificadoAutoevaluacion(
		"OMER CALDERON",
		"12119277",
		"601 - DOCTORADO INTERINSTITUCIONAL EN EDUCACIÓN",
		"18 del mes 05 de 2025",
	)
	if err != nil {
		log.Fatalf("Error creando PDF: %v", err)
	}
	fmt.Println("Base64 del PDF:", b64)

	
	//////////

	////////


	data := models.EmailTemplateData{
		NombreUsuario:     "BARON CAMACHO LUZ AMPARO",
		DocumentoUsuario:  "12.345.678",
		NombreEvaluacion:  "Heteroevaluación",
		NumeroPeriodo:     "61",
		FechaEvaRealizada: "07 de mayo de 2025",
		HoraEvaRealizada:  "09:40 AM",
	}

	adjuntos := []models.EmailAttachment{
		{
			ContentType: "application/pdf",
			FileName:    "certificado.pdf",
			Base64File:  b64, 
		},
	}

	respuesta := services.EnviarEmail(
		[]string{"jerodrigueza@udistrital.edu.co"},
		"condor@udistrital.edu.co",
		"PLANTILLA_EVALUACION_DOCENTE",
		data,
		adjuntos,
	)

	fmt.Println("Resultado del envío:", respuesta)

	// sirve la creación del pdf

	/*b64, err := helpers.CrearCertificadoAutoevaluacion(
		"OMER CALDERON",
		"12119277",
		"601 - DOCTORADO INTERINSTITUCIONAL EN EDUCACIÓN",
		"18 del mes 05 de 2025",
	)
	if err != nil {
		log.Fatalf("Error creando PDF: %v", err)
	}
	fmt.Println("Base64 del PDF:", b64)*/

	
	//////////
	evaluadoId := c.GetString("evaluado_id")
	periodoId, err1 := c.GetInt("periodo_id")
	procesoId, err2 := c.GetInt("proceso_id")

	if err1 != nil || err2 != nil || evaluadoId == "" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "Parámetros inválidos: se requieren evaluado_id, periodo_id y proceso_id",
		}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	resultado := services.GetReporteHeteroevaluacionConsejo(evaluadoId, periodoId, procesoId)
	c.Data["json"] = resultado
	c.ServeJSON()
}