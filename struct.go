package coven

import (
	"reflect"
	"unicode"
	"unsafe"
)

type structConverter struct {
	*convertType
	fieldConverters []*fieldConverter
}

// NewStructConverter finds convertible fields of the same name in convertType,
// and stores fieldConverters in structConverter, including nested anonymous fields.
func newStructConverter(convertType *convertType) (c converter) {
	dFieldIndex := fieldIndex(convertType.dstTyp, []int{})
	fieldConverters := make([]*fieldConverter, 0, len(dFieldIndex))
	for _, index := range dFieldIndex {
		df := convertType.dstTyp.FieldByIndex(index)
		df.Index = index
		if sf, ok := convertType.srcTyp.FieldByName(df.Name); ok {
			if fieldConverter := newFieldConverter(df, sf); fieldConverter != nil {
				fieldConverters = append(fieldConverters, fieldConverter)
			}
		}
	}

	if len(fieldConverters) > 0 {
		c = &structConverter{
			convertType,
			fieldConverters,
		}
	}

	return
}

// convert only affects fields stored in fieldConverters, the rest will remain unchanged.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by delegateConverter.Convert() and elemConverter.convert()
func (s *structConverter) convert(dPtr, sPtr unsafe.Pointer) {
	for _, fieldConverter := range s.fieldConverters {
		dPtr := unsafe.Pointer(uintptr(dPtr) + fieldConverter.dOffset)
		sPtr := unsafe.Pointer(uintptr(sPtr) + fieldConverter.sOffset)
		fieldConverter.convert(dPtr, sPtr)
	}
}

type fieldConverter struct {
	*elemConverter
	sOffset uintptr
	dOffset uintptr
	sIndex  []int
	dIndex  []int
	sName   string
	dName   string
}

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
