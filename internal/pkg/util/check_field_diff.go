package util

import (
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/oleiade/reflections.v1"
)

// CheckDiff checks what fields have been changed.
// Only checks "String", "Int" and "Float64" types.
func CheckDiff(old interface{}, new interface{}, fieldsToSkip map[string]bool) []string {
	modifiedFields := make([]string, 0)
	structItems, _ := reflections.Items(old)

	for field, oldValue := range structItems {
		if _, ok := fieldsToSkip[field]; ok {
			continue
		}
		fieldKind, _ := reflections.GetFieldKind(old, field)
		if fieldKind != reflect.String && fieldKind != reflect.Int && fieldKind != reflect.Float64 {
			continue
		}
		newValue, _ := reflections.GetField(new, field)
		if newValue != oldValue {
			if fieldKind == reflect.Int {
				modifiedFields = append(modifiedFields, field+": "+strconv.Itoa(oldValue.(int))+" -> "+strconv.Itoa(newValue.(int)))
			} else if fieldKind == reflect.Float64 {
				modifiedFields = append(modifiedFields, field+": "+fmt.Sprintf("%.2f", oldValue.(float64))+" -> "+fmt.Sprintf("%.2f", newValue.(float64)))
			} else {
				modifiedFields = append(modifiedFields, field+": "+oldValue.(string)+" -> "+newValue.(string))
			}
		}
	}

	return modifiedFields
}
