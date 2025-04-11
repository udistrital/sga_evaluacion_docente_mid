package services

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	//"github.com/udistrital/sga_evaluacion_docente_mid/models"
)

func GetReporteAutoevaluacionIITresConsejo(evaluadorId string, periodoId int, procesoId int) requestresponse.APIResponse {
	url := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "reporte_autoevaluacion_ii_tres_consejo?evaluador_id=" + fmt.Sprint(evaluadorId) + "&periodo_id=" + fmt.Sprint(periodoId) + "&proceso_id=" + fmt.Sprint(procesoId)
	// reporte_autoevaluacion_ii_tres_consejo?evaluador_id=80033827&periodo_id=61&proceso_id=6995

	fmt.Println("🌐 URL: ", url)

	var response []interface{}
	err := request.GetJson(url, &response)
	if err != nil {
		fmt.Println("❌ ERROR al hacer request.GetJson:", err)
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error en petición GET: %v", err))
	}

	fmt.Println("✅ Data recibida:", response)

	return requestresponse.APIResponseDTO(true, 200, response, "Datos obtenidos correctamente.")
}