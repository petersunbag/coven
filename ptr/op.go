package ptr

import (
	"reflect"
	"unsafe"
)

var (
	intAlign = unsafe.Alignof(int(1))
	cvtOps = make(map[convertKind]CvtOp)
)

type convertKind struct {
	srcTyp reflect.Kind
	dstTyp reflect.Kind
}

type CvtOp func(unsafe.Pointer, unsafe.Pointer)

func GetCvtOp(st, dt reflect.Type) CvtOp {
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

func Copy(dPtr, sPtr unsafe.Pointer, size uintptr) {
	align := uintptr(0)
	for ; align < size-intAlign; align += intAlign {
		*(*int)(unsafe.Pointer(uintptr(dPtr) + align)) = *(*int)(unsafe.Pointer(uintptr(sPtr) + align))
	}
	for ; align < size; align++ {
		*(*byte)(unsafe.Pointer(uintptr(dPtr) + align)) = *(*byte)(unsafe.Pointer(uintptr(sPtr) + align))
	}
}
