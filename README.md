# Go struct converter #

Package `converter` copies and converts a source struct value to a target struct on convertible fields.
Field type of nested pointer is supported.
Nested anonymous fields are supported, but field has same name with upper level field is ignored.
## Install ##

Use `go get` to install this package.

    go get -u github.com/petersunbag/struct-converter

## Usage ##

```go
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

c := NewConverter(new(FooBar), new(BarFoo))

s := "c"
foobar := FooBar{10, &Foo{Baz{1, "b"}, "B", &s}}
barFoo := BarFoo{}
c.Convert(&foobar, &barFoo)

bytes, _ := json.Marshal(barFoo)
fmt.Println(string(bytes))

// Output:
// {"Foo":{"A":1,"B":"B","C":"c"}}
```

## License ##

This package is licensed under MIT license. See LICENSE for details.