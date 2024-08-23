package services

import (
	"encoding/json"
	"fmt"

	"github.com/udistrital/utils_oas/requestresponse"
)

func GuardarRespuestas(data []byte) (APIResponseDTO requestresponse.APIResponse) {

	var dataSource map[string]interface{}
	var respuestas []map[string]interface{}
	var nuevaRespuesta map[string]interface{}

	if err := json.Unmarshal(data, &dataSource); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Error al parsear el JSON")
		return APIResponseDTO
	}

	if r, ok := dataSource["respuestas"].([]interface{}); ok {
		for _, resp := range r {
			if item, ok := resp.(map[string]interface{}); ok {
				metadata := make(map[string]interface{})

				if itemID, ok := item["item_id"]; ok {
					metadata["item_id"] = itemID
				}
				if valor, ok := item["valor"]; ok {
					metadata["valor"] = valor
				}
				if archivos, ok := item["archivos"]; ok {
					metadata["archivos"] = archivos
				}

				nuevaRespuesta = map[string]interface{}{
					"Activo":            true,
					"FechaCreacion":     "2024-08-23T00:00:00Z",
					"FechaModificacion": "2024-08-23T00:00:00Z",
					"Metadata":          metadata,
				}

				//err := guardar

				fmt.Println(nuevaRespuesta)
				respuestas = append(respuestas, nuevaRespuesta)
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "No se encontraron respuestas válidas")
		return APIResponseDTO
	}

	APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuestas, nil)
	return APIResponseDTO
}
