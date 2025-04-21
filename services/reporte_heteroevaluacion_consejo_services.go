package services

import (
	"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
)

type RespuestaHeteroevaluacionDTO struct {
	EspacioAcademicoId int         `json:"EspacioAcademicoId"`
	NombreEspacio      string      `json:"NombreEspacio"`
	Ambito1            interface{} `json:"Ambito1"`
	Ambito2            interface{} `json:"Ambito2"`
	Ambito3            interface{} `json:"Ambito3"`
	PromedioFinal      interface{} `json:"PromedioFinal"`
}

func GetReporteHeteroevaluacionConsejo(evaluadoId string, periodoId int, procesoId int) requestresponse.APIResponse {
	url := "http://" + beego.AppConfig.String("EvaluacionDocenteService") + "reporte_heteroevaluacion_consejo?evaluado_id=" + fmt.Sprint(evaluadoId) + "&periodo_id=" + fmt.Sprint(periodoId) + "&proceso_id=" + fmt.Sprint(procesoId)

	var evaluaciones []interface{}
	if err := request.GetJson(url, &evaluaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("Error en petición GET: %v", err))
	}

	idsEspacios := helpers.ExtractEspacioIDs(evaluaciones)

	mapaEspacios, err := helpers.GetMapaEspacios(idsEspacios)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	var respuestasOrdenadas []RespuestaHeteroevaluacionDTO
	for _, item := range evaluaciones {
		if dataMap, ok := item.(map[string]interface{}); ok {
			dto := RespuestaHeteroevaluacionDTO{
				EspacioAcademicoId: int(dataMap["EspacioAcademicoId"].(float64)),
				NombreEspacio:      mapaEspacios[strconv.Itoa(int(dataMap["EspacioAcademicoId"].(float64)))],
				Ambito1:            dataMap["Ambito1"],
				Ambito2:            dataMap["Ambito2"],
				Ambito3:            dataMap["Ambito3"],
				PromedioFinal:      dataMap["PromedioFinal"],
			}
			respuestasOrdenadas = append(respuestasOrdenadas, dto)
		}
	}

	result := map[string]interface{}{
		"RespuestasEvaluacion": respuestasOrdenadas,
	}

	return requestresponse.APIResponseDTO(true, 200, result, "Datos obtenidos correctamente.")
}

