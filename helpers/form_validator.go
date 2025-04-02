package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Normalize(val interface{}) string {
	if val == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", val))
}

func ToGrupoSlice(val interface{}) ([]map[string]string, error) {
	switch v := val.(type) {
	case []interface{}:
		var result []map[string]string
		for _, item := range v {
			obj, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("item no es un objeto")
			}
			grupo := map[string]string{
				"id_grupo": Normalize(obj["id_grupo"]),
				"grupo":    Normalize(obj["grupo"]),
			}
			result = append(result, grupo)
		}
		return result, nil

	case map[string]interface{}:
		return []map[string]string{
			{
				"id_grupo": Normalize(v["id_grupo"]),
				"grupo":    Normalize(v["grupo"]),
			},
		}, nil

	case string:
		var parsed interface{}
		if err := json.Unmarshal([]byte(v), &parsed); err != nil {
			return nil, fmt.Errorf("error al deserializar string JSON de grupos: %w", err)
		}
		return ToGrupoSlice(parsed)

	default:
		return nil, fmt.Errorf("tipo de valor no soportado: %T", val)
	}
}

func GruposIguales(a, b []map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	setA := make(map[string]struct{})
	setB := make(map[string]struct{})

	for _, g := range a {
		key := fmt.Sprintf("%s:%s", g["id_grupo"], g["grupo"])
		setA[key] = struct{}{}
	}
	for _, g := range b {
		key := fmt.Sprintf("%s:%s", g["id_grupo"], g["grupo"])
		setB[key] = struct{}{}
	}

	if len(setA) != len(setB) {
		return false
	}

	for k := range setA {
		if _, ok := setB[k]; !ok {
			return false
		}
	}

	return true
}
