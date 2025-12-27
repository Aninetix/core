package anware

import (
	"fmt"
	"reflect"
	"unicode"
)

// --- utils ---

func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// anHttp -> AnHttp
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// --- CONFIG ---

func extractSubConfig(appConfig any, moduleName string, expectedType any) any {
	return extractSubStruct(appConfig, moduleName, expectedType, "Config")
}

// --- CORE LOGIC ---

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

	// return pointer to field
	return field.Addr().Interface()
}
