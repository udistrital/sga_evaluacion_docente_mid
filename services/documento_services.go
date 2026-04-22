package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_evaluacion_docente_mid/helpers"
	"github.com/udistrital/sga_evaluacion_docente_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func ConsultarDocumentos(periodoId int, evaluadorId int) (APIResponseDTO requestresponse.APIResponse) {
	var res map[string]interface{}
	var documento models.Document
	var documentosBase64 []string

	query := "respuesta/document_uuids/" + fmt.Sprintf("%d", periodoId) + "/" + fmt.Sprintf("%d", evaluadorId)
	err1 := request.GetJson(beego.AppConfig.String("EvaluacionDocenteService")+query, &res)
	if err1 != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 403, nil, fmt.Sprintf("Error al consultar documentos en api_crud: %v", err1))
		return APIResponseDTO
	}

	uuids := res["Data"].([]interface{})
	fmt.Println("UUIDS:", uuids)

	if len(uuids) >= 1 {
		for _, uuid := range uuids {
			//se consultan los documentos por cada uno de los uuids con gestor documental mid
			fmt.Println("UUID:", uuid)
			query = "document/" + fmt.Sprintf("%s", uuid)
			err2 := request.GetJson(beego.AppConfig.String("GestorDocumentalService")+query, &documento)
			if err2 != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 403, nil, fmt.Sprintf("Error al consultar uuids: %v", err1))
				return APIResponseDTO
			} else {
				documentosBase64 = append(documentosBase64, documento.File)
			}
		}
		mergeBase64PDFs, _ := helpers.MergeBase64PDFs(documentosBase64)
		return requestresponse.APIResponseDTO(true, 200, mergeBase64PDFs, "Consulta exitosa")
	} else {
		return requestresponse.APIResponseDTO(true, 204, nil, nil)
	}
}
