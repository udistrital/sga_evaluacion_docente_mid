package helpers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
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
