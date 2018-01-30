package basic

import (
	"reflect"
	"unsafe"
)

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
