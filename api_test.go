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
	type Foo struct {
		A []int
		B map[int64]string
		C byte
	}

	type Bar struct {
		A []*int
		B map[string]*string
		C *byte
	}

	c := NewConverter(Bar{}, Foo{})

	foo := Foo{[]int{1, 2, 3}, map[int64]string{1: "a", 2: "b", 3: "c"}, 6}
	bar := Bar{}
	c.Convert(&bar, &foo)

	if expected := `{"A":[1,2,3],"B":{"1":"a","2":"b","3":"c"},"C":6}`; !reflect.DeepEqual(expected, jsonEncode(bar)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(bar))
	}

}

func jsonEncode(s interface{}) string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
