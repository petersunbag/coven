package coven

//import (
//	"reflect"
//	"unsafe"
//	"fmt"
//	"testing"
//)
//
//func TestConvertByte(t *testing.T) {
//	//type foo struct {
//	//	//c string
//	//	d byte
//	//	e byte
//	//	f byte
//	//	g byte
//	//	h byte
//	//	i byte
//	//	j byte
//	//	k byteR
//	//	l byte
//	//	m byte
//	//	n byte
//	//	//d byte
//	//	//d byte
//	//	//a int
//	//}
//	//f := foo{1,1,1,1,1,1,1,1,1,1,1}
//	//fmt.Println(unsafe.Alignof(f))
//	//
//	////type foo struct {
//	////	a byte
//	////}
//	////f:= foo{1}
//	//b := foo{d:2}
//	//fmt.Println(uintptr(unsafe.Pointer(&b))-uintptr(unsafe.Pointer(&f)))
//	//
//	//fmt.Println(unsafe.Sizeof(f))
//	//fmt.Println(reflect.TypeOf(foo{}).Size())
//	//
//	//
//	//align := unsafe.Alignof(int(1))
//	//fmt.Println(align)
//	//
//	//size := unsafe.Sizeof(map[string]string{"a":"a", "b":"b"})
//	//fmt.Println(size)
//
//	m := map[string]string{"a":"a", "b":"b"}
//
//	n := reflect.New(reflect.TypeOf(map[string]string{})).Interface()
//
//	//inter := (*emptyInterface)(unsafe.Pointer(&n))
//	//inter.word = unsafe.Pointer(&m)
//	//realInterface := (*interface{})(unsafe.Pointer(&inter))
//	//v := reflect.ValueOf(*realInterface)
//	//fmt.Println(realInterface)
//	//fmt.Println(inter.typ)
//	//fmt.Println(inter.word)
//
//	//mapInterface := encoder.mapInterface
//	//mapInterface.word = ptr
//	//realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
//	//mapVal := reflect.ValueOf(*realInterface)
//
//	//fmt.Println(v)
//
//	ptr := unsafe.Pointer(&m)
//
//	mapInterface := *(*emptyInterface)(unsafe.Pointer(&n))
//	mapInterface.word = ptr
//	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
//	fmt.Println(reflect.TypeOf(*realInterface))
//	mapVal := reflect.ValueOf(*realInterface)
//	fmt.Println(mapVal)
//	fmt.Println(mapVal.Elem().MapKeys())
//
//
//	fmt.Println(unsafe.Sizeof([]int{1,2,3,4,5,6,7,8}))
//}

// convert creates a value converted from src field and set it in dst field.
// The new value is first created as type of dDereferTyp,
// and then pointer nested for dReferDeep times to become a dTyp value.
//func (e *elemConverter) convert(dv, sv reflect.Value) {
//	if sv.Kind() == reflect.Ptr && sv.IsNil() {
//		sv = reflect.New(e.sDereferTyp).Elem()
//	} else {
//		for d := 0; d < e.sReferDeep; d++ {
//			sv = sv.Elem()
//		}
//	}
//
//	var v reflect.Value
//	v = reflect.New(e.dDereferTyp).Elem()
//	e.converter.convert(v, sv)
//
//	for t, d := e.dDereferTyp, 0; d < e.dReferDeep; d++ {
//		tmp := reflect.New(t).Elem()
//		tmp.Set(v)
//		v = tmp.Addr()
//		t = reflect.PtrTo(t)
//	}
//
//	dv.Set(v)
//}
