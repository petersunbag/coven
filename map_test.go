package coven

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestMapConverter_Convert(t *testing.T) {
	c := newMapConverter(&convertType{reflect.TypeOf(map[string]*string{}), reflect.TypeOf(map[int]int{})})
	n := map[int]int{1: 1, 2: 2, 3: 3}
	m := map[string]*string{}
	c.convert(unsafe.Pointer(&m), unsafe.Pointer(&n))

	if expected := `{"1":"1","2":"2","3":"3"}`; !reflect.DeepEqual(expected, jsonEncode(m)) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, jsonEncode(m))
	}
}

func TestPtrToMapValue(t *testing.T) {
	m := map[string]string{"a": "a", "b": "b"}
	n := reflect.New(reflect.TypeOf(map[string]string{})).Interface()
	inter := *(*emptyInterface)(unsafe.Pointer(&n))
	v := ptrToMapValue(&inter, unsafe.Pointer(&m)).Interface()

	if expected := map[string]string{"a": "a", "b": "b"}; !reflect.DeepEqual(expected, v) {
		t.Fatalf("[expected:%v] [actual:%v]", expected, v)
	}
}
