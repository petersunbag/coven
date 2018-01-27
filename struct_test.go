package coven

import (
	"encoding/json"
	"reflect"
	"testing"
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
	c.Convert(&bb, &foo)
	if expected := `{"A":1,"B":"a","C":"b","D":2,"E":[1,2,3]}`; expected != jsonEncode(bb) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(bb))
	}

	foo2 := Foo{1, "a", nil, nil, nil, 5}
	c.Convert(&bar, &foo2)
	if expected := `{"A":1,"B":"a","C":"","D":0,"E":null}`; expected != jsonEncode(bar) {
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
	c1.Convert(&barFoo, &foobar)
	if expected := `{"Foo":{"A":1,"B":"B","C":"c"}}`; expected != jsonEncode(barFoo) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(barFoo))
	}

	foobar = FooBar{}
	c2.Convert(&foobar, &barFoo)
	if expected := `{"A":0,"B":"B","C":"c"}`; expected != jsonEncode(foobar) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(foobar))
	}

	foobar = FooBar{10, nil}
	c1.Convert(&barFoo, &foobar)
	if expected := `{"Foo":{"A":0,"B":"","C":""}}`; expected != jsonEncode(barFoo) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(barFoo))
	}
	c2.Convert(&foobar, &barFoo)
	if expected := `{"A":10,"B":"","C":""}`; expected != jsonEncode(foobar) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(foobar))
	}
}

func TestFieldIndex(t *testing.T) {
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

	index := fieldIndex(reflect.TypeOf(foobar{}), []int{})
	if expected := [][]int{{1}, {2}, {0, 1}, {0, 0, 0}}; !reflect.DeepEqual(expected, index) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, index)
	}
}

func jsonEncode(s interface{}) string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

func stringPtr(s string) *string {
	return &s
}
