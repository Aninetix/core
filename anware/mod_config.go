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

func extractSubConfig(
	appConfig any,
	moduleName string,
	expectedType any,
) (any, error) {
	return extractSubStruct(appConfig, moduleName, expectedType)
}

func extractSubStruct(
	root any,
	moduleName string,
	expectedType any,
) (any, error) {

	if root == nil {
		panic("[ANWARE] FATAL: config root is nil")
	}

	rootVal := reflect.ValueOf(root)
	if rootVal.Kind() == reflect.Ptr {
		rootVal = rootVal.Elem()
	}

	if rootVal.Kind() != reflect.Struct {
		panic(fmt.Sprintf(
			"[ANWARE] FATAL: config root must be a struct, got %s",
			rootVal.Kind(),
		))
	}

	fieldName := toPascalCase(moduleName)
	field := rootVal.FieldByName(fieldName)

	if !field.IsValid() {
		return nil, fmt.Errorf(
			"[ANWARE] module '%s' disabled: missing config field %s.%s",
			moduleName,
			rootVal.Type().Name(),
			fieldName,
		)
	}

	expected := reflect.TypeOf(expectedType)
	actual := field.Type()

	if expected != actual {
		return nil, fmt.Errorf(
			"[ANWARE] module '%s' disabled: config type mismatch (expected %s, got %s)",
			moduleName,
			expected,
			actual,
		)
	}

	fieldPtr := field.Addr().Interface()
	fieldVal := reflect.ValueOf(fieldPtr).Elem()

	for i := 0; i < fieldVal.NumField(); i++ {
		subField := fieldVal.Type().Field(i)

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

	return fieldPtr, nil
}
