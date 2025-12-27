package anware

import (
	"fmt"
	"reflect"
	"unicode"
)

func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func extractSubConfig(appConfig any, moduleName string, expectedType any) any {
	return extractSubStruct(appConfig, moduleName, expectedType, "Config")
}

func extractSubStruct(root any, moduleName string, expectedType any, kind string) any {
	if root == nil {
		panic(fmt.Sprintf("[ANWARE] %s root is nil", kind))
	}

	rootVal := reflect.ValueOf(root)
	if rootVal.Kind() == reflect.Ptr {
		rootVal = rootVal.Elem()
	}

	if rootVal.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[ANWARE] %s root must be a struct", kind))
	}

	fieldName := toPascalCase(moduleName)
	field := rootVal.FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf(
			"[ANWARE] %s missing for module '%s' (expected field %s.%s)",
			kind,
			moduleName,
			rootVal.Type().Name(),
			fieldName,
		))
	}

	expected := reflect.TypeOf(expectedType)
	actual := field.Type()

	if expected != actual {
		panic(fmt.Sprintf(
			"[ANWARE] %s type mismatch for module '%s': expected %s, got %s",
			kind,
			moduleName,
			expected,
			actual,
		))
	}

	fieldPtr := field.Addr().Interface()
	fieldVal := reflect.ValueOf(fieldPtr).Elem()

	for i := 0; i < fieldVal.NumField(); i++ {
		subField := fieldVal.Type().Field(i)
		if subField.Name == "" {
			continue
		}

		if !fieldVal.Field(i).IsZero() {
			continue
		}

		if len(subField.Name) > 5 && subField.Name[len(subField.Name)-5:] == "_Glob" {
			globalName := subField.Name[:len(subField.Name)-5]
			globalVal := rootVal.FieldByName(globalName)
			if globalVal.IsValid() && globalVal.Type() == fieldVal.Field(i).Type() {
				fieldVal.Field(i).Set(globalVal)
			}
		}
	}

	return field.Addr().Interface()
}
