package coven

import (
	"github.com/petersunbag/coven/ptr"
	"unsafe"
)

// directConverter handles converting among convertible basic types
type basicConverter struct {
	*convertType
	cvtOp ptr.CvtOp
}

func newBasicConverter(convertType *convertType) (c converter) {
	if cvtOp := ptr.GetCvtOp(convertType.srcTyp, convertType.dstTyp); cvtOp != nil {
		c = &basicConverter{
			convertType: convertType,
			cvtOp:       cvtOp,
		}
	}

	return
}

// convert assigns converted source value to destination.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by Converter.Convert() and elemConverter.convert()
func (g *basicConverter) convert(dPtr, sPtr unsafe.Pointer) {
	g.cvtOp(sPtr, dPtr)
}
