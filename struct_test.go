package coven

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestSimpleConvert(t *testing.T) {
	type Foo struct {
		A int64
		B string
		C *string
		D *int
		E []int
		f int
	}
	type Bar struct {
		A int
		B *string
		C string
		D **int
		E []*int
		f int64
	}

	c := newStructConverter(&convertType{reflect.TypeOf(Bar{}), reflect.TypeOf(Foo{})})

	s := "b"
	i := 2

	bar := Bar{}
	bb := &bar

	foo := Foo{1, "a", &s, &i, []int{1, 2, 3}, 4}
	c.convert(unsafe.Pointer(dereferencedValue(&bb).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&foo).UnsafeAddr()))
	if expected := `{"A":1,"B":"a","C":"b","D":2,"E":[1,2,3]}`; expected != jsonEncode(bb) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(bb))
	}

	foo2 := Foo{1, "a", nil, nil, nil, 5}
	c.convert(unsafe.Pointer(dereferencedValue(&bar).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&foo2).UnsafeAddr()))
	if expected := `{"A":1,"B":"a","C":"","D":0,"E":[]}`; expected != jsonEncode(bar) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(bar))
	}
}

func TestNestedConvert(t *testing.T) {
	type Baz struct {
		A int
		B string
	}
	type Foo struct {
		Baz
		B string
		C *string
	}
	type Bar struct {
		Baz
		C string
	}

	type FooBar struct {
		A int64
		*Foo
	}

	type BarFoo struct {
		Foo Bar
	}

	c1 := newStructConverter(&convertType{reflect.TypeOf(BarFoo{}), reflect.TypeOf(FooBar{})})
	c2 := newStructConverter(&convertType{reflect.TypeOf(FooBar{}), reflect.TypeOf(BarFoo{})})

	barFoo := BarFoo{}

	foobar := FooBar{10, &Foo{Baz{1, "b"}, "B", stringPtr("c")}}
	c1.convert(unsafe.Pointer(dereferencedValue(&barFoo).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&foobar).UnsafeAddr()))
	if expected := `{"Foo":{"A":1,"B":"B","C":"c"}}`; expected != jsonEncode(barFoo) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(barFoo))
	}

	foobar = FooBar{}
	c2.convert(unsafe.Pointer(dereferencedValue(&foobar).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&barFoo).UnsafeAddr()))
	if expected := `{"A":0,"B":"B","C":"c"}`; expected != jsonEncode(foobar) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(foobar))
	}

	foobar = FooBar{10, nil}
	c1.convert(unsafe.Pointer(dereferencedValue(&barFoo).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&foobar).UnsafeAddr()))
	if expected := `{"Foo":{"A":0,"B":"","C":""}}`; expected != jsonEncode(barFoo) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(barFoo))
	}
	c2.convert(unsafe.Pointer(dereferencedValue(&foobar).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&barFoo).UnsafeAddr()))
	if expected := `{"A":10,"B":"","C":""}`; expected != jsonEncode(foobar) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(foobar))
	}
}

func TestExtractFields(t *testing.T) {
	type foo struct {
		A int
		B int
		D int
	}
	type bar struct {
		foo
		B int
		C int
	}
	type foobar struct {
		bar
		C   int
		D   int
		foo foo
	}

	align := unsafe.Alignof(1)

	fs, fm := extractFields(reflect.TypeOf(foobar{}), 0)
	if len(fs) != 4 {
		t.Fatalf("[expected:%v] [actual:%v]", 4, jsonEncode(fs))
	}
	if fm["A"].Offset != 0 {
		t.Fatalf("[expected:%v] [actual:%v]", 0, fm["A"].Offset)
	}
	if fm["B"].Offset != 3*align {
		t.Fatalf("[expected:%v] [actual:%v]", 0, fm["B"].Offset)
	}
	if fm["C"].Offset != 5*align {
		t.Fatalf("[expected:%v] [actual:%v]", 0, fm["C"].Offset)
	}
	if fm["D"].Offset != 6*align {
		t.Fatalf("[expected:%v] [actual:%v]", 0, fm["D"].Offset)
	}
}

func stringPtr(s string) *string {
	return &s
}
