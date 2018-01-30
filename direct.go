package coven

import (
	"github.com/petersunbag/coven/ptr"
	"unsafe"
)

// directConverter handles converting among convertible basic types
type directConverter struct {
	*convertType
	cvtOp ptr.CvtOp
}

func newDirectConverter(convertType *convertType) (c converter) {
	st := convertType.srcTyp
	dt := convertType.dstTyp

	if cvtOp := ptr.GetCvtOp(st, dt); cvtOp != nil {
		c = &directConverter{
			convertType: convertType,
			cvtOp:       cvtOp,
		}
	}

	return
}

// convert assigns converted source value to target.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by delegateConverter.Convert() and elemConverter.convert()
func (g *directConverter) convert(dPtr, sPtr unsafe.Pointer) {
	g.cvtOp(sPtr, dPtr)
}
