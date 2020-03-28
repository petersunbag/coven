package coven

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestGeneralConverter_Convert(t *testing.T) {
	dstTyp := dereferencedType(reflect.TypeOf(new(***int)))
	srcTyp := dereferencedType(reflect.TypeOf(new(***int)))
	c := newBasicConverter(&convertType{dstTyp, srcTyp, nil})
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

func TestStringCvt(t *testing.T) {
	a := []byte{'a', 'b', 'c'}
	b := ""
	c := newBasicConverter(&convertType{reflect.TypeOf(b), reflect.TypeOf(a), nil})
	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))

	if expected := "abc"; expected != b {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}

	d := []rune{'e', 'f', 'g'}
	c = newBasicConverter(&convertType{reflect.TypeOf(b), reflect.TypeOf(d), nil})
	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&d))

	if expected := "efg"; expected != b {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}

	c = newBasicConverter(&convertType{reflect.TypeOf(a), reflect.TypeOf(b), nil})
	c.convert(unsafe.Pointer(&a), unsafe.Pointer(&b))

	if expected := []byte{'e', 'f', 'g'}; !reflect.DeepEqual(expected, a) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, a)
	}

	b = "xyz"
	c = newBasicConverter(&convertType{reflect.TypeOf(d), reflect.TypeOf(b), nil})
	c.convert(unsafe.Pointer(&d), unsafe.Pointer(&b))

	if expected := []rune{'x', 'y', 'z'}; !reflect.DeepEqual(expected, d) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, d)
	}
}
