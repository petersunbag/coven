package coven

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestConvertOption(t *testing.T) {
	option := &StructOption{
		BannedFields: []string{"C"},
		AliasFields:  map[string]string{"D.B": "B1"},
	}
	oo := option.convert()
	expected := &structOption{
		BannedFields: map[string]struct{}{
			"C": {},
		},
		AliasFields: make(map[string]string),
		NestedOption: map[string]*structOption{
			"D": {
				AliasFields: map[string]string{
					"B": "B1",
				},
				NestedOption: map[string]*structOption{},
			},
		},
	}
	if !reflect.DeepEqual(expected, oo) {
		t.Fatalf("\n[expected:%v]\n[actual:%v]", jsonEncode(expected), jsonEncode(oo))
	}

	// make sure convert() is idempotent
	oo = option.convert()
	if !reflect.DeepEqual(expected, oo) {
		t.Fatalf("\n[expected:%v]\n[actual:%v]", jsonEncode(expected), jsonEncode(oo))
	}
}

func TestOptionParse(t *testing.T) {
	o := structOption{
		BannedFields: map[string]struct{}{
			"A":   {},
			"A.A": {},
			"B":   {},
			"B.B": {},
		},
		AliasFields: map[string]string{
			"A":   "a",
			"A.A": "a",
			"C":   "c",
			"C.C": "c",
		},
	}

	o.parse()

	expected := structOption{
		BannedFields: map[string]struct{}{
			"A": {},
			"B": {},
		},
		AliasFields: map[string]string{
			"A": "a",
			"C": "c",
		},
		NestedOption: map[string]*structOption{
			"A": {
				BannedFields: map[string]struct{}{
					"A": {},
				},
				AliasFields: map[string]string{
					"A": "a",
				},
				NestedOption: make(map[string]*structOption),
			},
			"B": {
				BannedFields: map[string]struct{}{
					"B": {},
				},
				NestedOption: make(map[string]*structOption),
			},
			"C": {
				AliasFields: map[string]string{
					"C": "c",
				},
				NestedOption: make(map[string]*structOption),
			},
		},
	}

	if !reflect.DeepEqual(expected, o) {
		t.Fatalf("\n[expected:%v]\n[actual:%v]", jsonEncode(expected), jsonEncode(o))
	}
}

func TestOptionConvert(t *testing.T) {
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
	barFoo := BarFoo{}

	c := newStructConverter(&convertType{reflect.TypeOf(BarFoo{}), reflect.TypeOf(FooBar{}), nil})
	c.convert(unsafe.Pointer(dereferencedValue(&barFoo).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&fooBar).UnsafeAddr()))
	if expected := `{"D":{"A":"a","B":""},"C":"c"}`; expected != jsonEncode(barFoo) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(barFoo))
	}

	option := &structOption{
		BannedFields: map[string]struct{}{"C": {}},
		AliasFields:  map[string]string{"D.B": "B1"},
	}
	option.parse()
	barFoo = BarFoo{}

	c = newStructConverter(&convertType{reflect.TypeOf(BarFoo{}), reflect.TypeOf(FooBar{}), option})
	c.convert(unsafe.Pointer(dereferencedValue(&barFoo).UnsafeAddr()), unsafe.Pointer(dereferencedValue(&fooBar).UnsafeAddr()))
	if expected := `{"D":{"A":"a","B":"b"},"C":""}`; expected != jsonEncode(barFoo) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(barFoo))
	}
}

func TestOptionSliceConvert(t *testing.T) {
	type Foo struct {
		A  string
		B1 string
	}
	type Bar struct {
		A string
		B string
	}

	a := []Foo{{"a", "b"}, {"aa", "bb"}}
	var b []Bar

	c := newSliceConverter(&convertType{reflect.TypeOf([]Bar{}), reflect.TypeOf([]Foo{}), nil})
	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if expected := []Bar{{A: "a"}, {A: "aa"}}; !reflect.DeepEqual(expected, b) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}

	option := &structOption{
		BannedFields: map[string]struct{}{"A": {}},
		AliasFields:  map[string]string{"B": "B1"},
	}
	option.parse()
	b = make([]Bar, 0)

	c = newSliceConverter(&convertType{reflect.TypeOf([]Bar{}), reflect.TypeOf([]Foo{}), option})
	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if expected := []Bar{{B: "b"}, {B: "bb"}}; !reflect.DeepEqual(expected, b) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}
}

func TestOptionMapConvert(t *testing.T) {
	type Foo struct {
		A  string
		B1 string
	}
	type Bar struct {
		A string
		B string
	}

	a := map[int]Foo{1: {"a", "b"}, 2: {"aa", "bb"}}
	b := make(map[int]Bar)

	c := newMapConverter(&convertType{reflect.TypeOf(map[int]Bar{}), reflect.TypeOf(map[int]Foo{}), nil})

	c.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if expected := map[int]Bar{1: {A: "a"}, 2: {A: "aa"}}; !reflect.DeepEqual(expected, b) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}

	option := &structOption{
		BannedFields: map[string]struct{}{"A": {}},
		AliasFields:  map[string]string{"B": "B1"},
	}
	option.parse()

	c1 := newMapConverter(&convertType{reflect.TypeOf(map[int]Bar{}), reflect.TypeOf(map[int]Foo{}), option})
	c1.convert(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if expected := map[int]Bar{1: {B: "b"}, 2: {B: "bb"}}; !reflect.DeepEqual(expected, b) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, b)
	}
}
