# coven #

[![Build Status](https://travis-ci.org/petersunbag/coven.svg?branch=master)](https://travis-ci.org/petersunbag/coven)
[![Coverage Status](https://coveralls.io/repos/github/petersunbag/coven/badge.svg?branch=master&98)](https://coveralls.io/github/petersunbag/coven?branch=master)

Support struct-to-struct, slice-to-slice and map-to-map converting.  
This package is inspired by https://github.com/thrift-iterator/go
* struct converting only affects destination fields of the same name with source fields, the rest will remain unchanged.nested anonymous fields are supported.
* map converting only affects destination map with keys that source map has, the rest will remain unchanged.
* slice converting will overwrite the whole destination slice.
* type with nested pointers is supported.
* except for map converting, use unsafe.pointer instead of reflect.Value to convert.
## Install ##

Use `go get` to install this package.

    go get -u github.com/petersunbag/coven

## Usage ##

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

var c = NewConverter(Bar{}, Foo{})

func demo(){
    foo := Foo{[]int{1, 2, 3}, map[int64][]byte{1: []byte{'a', 'b'}, 2: []byte{'b', 'a'}, 3: []byte{'c', 'd'}}, 6, foobar{11}}
    bar := Bar{}
    c.Convert(&bar, &foo)
    bytes, _ := json.Marshal(bar)
    fmt.Println(string(bytes))
}

// Output:
// {"A":[1,2,3],"B":{"1":"ab","2":"ba","3":"cd"},"C":6,"D":11}
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
## License ##

This package is licensed under MIT license. See LICENSE for details.