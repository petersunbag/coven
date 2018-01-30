package coven

import (
	"reflect"
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
