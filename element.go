package coven

import (
	"github.com/petersunbag/coven/ptr"
	"reflect"
	"unsafe"
)

// elemConverter is converter for struct field, slice element and map key&value.
// it deals with the nested pointers and store the converter of dereferenced types.
type elemConverter struct {
	sTyp                reflect.Type
	dTyp                reflect.Type
	sDereferTyp         reflect.Type
	dDereferTyp         reflect.Type
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

	if converter := newConverter(e.dDereferTyp, e.sDereferTyp, false); converter != nil {
		e.converter = converter
		e.sEmptyDereferValPtr = newValuePtr(e.sDereferTyp)
		ok = true
	}

	return
}

// convert accepts non-nil dPtr and sPtr pointer,
// it is assured by structConverter, sliceConverter and mapConverter
func (e *elemConverter) convert(dPtr, sPtr unsafe.Pointer) {
	for d := 0; d < e.sReferDeep; d++ {
		sPtr = unsafe.Pointer(*((**int)(sPtr)))
		if sPtr == nil {
			sPtr = e.sEmptyDereferValPtr
			break
		}
	}

	deep := 0
	for ; deep < e.dReferDeep; deep++ {
		oldPtr := dPtr
		dPtr = unsafe.Pointer(*((**int)(dPtr)))
		if dPtr == nil {
			dPtr = oldPtr
			break
		}
	}

	if deep := e.dReferDeep - deep; deep > 0 {
		v := newValuePtr(e.dDereferTyp)
		e.converter.convert(v, sPtr)
		for d := 0; d < deep; d++ {
			tmp := v
			v = unsafe.Pointer(&tmp)
		}
		*(**int)(dPtr) = *(**int)(v)
	} else {
		e.converter.convert(dPtr, sPtr)
	}
}

func referDeep(t reflect.Type) (reflect.Type, int) {
	d := 0
	for k := t.Kind(); k == reflect.Ptr; k = t.Kind() {
		t = t.Elem()
		d += 1
	}
	return t, d
}

func newValuePtr(t reflect.Type) unsafe.Pointer {
	var v unsafe.Pointer
	if v = ptr.NewValuePtr(t.Kind()); v == nil {
		v = unsafe.Pointer(reflect.New(t).Elem().UnsafeAddr())
	}
	return v
}
