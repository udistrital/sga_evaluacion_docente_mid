package helpers

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// Estructuras para mapear el XML
type Asignaturas struct {
	XMLName     xml.Name     `xml:"asignaturas"`
	Asignaturas []Asignatura `xml:"asignatura"`
}

type Asignatura struct {
	Codigo string `xml:"codigo_asignatura"`
	Nombre string `xml:"nombre_asignatura"`
}

// GetMapaEspacios obtiene un mapa de ID a nombre del espacio
func GetMapaEspacios(ids []int) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}

	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, strconv.Itoa(id))
	}
	idsJoin := strings.Join(idStrings, ",")

	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String("UrlWSO2") + beego.AppConfig.String("NsAcademica") + "/espacio_academico/" + idsJoin

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error haciendo petición GET")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error leyendo cuerpo de respuesta")
	}

	var espacios Asignaturas
	if err := xml.Unmarshal(body, &espacios); err != nil {
		return nil, fmt.Errorf("Error parseando XML")
	}

	mapa := make(map[string]string)
	for _, asignatura := range espacios.Asignaturas {
		mapa[asignatura.Codigo] = asignatura.Nombre
	}
	return mapa, nil
}

// ExtractEspacioIDs extrae los IDs de espacio académico de las evaluaciones
func ExtractEspacioIDs(data []interface{}) []int {
	var ids []int
	for _, item := range data {
		if dataMap, ok := item.(map[string]interface{}); ok {
			if idFloat, ok := dataMap["EspacioAcademicoId"].(float64); ok {
				ids = append(ids, int(idFloat))
			}
		}
	}
	return ids
}
