package services

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"strconv"
)

type RespuestaCoevaluacionIDTO struct {
	EspacioAcademicoId int    `json:"EspacioAcademicoId"`
	NombreEspacio      string `json:"NombreEspacio"`
	IdGrupo            string `json:"IdGrupo"`
	Grupo              string `json:"Grupo"`
	RespuestaPregunta1 string `json:"RespuestaPregunta1"`
	RespuestaPregunta2 string `json:"RespuestaPregunta2"`
	RespuestaPregunta3 string `json:"RespuestaPregunta3"`
	Enlace             string `json:"Enlace"`
}

func GetReporteCoevaluacionIConsejo(evaluadorId string, periodoId int, procesoId int) requestresponse.APIResponse {
	url := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "reporte_coevaluacion_i_consejo?evaluador_id=" + fmt.Sprint(evaluadorId) + "&periodo_id=" + fmt.Sprint(periodoId) + "&proceso_id=" + fmt.Sprint(procesoId)

	var evaluaciones []interface{}
	if err := request.GetJson(url, &evaluaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error en petición GET: %v", err))
	}

	idsEspacios := helpers.ExtractEspacioIDs(evaluaciones)

	mapaEspacios, err := helpers.GetMapaEspacios(idsEspacios)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	var respuestasOrdenadas []RespuestaCoevaluacionIDTO
	for _, item := range evaluaciones {
		if dataMap, ok := item.(map[string]interface{}); ok {
			dto := RespuestaCoevaluacionIDTO{
				EspacioAcademicoId: int(dataMap["EspacioAcademicoId"].(float64)),
				NombreEspacio:      mapaEspacios[strconv.Itoa(int(dataMap["EspacioAcademicoId"].(float64)))],
				IdGrupo:            dataMap["IdGrupo"].(string),
				Grupo:              dataMap["Grupo"].(string),
				RespuestaPregunta1: dataMap["RespuestaPregunta1"].(string),
				RespuestaPregunta2: dataMap["RespuestaPregunta2"].(string),
				RespuestaPregunta3: dataMap["RespuestaPregunta3"].(string),
				Enlace:             dataMap["Enlace"].(string),
			}
			respuestasOrdenadas = append(respuestasOrdenadas, dto)
		}
	}

	result := map[string]interface{}{
		"RespuestasEvaluacion": respuestasOrdenadas,
	}

	return requestresponse.APIResponseDTO(true, 200, result, "Datos obtenidos correctamente.")
}
