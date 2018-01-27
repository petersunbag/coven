package coven

import (
	"reflect"
	"testing"
)

func TestGeneralConverter_Convert(t *testing.T) {
	type foo struct {
		A int
	}

	dstTyp := dereferencedType(reflect.TypeOf(new(foo)))
	srcTyp := dereferencedType(reflect.TypeOf(new(foo)))

	c := newGeneralConverter(&convertType{dstTyp, srcTyp})
	foo1 := &foo{1}
	foo2 := foo{}
	c.Convert(&foo2, &foo1)
	if expected := `{"A":1}`; !reflect.DeepEqual(expected, jsonEncode(foo2)) {
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
	c.Convert(&q, &X)
	if expected := `1`; !reflect.DeepEqual(expected, jsonEncode(q)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(q))
	}
}
