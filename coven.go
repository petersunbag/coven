package coven

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var (
	createdConvertersMu sync.Mutex
	createdConverters   = make(map[convertType]*Converter)
)

type convertType struct {
	dstTyp reflect.Type
	srcTyp reflect.Type
}

// converter can handle converting among convertible basic types,
// and struct-struct, slice-slice, map-map converting too.
// type with nested pointer is supported.

// all methods in converter are thread-safe.
// we can define a global variable to hold a converter and use it in any goroutine.
type converter interface {
	convert(dPtr, sPtr unsafe.Pointer)
}

type Converter struct {
	*convertType
	converter
}

func (d *Converter) Convert(dst, src interface{}) error {
	if dst == nil || src == nil || reflect.ValueOf(dst).IsNil() || reflect.ValueOf(src).IsNil() {
		return nil
	}

	dv := dereferencedValue(dst)
	if !dv.CanSet() {
		return errors.New(fmt.Sprintf("[coven]destination should be a pointer. [actual:%v]", dv.Type()))
	}

	if dv.Type() != d.dstTyp {
		return errors.New(fmt.Sprintf("[coven]invalid destination type. [expected:%v] [actual:%v]", d.dstTyp, dv.Type()))
	}

	sv := dereferencedValue(src)
	if !sv.CanAddr() {
		return errors.New(fmt.Sprintf("[coven]source should be a pointer. [actual:%v]", sv.Type()))
	}

	if sv.Type() != d.srcTyp {
		return errors.New(fmt.Sprintf("[coven]invalid source type. [expected:%v] [actual:%v]", d.srcTyp, sv.Type()))
	}

	d.converter.convert(unsafe.Pointer(dv.UnsafeAddr()), unsafe.Pointer(sv.UnsafeAddr()))
	return nil
}

func NewConverter(dst, src interface{}) (*Converter, error) {
	dstTyp := reflect.TypeOf(dst)
	srcTyp := reflect.TypeOf(src)

	if c := newConverter(dstTyp, srcTyp, true); c == nil {
		return nil, errors.New(fmt.Sprintf("can't convert source type %s to destination type %s", srcTyp, dstTyp))
	} else {
		return c, nil
	}
}

func newConverter(dstTyp, srcTyp reflect.Type, lock bool) *Converter {
	if lock {
		createdConvertersMu.Lock()
		defer createdConvertersMu.Unlock()
	}

	dstTyp = dereferencedType(dstTyp)
	srcTyp = dereferencedType(srcTyp)

	cTyp := &convertType{dstTyp, srcTyp}
	if dc, ok := createdConverters[*cTyp]; ok {
		return dc
	}

	var c converter
	if c = newBasicConverter(cTyp); c == nil {
		switch sk, dk := srcTyp.Kind(), dstTyp.Kind(); {

		case sk == reflect.Struct && dk == reflect.Struct:
			c = newStructConverter(cTyp)

		case sk == reflect.Slice && dk == reflect.Slice:
			c = newSliceConverter(cTyp)

		case sk == reflect.Map && dk == reflect.Map:
			c = newMapConverter(cTyp)
		}
	}
	if c != nil {
		dc := &Converter{cTyp, c}
		createdConverters[*cTyp] = dc
		return dc
	}

	return nil
}
