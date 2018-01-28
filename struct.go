package coven

import (
	"reflect"
	"unicode"
	"unsafe"
)

// structConverter stores convertible fields of srcTyp and dstTyp.
// Field type of nested pointer is supported.
// srcTyp and dstTyp are both dereferenced reflect.Type

// All methods in structConverter are thread-safe.
// We can define a global variable to hold a structConverter and use it in any goroutine.
type structConverter struct {
	*convertType
	fieldConverters []*fieldConverter
}

// NewStructConverter analyzes type information of src and dst
// and returns a *structConverter with convertible fields of the same name.
// Field type of nested pointer is supported.
// It panics if src or dst is not a struct.
func newStructConverter(convertType *convertType) (c converter) {
	dFieldIndex := fieldIndex(convertType.dstTyp, []int{})
	fCvts := make([]*fieldConverter, 0, len(dFieldIndex))
	for _, index := range dFieldIndex {
		df := convertType.dstTyp.FieldByIndex(index)
		df.Index = index
		if sf, ok := convertType.srcTyp.FieldByName(df.Name); ok {
			if fCvt := newFieldConverter(df, sf); fCvt != nil {
				fCvts = append(fCvts, fCvt)
			}
		}
	}

	if len(fCvts) > 0 {
		c = &structConverter{
			convertType,
			fCvts,
		}
	}

	return
}

// Convert creates field values converted from src and set them in dst on fields stored in structConverter.
// Field type of nested pointer is supported.
// dst should be a pointer to a struct, otherwise Convert panics.
// dereferenced src and dst type should match their counterparts in structConverter.

//dv and sv must be dereferened value
func (s *structConverter) convert(dPtr, sPtr unsafe.Pointer) {
	for _, fCvt := range s.fieldConverters {
		dPtr := unsafe.Pointer(uintptr(dPtr) + fCvt.dOffset)
		sPtr := unsafe.Pointer(uintptr(sPtr) + fCvt.sOffset)
		fCvt.convert(dPtr, sPtr)
	}
}

// fieldConverter stores information of convertible field from one type to another.
// sTyp and dTyp are original types of src and dst.
// sDereferTyp and dDereferTyp are dereferenced types of src and dst.
// sReferDeep and dReferDeep are levels of nested pointer of src and dst.
// If src and dst are different struct, make them a structConverter.
type fieldConverter struct {
	*elemConverter
	sOffset uintptr
	dOffset uintptr
	sIndex  []int
	dIndex  []int
	sName   string
	dName   string
}

// newFieldConverter analyzes information of src and dst field
// and returns a *fieldConverter, if they are convertible, ok is true.
// Field type of nested pointer is supported.
func newFieldConverter(df, sf reflect.StructField) (f *fieldConverter) {
	if elemConverter, ok := newElemConverter(df.Type, sf.Type); ok {
		return &fieldConverter{
			elemConverter: elemConverter,
			sOffset:       sf.Offset,
			dOffset:       df.Offset,
			sIndex:        sf.Index,
			dIndex:        df.Index,
			sName:         sf.Name,
			dName:         df.Name,
		}
	}

	return nil
}

// fieldIndex returns indices of every field in a struct, including nested anonymous fields.
// Field has same name with upper level field is not returned.
func fieldIndex(t reflect.Type, prefixIndex []int) (indices [][]int) {
	t = dereferencedType(t)
	fName := make(map[string]struct{})
	anonymous := make([]int, 0, t.NumField())
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)

		if unicode.IsUpper(rune(f.Name[0])) {
			indices = append(indices, append(prefixIndex, i))
			fName[f.Name] = struct{}{}
		}

		if f.Anonymous {
			anonymous = append(anonymous, i)
		}
	}

	for _, i := range anonymous {
		for _, index := range fieldIndex(t.Field(i).Type, []int{i}) {
			name := t.FieldByIndex(index).Name
			if _, ok := fName[name]; ok {
				continue
			}
			fName[name] = struct{}{}
			indices = append(indices, append(prefixIndex, index...))
		}
	}

	return
}
