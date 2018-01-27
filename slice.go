package coven

import (
	"reflect"
	"unsafe"
)

type sliceConverter struct {
	*convertType
	*elemConverter
}

func newSliceConverter(convertType *convertType) converter {
	if elemConverter, ok := newElemConverter(convertType.dstTyp.Elem(), convertType.srcTyp.Elem()); ok {
		s := &sliceConverter{
			convertType,
			elemConverter,
		}
		return s
	}
	return nil
}

func (s *sliceConverter) Convert(dst, src interface{}) {
	dv := dereferencedValue(dst)
	sv := dereferencedValue(src)

	dSlice := (*sliceHeader)(unsafe.Pointer(dv.UnsafeAddr()))
	sSlice := (*sliceHeader)(unsafe.Pointer(sv.UnsafeAddr()))

	dOffset, sOffset := uintptr(0), uintptr(0)
	length := sv.Len()
	dSlice.Len = length

	if dSlice.Cap < length {
		newVal := reflect.MakeSlice(s.dstTyp, 0, length)
		dSlice.Data = unsafe.Pointer(newVal.Pointer())
		dSlice.Cap = length
	}

	for i := 0; i < length; i++ {
		if s.elemConverter.cvtOp != nil {
			dElemPtr := unsafe.Pointer(uintptr(dSlice.Data) + dOffset)
			sElemPtr := unsafe.Pointer(uintptr(sSlice.Data) + sOffset)
			s.elemConverter.convertByPtr(dElemPtr, sElemPtr)
			dOffset += s.elemConverter.dDereferTyp.Size()
			sOffset += s.elemConverter.sDereferTyp.Size()
		} else {
			sElemV := sv.Index(i)
			dElemV := dv.Index(i)
			s.elemConverter.convert(dElemV, sElemV)
		}
	}
}
