package coven

import (
	"reflect"
	"unsafe"
)

type elemConverter struct {
	sTyp        reflect.Type
	dTyp        reflect.Type
	sDereferTyp reflect.Type
	dDereferTyp reflect.Type
	sReferDeep  int
	dReferDeep  int
	cvtOp       cvtOp
	converter   converter
}

func newElemConverter(dt, st reflect.Type) (e *elemConverter, ok bool) {
	e = &elemConverter{
		sTyp:        st,
		dTyp:        dt,
		sDereferTyp: st,
		dDereferTyp: dt,
		sReferDeep:  0,
		dReferDeep:  0,
		cvtOp:       nil,
		converter:   nil,
	}

	e.sDereferTyp, e.sReferDeep = referDeep(e.sDereferTyp)
	e.dDereferTyp, e.dReferDeep = referDeep(e.dDereferTyp)

	if e.cvtOp = cvtOps[convertKind{e.sDereferTyp.Kind(), e.dDereferTyp.Kind()}]; e.cvtOp != nil {
		ok = true
	} else if e.converter = newConverter(e.dDereferTyp, e.sDereferTyp, false); e.converter != nil {
		ok = true
	}

	return
}

// convert creates a value converted from src field and set it in dst field.
// The new value is first created as type of dDereferTyp,
// and then pointer nested for dReferDeep times to become a dTyp value.
func (f *elemConverter) convert(dv, sv reflect.Value) {
	if sv.Kind() == reflect.Ptr && sv.IsNil() {
		sv = reflect.New(f.sDereferTyp).Elem()
	} else {
		for d := 0; d < f.sReferDeep; d++ {
			sv = sv.Elem()
		}
	}

	var v reflect.Value

	v = reflect.New(f.dDereferTyp)
	f.converter.Convert(v.Interface(), sv.Addr().Interface())
	v = v.Elem()

	for t, d := f.dDereferTyp, 0; d < f.dReferDeep; d++ {
		tmp := reflect.New(t).Elem()
		tmp.Set(v)
		v = tmp.Addr()
		t = reflect.PtrTo(t)
	}

	dv.Set(v)
}

func (f *elemConverter) convertByPtr(dPtr, sPtr unsafe.Pointer) {
	if *((**int)(sPtr)) == nil {
		sPtr = newValue(f.sDereferTyp.Kind())
	} else {
		for d := 0; d < f.sReferDeep; d++ {
			sPtr = unsafe.Pointer(*((**int)(sPtr)))
		}
	}

	if f.dReferDeep > 0 {
		v := newValue(f.dDereferTyp.Kind())
		f.cvtOp(sPtr, v)
		for d := 0; d < f.dReferDeep; d++ {
			tmp := v
			v = unsafe.Pointer(&tmp)
		}
		*((*int)(dPtr)) = *(*int)(v)
	} else {
		sPtr := unsafe.Pointer(sPtr)
		dPtr := unsafe.Pointer(dPtr)
		f.cvtOp(sPtr, dPtr)
	}
}
