package coven

import (
	"reflect"
	"unsafe"
)

func dereferencedType(t reflect.Type) reflect.Type {
	for k := t.Kind(); k == reflect.Ptr || k == reflect.Interface; k = t.Kind() {
		t = t.Elem()
	}

	return t
}

func dereferencedValue(value interface{}) reflect.Value {
	v := reflect.ValueOf(value)

	for k := v.Kind(); k == reflect.Ptr || k == reflect.Interface; k = v.Kind() {
		v = v.Elem()
	}

	return v
}

func referDeep(t reflect.Type) (reflect.Type, int) {
	d := 0
	for k := t.Kind(); k == reflect.Ptr; k = t.Kind() {
		t = t.Elem()
		d += 1
	}
	return t, d
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
