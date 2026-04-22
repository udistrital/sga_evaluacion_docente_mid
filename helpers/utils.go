package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/phpdave11/gofpdf"
	"io/ioutil"
	"os"
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

func CrearCertificadoPdfEvaluacion(
	nombre, cedula, espacioAcademico, fecha, nombreEva, periodoAcademico, nombreEvaluado, grupo string,
) (string, error) {
	const (
		margenIzquierdo = 25.0
		margenSuperior  = 20.0
		margenDerecho   = 25.0
		anchoPagina     = 216.0
		altoPagina      = 280.0
		col1Width       = 96.0
		col2Width       = 70.0
	)

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: anchoPagina, Ht: altoPagina},
	})
	pdf.SetMargins(margenIzquierdo, margenSuperior, margenDerecho)
	pdf.AddPage()

	pdf.AddUTF8Font("arial", "", "utils/fonts/arial.ttf")
	pdf.AddUTF8Font("arial", "B", "utils/fonts/arialbd.ttf")
	pdf.SetFont("arial", "", 11)

	pdf.ImageOptions("utils/images/headerPdf.png", 10, 10, 196, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.Ln(50)

	pdf.SetFont("arial", "B", 10)
	pdf.CellFormat(0, 10, "CERTIFICA QUE:", "", 1, "C", false, 0, "")
	pdf.Ln(14)

	agregarContenido := func(textoInicial, evaluado string) {
		pdf.SetFont("arial", "", 10)
		pdf.Write(6, textoInicial)
		pdf.SetFont("arial", "B", 10)
		pdf.Write(6, evaluado)
		pdf.SetFont("arial", "", 10)
	}

	agregarTabla := func(col1, col2 string) {
		pdf.SetFont("arial", "B", 10)
		pdf.CellFormat(col1Width, 6, col1, "1", 0, "C", false, 0, "")
		pdf.CellFormat(col2Width, 6, "TIPO EVALUACIÓN", "1", 1, "C", false, 0, "")
		pdf.SetFont("arial", "", 10)

		x := pdf.GetX()
		y := pdf.GetY()

		pdf.MultiCell(col1Width, 5.5, col2, "1", "C", false)
		h := pdf.GetY() - y
		pdf.SetXY(x+col1Width, y)
		pdf.CellFormat(col2Width, h, nombreEva, "1", 0, "C", false, 0, "")
		pdf.SetY(y + h)
		pdf.Ln(14)
	}

	switch nombreEva {
	case "Autoevaluación I", "Heteroevaluación":
		agregarContenido("El (La) estudiante ", nombre)
		agregarContenido(", identificado (a) con No. ", cedula)
		pdf.Write(6, fmt.Sprintf(", ha realizado la evaluación docente para el periodo académico %s así:", periodoAcademico))
		pdf.Ln(14)
		agregarTabla("ESPACIO ACADÉMICO", espacioAcademico)

	case "Coevaluación I":
		agregarContenido("La persona docente ", nombre)
		agregarContenido(", identificado (a) con cédula de ciudadanía No. ", cedula)
		pdf.Write(6, fmt.Sprintf(", ha realizado la evaluación docente para el grupo %s en el periodo académico %s así:", grupo, periodoAcademico))
		pdf.Ln(14)
		agregarTabla("ESPACIO ACADÉMICO", espacioAcademico)

	case "Autoevaluación II 3", "Autoevaluación II 2", "Autoevaluación II 1":
		agregarContenido("La persona docente ", nombre)
		agregarContenido(", identificado (a) con cédula de ciudadanía No. ", cedula)
		pdf.Write(6, fmt.Sprintf(", ha realizado la evaluación docente para el periodo académico %s así:", periodoAcademico))
		pdf.Ln(14)
		agregarTabla("ESPACIO ACADÉMICO", espacioAcademico)

	case "Coevaluación II":
		agregarContenido("La persona ", nombre)
		agregarContenido(", identificado (a) con No. ", cedula)
		pdf.Write(6, fmt.Sprintf(", ha realizado la evaluación docente para el periodo académico %s así:", periodoAcademico))
		pdf.Ln(14)
		agregarTabla("DOCENTE EVALUADO (A)", nombreEvaluado)

	default:
		return "", fmt.Errorf("tipo de evaluación no reconocido: %s", nombreEva)
	}

	pdf.SetFont("arial", "", 10)
	pdf.MultiCell(0, 6, fmt.Sprintf("La presente certificación se expide automáticamente desde el Sistema de Gestión Académica el día %s.", fecha), "", "L", false)
	pdf.Ln(7)
	pdf.MultiCell(0, 6, "Atento saludo,", "", "L", false)
	pdf.Ln(7)
	pdf.SetFont("arial", "B", 10)
	pdf.MultiCell(0, 6, "COORDINACIÓN OFICINA DE EVALUACIÓN DOCENTE\nUNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS\nSISTEMA SOPORTADO POR LA OATI", "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return "", fmt.Errorf("error generando PDF: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
