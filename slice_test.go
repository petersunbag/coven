package coven

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestSliceConverter_Convert(t *testing.T) {
	type foo struct {
		A int
	}

	type bar struct {
		A *int
	}

	c := newSliceConverter(&convertType{reflect.TypeOf([]bar{}), reflect.TypeOf([]foo{})})

	s := []foo{foo{1}, foo{2}, foo{3}}
	d := make([]bar, 0)

	c.convert(unsafe.Pointer(dereferencedValue(&d).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&s).UnsafeAddr()))

	if expected := `[{"A":1},{"A":2},{"A":3}]`; !reflect.DeepEqual(expected, jsonEncode(d)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(d))
	}
}
