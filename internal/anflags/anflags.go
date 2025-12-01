package anflags

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ParseFlags(s any) error {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("ParseFlags: argument must be pointer to struct")
	}
	st := rv.Elem()
	typ := st.Type()

	// map index -> *string (valeur temporaire de flag)
	tmp := make(map[int]*string)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// n'enregistrer que les champs exportés
		if !field.IsExported() {
			continue
		}

		name := field.Tag.Get("flag")
		if name == "" {
			// fallback : nom du champ en minuscules
			name = lowerFirst(field.Name)
		}
		def := field.Tag.Get("default")
		usage := field.Tag.Get("usage")

		// register as string flag, we'll convert after Parse
		ptr := flag.String(name, def, usage)
		tmp[i] = ptr
	}

	flag.Parse()

	// après Parse(), convertir et assigner
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		ptrStr, ok := tmp[i]
		if !ok {
			continue
		}
		valStr := *ptrStr
		fv := st.Field(i)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(valStr)

		case reflect.Bool:
			b, err := strconv.ParseBool(valStr)
			if err != nil {
				return fmt.Errorf("flag %s: %w", field.Name, err)
			}
			fv.SetBool(b)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// handle time.Duration specially if tagged type is time.Duration
			if field.Type == reflect.TypeOf(time.Duration(0)) {
				d, err := time.ParseDuration(valStr)
				if err != nil {
					return fmt.Errorf("flag %s (duration): %w", field.Name, err)
				}
				fv.SetInt(int64(d))
				continue
			}
			bits := fv.Type().Bits()
			n, err := strconv.ParseInt(valStr, 10, bits)
			if err != nil {
				return fmt.Errorf("flag %s: %w", field.Name, err)
			}
			fv.SetInt(n)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bits := fv.Type().Bits()
			n, err := strconv.ParseUint(valStr, 10, bits)
			if err != nil {
				return fmt.Errorf("flag %s: %w", field.Name, err)
			}
			fv.SetUint(n)

		case reflect.Float32, reflect.Float64:
			bits := fv.Type().Bits()
			f, err := strconv.ParseFloat(valStr, bits)
			if err != nil {
				return fmt.Errorf("flag %s: %w", field.Name, err)
			}
			fv.SetFloat(f)

		default:
			return fmt.Errorf("flag %s: type %s non supportée", field.Name, fv.Kind())
		}
	}

	return nil
}

// Helper: lowercase first letter (simple fallback)
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if 'A' <= r[0] && r[0] <= 'Z' {
		r[0] = r[0] - 'A' + 'a'
	}
	return string(r)
}
