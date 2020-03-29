# coven #

[![Build Status](https://travis-ci.org/petersunbag/coven.svg?branch=master)](https://travis-ci.org/petersunbag/coven)
[![Coverage Status](https://coveralls.io/repos/github/petersunbag/coven/badge.svg?branch=master&98)](https://coveralls.io/github/petersunbag/coven?branch=master)

Support struct-to-struct, slice-to-slice and map-to-map converting.  
This package is inspired by https://github.com/thrift-iterator/go
* struct converting only affects destination fields of the same name with source fields, the rest will remain unchanged.nested anonymous fields are supported.
* map converting only affects destination map with keys that source map has, the rest will remain unchanged.
* slice converting will overwrite the whole destination slice.
* type with nested pointers is supported.
* except for map converting, use unsafe.pointer instead of reflect.Value to convert, so it can convert faster.
## Install ##

Use `go get` to install this package.

    go get -u github.com/petersunbag/coven

## Usage ##
### Basic usage ###
```go
type foobar struct {
    D int
}
type Foo struct {
    A []int
    B map[int64][]byte
    C byte
    foobar
}

type Bar struct {
    A []*int
    B map[string]*string
    C *byte
    D int64
}

var c, err = NewConverter(Bar{}, Foo{})
if err != nil {
    panic(err)
}

func demo(){
    foo := Foo{[]int{1, 2, 3}, map[int64][]byte{1: []byte{'a', 'b'}, 2: []byte{'b', 'a'}, 3: []byte{'c', 'd'}}, 6, foobar{11}}
    bar := Bar{}
    err := c.Convert(&bar, &foo)
    if err != nil {
        panic(err)
    }
    bytes, _ := json.Marshal(bar)
    fmt.Println(string(bytes))
}

// Output:
// {"A":[1,2,3],"B":{"1":"ab","2":"ba","3":"cd"},"C":6,"D":11}
```
### Use `StructOption` to control struct converting behavior ###
To achieve a more flexible struct converting, you can use `StructOption` to control the behavior of a field in dst struct. An option describes if a filed is allowed to be converted or if it should use a alias when converting.
```go
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

option := &StructOption{
    BannedFields: []string{"C"},
    AliasFields:  map[string]string{"D.B": "B1"},
}

c, err := NewConverterOption(BarFoo{}, FooBar{}, option)
if err != nil {
    panic(err)
}

func demo(){
    fooBar := FooBar{
        D: Foo{
            A:  "a",
            B1: "b",
        },
        C: "c",
    }
	var barFoo BarFoo
    if err = c.Convert(&barFoo, &fooBar); err != nil {
        panic(err)
    }
    bytes, _ := json.Marshal(barFoo)
    fmt.Println(string(bytes))
}

// Output:
// {"D":{"A":"a","B":"b"},"C":""}
```
## Benchmark ##

Direct ptr operation is faster than using reflection for basic types, such as int to float conversion, and even for identical slice or struct.

| ptr int-float | reflection int-float | ptr struct  | reflection struct | ptr slice  | reflection slice |
| ---           | ---                  | ---         | ---               | ---        | ---              |
| 57.9 ns/op    | 98.6 ns/op           | 86.9 ns/op  | 118 ns/op         | 80.2 ns/op | 99.7 ns/op       |
| 0 B/op        | 16 B/op              | 0 B/op      | 64 B/op           | 0 B/op     | 32 B/op          |
| 0 allocs/op   | 2 allocs/op          | 0 allocs/op | 1 allocs/op       | 0 allocs/op| 1 allocs/op      |

Test cases above don't include map type, because there is no way to operate map through ptr in Go.  
See benchmark_test.go for details.
## FAQ ##
### Why not use tag instead of `StructOption`?  ###
Yes, using tag may seem more convenient and intuitive. But tag has its limitation.
 You can only use tag in your OWN struct. What if the dst struct is not user 
 defined? It may come from a third-party lib, and you can't add tag on it.

`coven` is just a converter, which may not expected that dst struct is user defined.
That is why I choose to use an extra `StructOption`.
## License ##

This package is licensed under MIT license. See LICENSE for details.