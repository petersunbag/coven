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

	_ = NewConverter(new(foo), new(bar))
	_ = NewConverter(new(foo), new(bar))

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
	}

	type Bar struct {
		A []*int
		B map[string]*string
		C *byte
		D int64
	}

	c := NewConverter(Bar{}, Foo{})

	foo := Foo{[]int{1, 2, 3}, map[int64][]byte{1: []byte{'a', 'b'}, 2: []byte{'b', 'a'}, 3: []byte{'c', 'd'}}, 6, foobar{11}}
	bar := Bar{}
	c.Convert(&bar, &foo)

	if expected := `{"A":[1,2,3],"B":{"1":"ab","2":"ba","3":"cd"},"C":6,"D":11}`; !reflect.DeepEqual(expected, jsonEncode(bar)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(bar))
	}

}

func jsonEncode(s interface{}) string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
