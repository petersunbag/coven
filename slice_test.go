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

	c := newSliceConverter(&convertType{reflect.TypeOf([]bar{}), reflect.TypeOf([]foo{}), nil})

	s := []foo{foo{1}, foo{2}, foo{3}}

	d := make([]bar, 0)
	d = nil

	c.convert(unsafe.Pointer(dereferencedValue(&d).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&s).UnsafeAddr()))

	if expected := `[{"A":1},{"A":2},{"A":3}]`; !reflect.DeepEqual(expected, jsonEncode(d)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(d))
	}

	a := []int{1, 2, 3}
	b := []*byte{}

	c = newSliceConverter(&convertType{reflect.TypeOf([]*byte{}), reflect.TypeOf([]int{}), nil})
	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if expected := `[1,2,3]`; !reflect.DeepEqual(expected, jsonEncode(b)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(b))
	}
}

func TestSameSliceConvert(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4}

	c := newSliceConverter(&convertType{reflect.TypeOf([]int{}), reflect.TypeOf([]int{}), nil})

	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if expected := []int{1, 2, 3}; !reflect.DeepEqual(expected, b) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}
}

func TestOptionSliceStructConvert(t *testing.T) {
	type Foo struct {
		A  string
		B1 string
	}
	type Bar struct {
		A string
		B string
	}

}
