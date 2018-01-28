package coven

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestGeneralConverter_Convert(t *testing.T) {
	type foo struct {
		A int
		B byte
	}

	dstTyp := dereferencedType(reflect.TypeOf(new(foo)))
	srcTyp := dereferencedType(reflect.TypeOf(new(foo)))

	c := newGeneralConverter(&convertType{dstTyp, srcTyp})
	foo1 := &foo{1, 2}
	foo2 := foo{}
	c.convert(unsafe.Pointer(dereferencedValue(&foo2).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&foo1).UnsafeAddr()))
	if expected := `{"A":1,"B":2}`; !reflect.DeepEqual(expected, jsonEncode(foo2)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(foo2))
	}

	dstTyp = dereferencedType(reflect.TypeOf(new(***int)))
	srcTyp = dereferencedType(reflect.TypeOf(new(***int)))
	c = newGeneralConverter(&convertType{dstTyp, srcTyp})
	x := 1
	y := &x
	z := &y
	X := &z

	o := 0
	p := &o
	q := &p
	//Y := &q
	c.convert(unsafe.Pointer(dereferencedValue(&q).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&X).UnsafeAddr()))
	if expected := `1`; !reflect.DeepEqual(expected, jsonEncode(q)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(q))
	}
}
