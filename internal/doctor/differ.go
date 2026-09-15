package doctor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/config"
)

// DifferenceType define el tipo de cambio detectado entre la configuración existente y la detectada.
type DifferenceType string

const (
	DiffAdded    DifferenceType = "ADDED"    // Propiedad agregada
	DiffModified DifferenceType = "MODIFIED" // Propiedad modificada
	DiffRemoved  DifferenceType = "REMOVED"  // Propiedad eliminada
)

// Difference representa una diferencia individual en una propiedad.
type Difference struct {
	Type     DifferenceType `json:"type"`
	Section  string         `json:"section"`
	Field    string         `json:"field"`
	OldValue string         `json:"oldValue"`
	NewValue string         `json:"newValue"`
}

// ComputeDifferences compara la configuración existente con la detectada y devuelve una lista de diferencias.
func ComputeDifferences(existing *config.KokoConfig, detected config.KokoConfig) []Difference {
	var diffs []Difference

	if existing == nil {
		diffs = append(diffs, Difference{
			Type:     DiffAdded,
			Section:  "project",
			Field:    "koko.config.json",
			OldValue: "<missing>",
			NewValue: "new manifest will be created",
		})
		return diffs
	}

	oldMap := toMap(existing)
	newMap := toMap(detected)

	// Omitir campos que no forman parte de la arquitectura técnica dinámica
	delete(oldMap, "project")
	delete(newMap, "project")
	delete(oldMap, "$schema")
	delete(newMap, "$schema")

	diffMaps(&diffs, "", oldMap, newMap)

	return diffs
}

func toMap(v interface{}) map[string]interface{} {
	data, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(data, &m)
	if m == nil {
		m = make(map[string]interface{})
	}
	return m
}

func diffMaps(diffs *[]Difference, prefix string, oldMap, newMap map[string]interface{}) {
	// 1. Revisar campos en newMap (ADDED o MODIFIED)
	for key, newVal := range newMap {
		currentPath := key
		if prefix != "" {
			currentPath = prefix + "." + key
		}

		oldVal, exists := oldMap[key]

		if !exists || oldVal == nil {
			if newVal != nil && !isEmpty(newVal) {
				// Si newVal es un sub-mapa, descender recursivamente para mostrar cada propiedad limpia
				if newSub, isNewSub := newVal.(map[string]interface{}); isNewSub {
					diffMaps(diffs, currentPath, make(map[string]interface{}), newSub)
					continue
				}
				addDiff(diffs, DiffAdded, currentPath, "<none>", formatValue(newVal))
			}
			continue
		}

		// Si ambos son sub-mapas, descender recursivamente
		oldSub, isOldSub := oldVal.(map[string]interface{})
		newSub, isNewSub := newVal.(map[string]interface{})
		if isOldSub && isNewSub {
			diffMaps(diffs, currentPath, oldSub, newSub)
			continue
		}

		// Comparar valores primitivos o slices
		if !reflect.DeepEqual(oldVal, newVal) {
			if !isEmpty(oldVal) || !isEmpty(newVal) {
				addDiff(diffs, DiffModified, currentPath, formatValue(oldVal), formatValue(newVal))
			}
		}
	}

	// 2. Revisar campos en oldMap que ya no existen en newMap (REMOVED)
	for key, oldVal := range oldMap {
		currentPath := key
		if prefix != "" {
			currentPath = prefix + "." + key
		}

		if _, exists := newMap[key]; !exists {
			if oldVal != nil && !isEmpty(oldVal) {
				if oldSub, isOldSub := oldVal.(map[string]interface{}); isOldSub {
					diffMaps(diffs, currentPath, oldSub, make(map[string]interface{}))
					continue
				}
				addDiff(diffs, DiffRemoved, currentPath, formatValue(oldVal), "<none>")
			}
		}
	}
}

func formatValue(val interface{}) string {
	if val == nil {
		return "<none>"
	}
	if sl, ok := val.([]interface{}); ok {
		var items []string
		for _, item := range sl {
			items = append(items, fmt.Sprintf("%v", item))
		}
		return strings.Join(items, ", ")
	}
	return fmt.Sprintf("%v", val)
}

func addDiff(diffs *[]Difference, diffType DifferenceType, fullPath, oldVal, newVal string) {
	parts := strings.Split(fullPath, ".")
	section := "root"
	field := fullPath
	if len(parts) > 1 {
		section = strings.Join(parts[:len(parts)-1], " > ")
		field = parts[len(parts)-1]
	}

	*diffs = append(*diffs, Difference{
		Type:     diffType,
		Section:  section,
		Field:    field,
		OldValue: oldVal,
		NewValue: newVal,
	})
}

func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}
	if s, ok := val.(string); ok && s == "" {
		return true
	}
	if m, ok := val.(map[string]interface{}); ok && len(m) == 0 {
		return true
	}
	if sl, ok := val.([]interface{}); ok && len(sl) == 0 {
		return true
	}
	return false
}
