package helpers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"bytes"
	"github.com/phpdave11/gofpdf"
	//"math"

)

func MergeBase64PDFs(base64Docs []string) (string, error) {
	var tempFiles []string

	for i, b64 := range base64Docs {
		data, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return "", fmt.Errorf("decodificando base64 #%d: %w", i, err)
		}
		tmpPath := fmt.Sprintf("tmp_doc_%d.pdf", i)
		if err := ioutil.WriteFile(tmpPath, data, 0o644); err != nil {
			return "", fmt.Errorf("escribiendo %s: %w", tmpPath, err)
		}
		tempFiles = append(tempFiles, tmpPath)
	}

	outPath := "merged.pdf"
	if err := api.MergeCreateFile(tempFiles, outPath, false, nil); err != nil {
		return "", fmt.Errorf("error mergeando PDFs: %w", err)
	}

	mergedData, err := ioutil.ReadFile(outPath)
	if err != nil {
		return "", fmt.Errorf("leyendo %s: %w", outPath, err)
	}
	mergedB64 := base64.StdEncoding.EncodeToString(mergedData)

	for _, f := range append(tempFiles, outPath) {
		_ = os.Remove(f)
	}

	return mergedB64, nil
}

func CrearCertificadoAutoevaluacion(nombre string, cedula string, proyecto string, fecha string) (string, error) {
	//espacio académico = materia
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: 216, Ht: 280}, 
	})
	pdf.SetMargins(25, 20, 25)
	pdf.AddPage()

	pdf.AddUTF8Font("arial", "", "utils/fonts/arial.ttf")
	pdf.AddUTF8Font("arial", "B", "utils/fonts/arialbd.ttf")
	pdf.SetFont("arial", "", 11)

	pdf.ImageOptions("utils/images/headerPdf.png", 10, 10, 196, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.Ln(50)

	pdf.SetFont("arial", "B", 10)
	pdf.CellFormat(0, 10, "CERTIFICA QUE:", "", 1, "C", false, 0, "")
	pdf.Ln(14)

	pdf.SetFont("arial", "", 10)
	pdf.Write(6, "El docente ")
	pdf.SetFont("arial", "B", 10)
	pdf.Write(6, nombre)
	pdf.SetFont("arial", "", 10)
	pdf.Write(6, ", identificado con cédula de ciudadanía No. ")
	pdf.SetFont("arial", "B", 10)
	pdf.Write(6, cedula)
	pdf.SetFont("arial", "", 10)
	pdf.Write(6, ", inscrito en el proyecto curricular listado a continuación, ha realizado la autoevaluación docente para el periodo académico 2025-1 así:")
	pdf.Ln(14)
	
	col1Width := 96.0 
	col2Width := 70.0 
	pdf.SetFont("arial", "B", 10)
	pdf.CellFormat(col1Width, 6, "PROYECTO CURRICULAR", "1", 0, "C", false, 0, "")
	pdf.CellFormat(col2Width, 6, "EVALUACIÓN DOCENTE", "1", 1, "C", false, 0, "")
	pdf.SetFont("arial", "", 10)
	x := pdf.GetX()
	y := pdf.GetY()
	pdf.MultiCell(col1Width, 5.5, proyecto, "1", "C", false)
	h1 := pdf.GetY() - y 
	pdf.SetXY(x+col1Width, y)
	pdf.CellFormat(col2Width, h1, fecha, "1", 0, "C", false, 0, "")
	pdf.SetY(y + h1)
	pdf.Ln(14)

	pdf.SetFont("arial", "", 10)
	textoFinalNormal := fmt.Sprintf("La presente certificación se expide automáticamente desde el Sistema de Gestión Académica el día %s.",
		fecha)
	pdf.MultiCell(0, 6, textoFinalNormal, "", "L", false)
	pdf.Ln(7)

	pdf.SetFont("arial", "", 10)
	pdf.MultiCell(0, 6, "Atento saludo,", "", "L", false)
	pdf.Ln(7)
	
	pdf.SetFont("arial", "B", 10)
	pdf.MultiCell(0, 6, "COORDINACIÓN OFICINA DE EVALUACIÓN DOCENTE\nUNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS\nSISTEMA SOPORTADO POR LA OATI", "", "L", false)

	// Guardar en buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return "", fmt.Errorf("error generando PDF: %w", err)
	}

	// Guardar PDF para revisión visual
	if err := os.WriteFile("Certificado-Autoevaluacion.pdf", buf.Bytes(), 0o644); err != nil {
		return "", fmt.Errorf("error escribiendo archivo: %w", err)
	}

	// Codificar base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil

}