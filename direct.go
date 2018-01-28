package coven

import (
	"reflect"
	"unsafe"
)

type directConverter struct {
	*convertType
	cvtOp
	size uintptr
}

var intAlign = unsafe.Alignof(int(1))

func newGeneralConverter(convertType *convertType) (c converter) {
	st := convertType.srcTyp
	dt := convertType.dstTyp
	sk := st.Kind()
	dk := dt.Kind()
	if sk == reflect.Slice || sk == reflect.Map {
		return
	}

	if cvtOp := cvtOps[convertKind{sk, dk}]; cvtOp != nil {
		c = &directConverter{
			convertType: convertType,
			cvtOp:       cvtOp,
		}
		return
	}
	if st == dt && sk == reflect.Struct {
		c = &directConverter{
			convertType: convertType,
			size:        st.Size(),
		}
		return
	}

	return
}

//dv and sv must be dereferened value
func (g *directConverter) convert(dPtr, sPtr unsafe.Pointer) {
	if g.cvtOp != nil {
		g.cvtOp(sPtr, dPtr)
	} else { //dst and src have same type, exclude map and slice
		size := g.size
		align := uintptr(0)
		for ; align < size-intAlign; align += intAlign {
			*(*int)(unsafe.Pointer(uintptr(dPtr) + align)) = *(*int)(unsafe.Pointer(uintptr(sPtr) + align))
		}
		for ; align < size; align++ {
			*(*byte)(unsafe.Pointer(uintptr(dPtr) + align)) = *(*byte)(unsafe.Pointer(uintptr(sPtr) + align))
		}
	}
}
