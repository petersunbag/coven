package coven

import (
	"fmt"
	"reflect"
	"sync"
)

type convertType struct {
	dstTyp reflect.Type
	srcTyp reflect.Type
}

var (
	createdConvertersMu sync.Mutex
	createdConverters   = make(map[convertType]converter)
)

type converter interface {
	Convert(dst, src interface{})
}

func NewConverter(dst, src interface{}) converter {
	dstTyp := reflect.TypeOf(dst)
	srcTyp := reflect.TypeOf(src)

	if c := newConverter(dstTyp, srcTyp, true); c == nil {
		panic(fmt.Sprintf("can't convert source type %s to target type %s", srcTyp, dstTyp))
	} else {
		return c
	}
}

func newConverter(dstTyp, srcTyp reflect.Type, lock bool) converter {
	if lock {
		createdConvertersMu.Lock()
		defer createdConvertersMu.Unlock()
	}

	dstTyp = dereferencedType(dstTyp)
	srcTyp = dereferencedType(srcTyp)

	cTyp := &convertType{dstTyp, srcTyp}
	if c, ok := createdConverters[*cTyp]; ok {
		return c
	}

	var c converter
	switch sk, dk := srcTyp.Kind(), dstTyp.Kind(); {

	case srcTyp.ConvertibleTo(dstTyp):
		c = newGeneralConverter(cTyp)

	case sk == reflect.Struct && dk == reflect.Struct:
		c = newStructConverter(cTyp)

	case sk == reflect.Slice && dk == reflect.Slice:
		c = newSliceConverter(cTyp)

	default:
		c = nil
	}

	if c != nil {
		c = &delegateConverter{cTyp, c}
		createdConverters[*cTyp] = c
	}

	return c
}

type delegateConverter struct {
	*convertType
	converter
}

func (d *delegateConverter) Convert(dst, src interface{}) {
	dv := dereferencedValue(dst)
	if !dv.CanSet() {
		panic(fmt.Sprintf("target should be a pointer. [actual:%v]", dv.Type()))
	}

	if dv.Type() != d.dstTyp {
		panic(fmt.Sprintf("invalid target type. [expected:%v] [actual:%v]", d.dstTyp, dv.Type()))
	}

	sv := dereferencedValue(src)
	if !sv.CanAddr() {
		panic(fmt.Sprintf("source should be a pointer. [actual:%v]", dv.Type()))
	}

	if sv.Type() != d.srcTyp {
		panic(fmt.Sprintf("invalid source type. [expected:%v] [actual:%v]", d.srcTyp, sv.Type()))
	}

	d.converter.Convert(dv.Addr().Interface(), sv.Addr().Interface())
}
