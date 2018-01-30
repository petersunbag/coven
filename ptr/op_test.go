package ptr

import (
	"testing"
	"unsafe"
)

func TestBasicPtrCvt(t *testing.T) {
	a := int64(-1)
	b := float32(-2.9)
	c := true
	d := ""
	e := complex128(1 + 2i)

	cvtIntFloat32(unsafe.Pointer(&a), unsafe.Pointer(&b))
	if b != -1 {
		t.Fatalf("[expected:%v] [actual:%v]", -1, b)
	}

	b = -2.9
	cvtFloat32Int64(unsafe.Pointer(&b), unsafe.Pointer(&a))
	if a != -2 {
		t.Fatalf("[expected:%v] [actual:%v]", -2, a)
	}

	cvtBoolString(unsafe.Pointer(&c), unsafe.Pointer(&d))
	if d != "true" {
		t.Fatalf("[expected:%v] [actual:%v]", "true", d)
	}

	cvtComplex128String(unsafe.Pointer(&e), unsafe.Pointer(&d))
	if d != "(1+2i)" {
		t.Fatalf("[expected:%v] [actual:%v]", "(1+2i)", d)
	}

}
