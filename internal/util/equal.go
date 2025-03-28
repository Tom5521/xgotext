package util

import (
	"math"
	"path/filepath"
	"reflect"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

func FuzzyEqual(x, y string) bool {
	return fuzzy.Ratio(x, y) >= 80
}

func IsSimilarButNotIdentical(x, y string) bool {
	ratio := fuzzy.Ratio(x, y)
	return ratio >= 80 && ratio < 100
}

func EqualPaths(x, y string) bool {
	abs1, err := filepath.Abs(x)
	abs2, err1 := filepath.Abs(y)
	if err != nil || err1 != nil {
		return x == y
	}

	return abs1 == abs2
}

func Equal[X, Y any](x X, y Y) bool {
	typeX, typeY := reflect.TypeOf(x), reflect.TypeOf(y)
	valueX, valueY := reflect.ValueOf(x), reflect.ValueOf(y)

	if typeX.Kind() != typeY.Kind() {
		return false
	}

	if typeX.Kind() == reflect.Pointer {
		if valueX.IsNil() || valueY.IsNil() {
			return valueX.IsNil() == valueY.IsNil()
		}
		typeX = typeX.Elem()
		valueX = valueX.Elem()
		typeY = typeY.Elem()
		valueY = valueY.Elem()
	}
	if typeX != typeY {
		return false
	}

	switch typeX.Kind() {
	case reflect.Bool:
		return valueX.Bool() == valueY.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueX.Int() == valueY.Int()
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return valueX.Uint() == valueY.Uint()
	case reflect.Float32, reflect.Float64:
		return floatEqual(valueX.Float(), valueY.Float())
	case reflect.Complex64, reflect.Complex128:
		return complexEqual(valueX.Complex(), valueY.Complex())
	case reflect.String:
		return valueX.String() == valueY.String()
	case reflect.Array, reflect.Slice:
		return sliceEqual(valueX, valueY)
	case reflect.Map:
		return mapEqual(valueX, valueY)
	case reflect.Struct:
		return structEqual(valueX, valueY)
	case reflect.Interface:
		return interfaceEqual(valueX, valueY)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return false
	default:
		return false
	}
}

func floatEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

func complexEqual(a, b complex128) bool {
	return floatEqual(real(a), real(b)) && floatEqual(imag(a), imag(b))
}

func sliceEqual(v1, v2 reflect.Value) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	for i := 0; i < v1.Len(); i++ {
		if !Equal(v1.Index(i).Interface(), v2.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func mapEqual(v1, v2 reflect.Value) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	for _, key := range v1.MapKeys() {
		val1 := v1.MapIndex(key)
		val2 := v2.MapIndex(key)
		if !val2.IsValid() || !Equal(val1.Interface(), val2.Interface()) {
			return false
		}
	}
	return true
}

func structEqual(v1, v2 reflect.Value) bool {
	for i := 0; i < v1.NumField(); i++ {
		field := v1.Type().Field(i)
		if field.IsExported() {
			if !Equal(v1.Field(i).Interface(), v2.Field(i).Interface()) {
				return false
			}
		}
	}
	return true
}

func interfaceEqual(v1, v2 reflect.Value) bool {
	if v1.IsNil() || v2.IsNil() {
		return v1.IsNil() == v2.IsNil()
	}
	return Equal(v1.Elem().Interface(), v2.Elem().Interface())
}
