package helpers

import "reflect"

func GetFieldString(s any, name string) string {
	rv := reflect.ValueOf(s).Elem()
	fv := rv.FieldByName(name)
	if fv.IsValid() && fv.Kind() == reflect.String {
		return fv.String()
	}
	return ""
}

func GetFieldBool(s any, name string) bool {
	rv := reflect.ValueOf(s).Elem()
	fv := rv.FieldByName(name)
	if fv.IsValid() && fv.Kind() == reflect.Bool {
		return fv.Bool()
	}
	return false
}

func GetFieldInt(s any, name string) int {
	rv := reflect.ValueOf(s).Elem()
	fv := rv.FieldByName(name)

	if !fv.IsValid() {
		return 0
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(fv.Int())
	}

	return 0
}

func GetFieldIntUniversal(s any, name string) int64 {
	rv := reflect.ValueOf(s).Elem()
	fv := rv.FieldByName(name)

	if !fv.IsValid() {
		return 0
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fv.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(fv.Uint())
	}

	return 0
}
