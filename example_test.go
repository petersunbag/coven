package coven

import (
	"encoding/json"
	"fmt"
)

func ExampleConverter_Convert() {
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

	foo := Foo{[]int{1, 2, 3}, map[int64][]byte{1: []byte{'a', 'b'}, 2: []byte{'b', 'a'}, 3: []byte{'c', 'd'}}, 6, foobar{11}}
	bar := Bar{}
	c.Convert(&bar, &foo)
	bytes, _ := json.Marshal(bar)
	fmt.Println(string(bytes))

	// Output:
	// {"A":[1,2,3],"B":{"1":"ab","2":"ba","3":"cd"},"C":6,"D":11}
}