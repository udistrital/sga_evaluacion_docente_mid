package services

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"strconv"
)

type RespuestaAutoevaluacionIItresDTO struct {
	Documento          string      `json:"Documento"`
	EspacioAcademicoId int         `json:"EspacioAcademicoId"`
	NombreEspacio      string      `json:"NombreEspacio"`
	Promedio           float64     `json:"Promedio"`
	RespuestaPregunta1 interface{} `json:"RespuestaPregunta1"`
	RespuestaPregunta2 interface{} `json:"RespuestaPregunta2"`
	RespuestaPregunta3 interface{} `json:"RespuestaPregunta3"`
	Enlace             interface{} `json:"Enlace"`
}

func GetReporteAutoevaluacionIITresConsejo(evaluadorId string, periodoId int, procesoId int, nombreEvaluador string) requestresponse.APIResponse {
	url := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "reporte_autoevaluacion_ii_tres_consejo?evaluador_id=" + fmt.Sprint(evaluadorId) + "&periodo_id=" + fmt.Sprint(periodoId) + "&proceso_id=" + fmt.Sprint(procesoId)

	var evaluaciones []interface{}
	if err := request.GetJson(url, &evaluaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error en petición GET: %v", err))
	}

	idsEspacios := helpers.ExtractEspacioIDs(evaluaciones)

	mapaEspacios, err := helpers.GetMapaEspacios(idsEspacios)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	var respuestasOrdenadas []RespuestaAutoevaluacionIItresDTO
	for _, item := range evaluaciones {
		if dataMap, ok := item.(map[string]interface{}); ok {
			dto := RespuestaAutoevaluacionIItresDTO{
				Documento:          fmt.Sprint(dataMap["Documento"]),
				EspacioAcademicoId: int(dataMap["EspacioAcademicoId"].(float64)),
				NombreEspacio:      mapaEspacios[strconv.Itoa(int(dataMap["EspacioAcademicoId"].(float64)))],
				Promedio:           dataMap["Promedio"].(float64),
				RespuestaPregunta1: dataMap["RespuestaPregunta1"],
				RespuestaPregunta2: dataMap["RespuestaPregunta2"],
				RespuestaPregunta3: dataMap["RespuestaPregunta3"],
				Enlace:             dataMap["Enlace"],
			}
			respuestasOrdenadas = append(respuestasOrdenadas, dto)
		}
	}

	result := map[string]interface{}{
		"NombreEvaluador":      nombreEvaluador,
		"RespuestasEvaluacion": respuestasOrdenadas,
	}

	return requestresponse.APIResponseDTO(true, 200, result, "Datos obtenidos correctamente.")
}
