# coven #
support struct-to-struct, slice-to-slice and map-to-map converting.
this project is inspired by https://github.com/thrift-iterator/go
* struct converting only affects target fields of the same name with source fields, the rest will remain unchanged.nested anonymous fields are supported.
* map converting only affects target with keys that source map has, the rest will remain unchanged.
* slice converting will overwrite the whole target slice.
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

## License ##

This package is licensed under MIT license. See LICENSE for details.