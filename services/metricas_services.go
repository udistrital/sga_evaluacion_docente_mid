package services

import (
	"encoding/json"
	"fmt"

	"github.com/udistrital/utils_oas/requestresponse"
)

func MetricasHeteroevaluacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	tipoReporte, ok := dataSource["tipo_reporte"].(string)
	if !ok || (tipoReporte != "global" && tipoReporte != "facultad") {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Tipo de reporte no válido o no especificado")
		return APIResponseDTO
	}

	if tipoReporte == "global" {
		return ReporteGlobal(data)
	} else if tipoReporte == "facultad" {
		return ReporteFacultad(data)
	}

	return requestresponse.APIResponseDTO(false, 500, nil, "Error inesperado al procesar el reporte")
}

func ReporteGlobal(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte global procesado exitosamente")
	return APIResponseDTO
}

func ReporteFacultad(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte de facultad procesado exitosamente")
	return APIResponseDTO
}

func MetricasAutoevaluacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	tipoReporte, ok := dataSource["tipo_reporte"].(string)
	if !ok || (tipoReporte != "global" && tipoReporte != "facultad") {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Tipo de reporte no válido o no especificado")
		return APIResponseDTO
	}

	if tipoReporte == "estudiante" {
		return ReporteEstudiante(data)
	} else if tipoReporte == "docente" {
		return ReporteDocente(data)
	}

	return requestresponse.APIResponseDTO(false, 500, nil, "Error inesperado al procesar el reporte")
}

func ReporteEstudiante(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte global procesado exitosamente")
	return APIResponseDTO
}

func ReporteDocente(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte de facultad procesado exitosamente")
	return APIResponseDTO
}
func MetricasCoevaluacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var dataSource map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Sprintf("Error al parsear el JSON: %v", err))
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, dataSource, "Reporte de facultad procesado exitosamente")
	return APIResponseDTO
}
