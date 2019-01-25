package coven

import (
	"reflect"
	"testing"
)

func BenchmarkCoven(b *testing.B) {
	type foobar struct {
		D int
	}
	type Foo struct {
		A []int
		B map[int64]int
		C byte
		foobar
	}

	type Bar struct {
		A []*int
		B map[string]*string
		C *byte
		D int64
	}

	c, err := NewConverter(Bar{}, Foo{})
	if err != nil {
		panic(err)
	}

	foo := Foo{[]int{1, 2, 3}, map[int64]int{1: 1, 2: 2, 3: 3}, 6, foobar{11}}
	bar := Bar{}

	for i := 0; i < b.N; i++ {
		err = c.Convert(&bar, &foo)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkCovenWithoutMap(b *testing.B) {
	type foobar struct {
		D int
	}
	type Foo struct {
		A []int
		//B map[int64][]byte
		C byte
		foobar
	}

	type Bar struct {
		A []*int
		//B map[string]*string
		C *byte
		D int64
	}

	c, err := NewConverter(Bar{}, Foo{})
	if err != nil {
		panic(err)
	}

	foo := Foo{[]int{1, 2, 3}, 6, foobar{11}}
	bar := Bar{}

	for i := 0; i < b.N; i++ {
		err = c.Convert(&bar, &foo)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkStructConvert(b *testing.B) {
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

	foobar := FooBar{10, &Foo{Baz{1, "b"}, "B", stringPtr("c")}}
	barFoo := BarFoo{}

	c, err := NewConverter(BarFoo{}, FooBar{})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		err = c.Convert(&barFoo, &foobar)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSameStruct(b *testing.B) {
	type bar struct {
		A int
		B byte
	}

	type foo struct {
		bar
		C string
		D []int
	}

	Foo := foo{}

	c, err := NewConverter(Foo, Foo)
	if err != nil {
		panic(err)
	}
	foo1 := &foo{bar{1, 2}, "abc", []int{1, 2, 3}}
	foo2 := foo{}

	for i := 0; i < b.N; i++ {
		err = c.Convert(&foo2, &foo1)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSameStructReflect(b *testing.B) {
	type bar struct {
		A int
		B byte
	}

	type foo struct {
		bar
		C string
		D []int
	}

	foo1 := foo{bar{1, 2}, "abc", []int{1, 2, 3}}
	foo2 := foo{}
	t := reflect.TypeOf(foo1)

	for i := 0; i < b.N; i++ {
		reflect.ValueOf(&foo2).Elem().Set(reflect.ValueOf(foo1).Convert(t))
	}
}

func BenchmarkSameSlice(b *testing.B) {
	a := []int{1, 2, 3}
	d := []int{4}

	c, err := NewConverter([]int{}, []int{})
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		err = c.Convert(&d, &a)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSameSliceReflect(b *testing.B) {
	a := []int{1, 2, 3}
	d := []int{4}

	t := reflect.TypeOf(d)

	for i := 0; i < b.N; i++ {
		reflect.ValueOf(&d).Elem().Set(reflect.ValueOf(a).Convert(t))
	}
}

func BenchmarkBasic(b *testing.B) {
	x := 1
	y := 2.2
	c, err := NewConverter(y, x)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		err = c.Convert(&y, &x)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkBasicReflect(b *testing.B) {
	x := 1
	y := 2.2
	t := reflect.TypeOf(y)

	for i := 0; i < b.N; i++ {
		reflect.ValueOf(&y).Elem().Set(reflect.ValueOf(x).Convert(t))
	}
}
