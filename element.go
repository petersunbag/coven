package coven

import (
	"reflect"
	"unsafe"
)

type elemConverter struct {
	sTyp                reflect.Type
	dTyp                reflect.Type
	sDereferTyp         reflect.Type
	dDereferTyp         reflect.Type
	sDereferSize        uintptr
	dDereferSize        uintptr
	sReferDeep          int
	dReferDeep          int
	sEmptyDereferValPtr unsafe.Pointer
	converter           converter
}

func newElemConverter(dt, st reflect.Type) (e *elemConverter, ok bool) {
	e = &elemConverter{
		sTyp:        st,
		dTyp:        dt,
		sDereferTyp: st,
		dDereferTyp: dt,
		sReferDeep:  0,
		dReferDeep:  0,
		converter:   nil,
	}

	e.sDereferTyp, e.sReferDeep = referDeep(e.sDereferTyp)
	e.dDereferTyp, e.dReferDeep = referDeep(e.dDereferTyp)
	e.sDereferSize = e.sDereferTyp.Size()
	e.dDereferSize = e.dDereferTyp.Size()

	if converter := newConverter(e.dDereferTyp, e.sDereferTyp, false); converter != nil {
		e.converter = converter
		e.sEmptyDereferValPtr = newValuePtr(e.sDereferTyp)
		ok = true
	}

	return
}

func (e *elemConverter) convert(dPtr, sPtr unsafe.Pointer) {
	for d := 0; d < e.sReferDeep && sPtr != nil; d++ {
		sPtr = unsafe.Pointer(*((**int)(sPtr)))
	}

	if sPtr == nil {
		sPtr = e.sEmptyDereferValPtr
	}

	if e.dReferDeep > 0 {
		v := newValuePtr(e.dDereferTyp)
		e.converter.convert(v, sPtr)
		for d := 0; d < e.dReferDeep; d++ {
			tmp := v
			v = unsafe.Pointer(&tmp)
		}
		*(**int)(dPtr) = *(**int)(v)
	} else {
		e.converter.convert(dPtr, sPtr)
	}
}

func newValuePtr(t reflect.Type) unsafe.Pointer {
	var v unsafe.Pointer
	if v = newBasicValuePtr(t.Kind()); v == nil {
		v = unsafe.Pointer(reflect.New(t).Elem().UnsafeAddr())
	}
	return v
}
