package coven

import (
	"reflect"
	"unsafe"
)

type generalConverter struct {
	*convertType
	cvtOp
}

func newGeneralConverter(convertType *convertType) (c converter) {
	c = &generalConverter{
		convertType,
		cvtOps[convertKind{convertType.srcTyp.Kind(), convertType.dstTyp.Kind()}],
	}
	return
}

func (g *generalConverter) Convert(dst, src interface{}) {
	dv := dereferencedValue(dst)
	sv := dereferencedValue(src)

	if g.cvtOp != nil {
		sPtr := unsafe.Pointer(sv.UnsafeAddr())
		dPtr := unsafe.Pointer(dv.UnsafeAddr())
		g.cvtOp(sPtr, dPtr)
	} else {
		var v reflect.Value
		v = sv.Convert(dv.Type())
		dv.Set(v)
	}
}
