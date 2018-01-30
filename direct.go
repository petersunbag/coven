package coven

import (
	"reflect"
	"unsafe"
)

// directConverter handles converting among convertible basic types,
// and of the identical struct type.
type directConverter struct {
	*convertType
	cvtOp
	size uintptr
}

var intAlign = unsafe.Alignof(int(1))

func newDirectConverter(convertType *convertType) (c converter) {
	st := convertType.srcTyp
	dt := convertType.dstTyp

	if cvtOp := getCvtOp(st, dt); cvtOp != nil {
		c = &directConverter{
			convertType: convertType,
			cvtOp:       cvtOp,
		}
		return
	}

	if st == dt && st.Kind() == reflect.Struct {
		c = &directConverter{
			convertType: convertType,
			size:        st.Size(),
		}
		return
	}

	return
}

// convert assigns converted source value to target.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by delegateConverter.Convert() and elemConverter.convert()
func (g *directConverter) convert(dPtr, sPtr unsafe.Pointer) {
	if g.cvtOp != nil {
		g.cvtOp(sPtr, dPtr)
	} else { // same struct type
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

func getCvtOp(st, dt reflect.Type) cvtOp {
	sk := st.Kind()
	dk := dt.Kind()
	if cvtOp := cvtOps[convertKind{sk, dk}]; cvtOp != nil {
		return cvtOp
	}

	switch sk {
	case reflect.Slice:
		if dk == reflect.String && st.Elem().PkgPath() == "" {
			switch st.Elem().Kind() {
			case reflect.Uint8:
				return cvtBytesString
			case reflect.Int32:
				return cvtRunesString
			}
		}

	case reflect.String:
		if dk == reflect.Slice && dt.Elem().PkgPath() == "" {
			switch dt.Elem().Kind() {
			case reflect.Uint8:
				return cvtStringBytes
			case reflect.Int32:
				return cvtStringRunes
			}
		}
	}

	return nil
}

func cvtRunesString(sPtr unsafe.Pointer, dPtr unsafe.Pointer) {
	*(*string)(dPtr) = (string)(*(*[]rune)(sPtr))
}

func cvtBytesString(sPtr unsafe.Pointer, dPtr unsafe.Pointer) {
	*(*string)(dPtr) = (string)(*(*[]byte)(sPtr))
}

func cvtStringRunes(sPtr unsafe.Pointer, dPtr unsafe.Pointer) {
	*(*[]rune)(dPtr) = ([]rune)(*(*string)(sPtr))
}

func cvtStringBytes(sPtr unsafe.Pointer, dPtr unsafe.Pointer) {
	*(*[]byte)(dPtr) = ([]byte)(*(*string)(sPtr))
}
