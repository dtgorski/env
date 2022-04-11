// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2021

package env

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var errRequirePtrStruct = errors.New("env: unmarshal requires pointer to struct")

// Unmarshal populates the interface with values from process environment.
func Unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr {
		return errRequirePtrStruct
	}
	if rv.Elem().Kind() != reflect.Struct {
		return errRequirePtrStruct
	}

	walk(rv.Elem())
	return nil
}

func walk(v reflect.Value) {
	if v.Type().Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		fld := v.Field(i)
		if !fld.IsValid() || !fld.CanSet() {
			continue
		}
		key, opt := (structTag)(v.Type().Field(i).Tag).Get("env")
		str := loadVal(key, opt.Contains("file"))

		switch fld.Type().Kind() {
		case reflect.Bool:
			fld.SetBool(coerceBool(str))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fld.SetInt(coerceInt64(str))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fld.SetUint(coerceUint64(str))

		case reflect.Float32, reflect.Float64:
			fld.SetFloat(coerceFloat64(str))

		case reflect.Map:
			// FIXME: YAGNI?

		case reflect.Slice:
			switch typ := fld.Type(); typ.Elem().Kind() {
			case reflect.String:
				fld.Set(coerceStringSlice(str))

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fld.Set(coerceIntSlice(typ, str))

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fld.Set(coerceUintSlice(typ, str))

			case reflect.Float32, reflect.Float64:
				fld.Set(coerceFloatSlice(typ, str))
			}

		case reflect.Struct:
			walk(fld)

		case reflect.String:
			fld.SetString(str)
		}
	}
}

func loadVal(key string, fromFile bool) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	if fromFile {
		if val, ok := os.LookupEnv(key + "_FILE"); ok {
			b, _ := ioutil.ReadFile(val)
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}

func coerceStringSlice(s string) reflect.Value {
	sub := split(s, ",")
	val := reflect.MakeSlice(reflect.TypeOf(sub), len(sub), len(sub))
	for i, v := range sub {
		val.Index(i).SetString(v)
	}
	return val
}

func coerceIntSlice(intType reflect.Type, s string) reflect.Value {
	sub := split(s, ",")
	val := reflect.MakeSlice(intType, len(sub), len(sub))
	for i, v := range sub {
		val.Index(i).SetInt(coerceInt64(v))
	}
	return val
}

func coerceUintSlice(uintType reflect.Type, s string) reflect.Value {
	sub := split(s, ",")
	val := reflect.MakeSlice(uintType, len(sub), len(sub))
	for i, v := range sub {
		val.Index(i).SetUint(coerceUint64(v))
	}
	return val
}

func coerceFloatSlice(floatType reflect.Type, s string) reflect.Value {
	sub := split(s, ",")
	val := reflect.MakeSlice(floatType, len(sub), len(sub))
	for i, v := range sub {
		val.Index(i).SetFloat(coerceFloat64(v))
	}
	return val
}

func coerceInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func coerceUint64(s string) uint64 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return u
}

func coerceFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func coerceBool(s string) bool {
	if len(s) > 0 {
		s = strings.ToLower(s)
		for _, v := range []string{"true", "on", "yes", "1"} {
			if v == s {
				return true
			}
		}
	}
	return false
}

func split(s, d string) []string {
	sub := strings.Split(s, d)
	elm := make([]string, 0)

	for _, str := range sub {
		s := strings.TrimSpace(str)
		if s != "" {
			elm = append(elm, s)
		}
	}
	return elm
}
