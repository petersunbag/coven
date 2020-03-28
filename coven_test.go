package coven

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	type foo struct {
		A int
	}

	type bar struct {
		A int
	}

	_, err := NewConverter(new(foo), new(bar))
	if err != nil {
		panic(err)
	}
	_, err = NewConverter(new(foo), new(bar))
	if err != nil {
		panic(err)
	}

	if len(createdConverters) != 2 {
		t.Fatalf("cache fail")
	}
}

func TestDelegateConverter_Convert(t *testing.T) {
	type foobar struct {
		D int
	}
	type Foo struct {
		A []int
		B map[int64][]byte
		C byte
		foobar
		E []int
	}

	type Bar struct {
		A []*int
		B map[string]*string
		C *byte
		D int64
		E []int
	}

	c, err := NewConverter(Bar{}, Foo{})
	if err != nil {
		panic(err)
	}

	foo := Foo{[]int{1, 2, 3}, map[int64][]byte{1: []byte{'a', 'b'}, 2: []byte{'b', 'a'}, 3: []byte{'c', 'd'}}, 6, foobar{11}, []int{}}
	bar := Bar{}
	err = c.Convert(&bar, &foo)
	if err != nil {
		panic(err)
	}
	err = c.Convert(&bar, nil)
	if err != nil {
		panic(err)
	}

	if expected := `{"A":[1,2,3],"B":{"1":"ab","2":"ba","3":"cd"},"C":6,"D":11,"E":[]}`; !reflect.DeepEqual(expected, jsonEncode(bar)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(bar))
	}

}

func TestNewConverterOption(t *testing.T) {
	type Foo struct {
		A  string
		B1 string
	}
	type Bar struct {
		A string
		B string
	}

	type FooBar struct {
		D Foo
		C string
	}

	type BarFoo struct {
		D Bar
		C string
	}

	fooBar := FooBar{
		D: Foo{
			A:  "a",
			B1: "b",
		},
		C: "c",
	}

	option := &StructOption{
		BannedFields: []string{"C"},
		AliasFields:  map[string]string{"D.B": "B1"},
	}

	c, err := NewConverterOption(BarFoo{}, FooBar{}, option)
	if err != nil {
		panic(err)
	}
	var barFoo BarFoo
	if err = c.Convert(&barFoo, &fooBar); err != nil {
		panic(err)
	}
	expected := BarFoo{
		D: Bar{
			A: "a",
			B: "b",
		},
	}
	if expected != barFoo {
		t.Fatalf("[expected:%v] [actual:%v]", expected, barFoo)
	}

	c, err = NewConverterOption([]BarFoo{}, []FooBar{}, option)
	if err != nil {
		panic(err)
	}
	a := []FooBar{fooBar}
	b := make([]BarFoo, 0)
	if err = c.Convert(&b, &a); err != nil {
		panic(err)
	}
	expected1 := []BarFoo{expected}
	if !reflect.DeepEqual(expected1, b) {
		t.Fatalf("[expected:%v] [actual:%v]", expected1, b)
	}

	c, err = NewConverterOption(map[int]BarFoo{}, map[int]FooBar{}, option)
	if err != nil {
		panic(err)
	}
	m := map[int]FooBar{1: fooBar}
	n := make(map[int]BarFoo)
	if err = c.Convert(&n, &m); err != nil {
		panic(err)
	}
	expected2 := map[int]BarFoo{1: expected}
	if !reflect.DeepEqual(expected2, n) {
		t.Fatalf("[expected:%v] [actual:%v]", expected2, n)

	}
}

func jsonEncode(s interface{}) string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
