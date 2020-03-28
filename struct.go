package coven

import (
	"github.com/petersunbag/coven/ptr"
	"reflect"
	"unicode"
	"unsafe"
)

type structConverter struct {
	*convertType
	fieldConverters []*fieldConverter
	size            uintptr
}

// NewStructConverter finds convertible fields of the same name in convertType,
// and stores fieldConverters in structConverter, including nested anonymous fields.
func newStructConverter(convertType *convertType) (s converter) {
	if convertType.srcTyp == convertType.dstTyp {
		s = &structConverter{
			convertType: convertType,
			size:        convertType.srcTyp.Size(),
		}
		return
	}
	_, sFields := extractFields(convertType.srcTyp, 0)
	dFieldIndex, _ := extractFields(convertType.dstTyp, 0)
	if convertType.option != nil {
		dFieldIndex = filterField(dFieldIndex, convertType.option.BannedFields)
		dFieldIndex = aliasField(dFieldIndex, convertType.option.AliasFields)
	}
	fieldConverters := make([]*fieldConverter, 0, len(dFieldIndex))

	for _, df := range dFieldIndex {
		if sf, ok := sFields[df.Name]; ok {
			var nestOption *structOption = nil
			if convertType.option != nil && convertType.option.NestedOption != nil {
				nestOption = convertType.option.NestedOption[df.Name]
			}
			if fieldConverter := newFieldConverter(*df, *sf, nestOption); fieldConverter != nil {
				fieldConverters = append(fieldConverters, fieldConverter)
			}
		}
	}

	if len(fieldConverters) > 0 {
		s = &structConverter{
			convertType:     convertType,
			fieldConverters: fieldConverters,
		}
	}

	return
}

func filterField(fields []*reflect.StructField, bannedFields map[string]struct{}) []*reflect.StructField {
	if bannedFields == nil {
		return fields
	}

	filtered := make([]*reflect.StructField, 0, len(fields))
	for _, f := range fields {
		if _, ok := bannedFields[f.Name]; ok {
			continue
		} else {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func aliasField(fields []*reflect.StructField, aliasFields map[string]string) []*reflect.StructField {
	if aliasFields == nil {
		return fields
	}

	for _, f := range fields {
		if alias, ok := aliasFields[f.Name]; ok {
			f.Name = alias
		}
	}
	return fields
}

// convert only affects fields stored in fieldConverters, the rest will remain unchanged.
// dPtr and sPtr must pointed to a non-pointer value,
// it is assured by Converter.Convert() and elemConverter.convert()
func (s *structConverter) convert(dPtr, sPtr unsafe.Pointer) {
	if s.dstTyp == s.srcTyp {
		ptr.Copy(dPtr, sPtr, s.size)
		return
	}

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
	sName   string
	dName   string
}

func newFieldConverter(df, sf reflect.StructField, option *structOption) (f *fieldConverter) {
	if elemConverter, ok := newElemConverter(df.Type, sf.Type, option); ok {
		return &fieldConverter{
			elemConverter: elemConverter,
			sOffset:       sf.Offset,
			dOffset:       df.Offset,
			sName:         sf.Name,
			dName:         df.Name,
		}
	}

	return nil
}

// extractFields returns every exported field of a struct, including nested anonymous fields.
// Field has same name with upper level field is not returned.
func extractFields(t reflect.Type, offset uintptr) (fieldSlice []*reflect.StructField, fieldMap map[string]*reflect.StructField) {
	fieldMap = make(map[string]*reflect.StructField)
	anonymous := make([]*reflect.StructField, 0, t.NumField())
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		f.Offset = f.Offset + offset
		if unicode.IsUpper(rune(f.Name[0])) {
			fieldSlice = append(fieldSlice, &f)
			fieldMap[f.Name] = &f
		}

		if f.Anonymous {
			anonymous = append(anonymous, &f)
		}
	}

	for _, af := range anonymous {
		afTyp := dereferencedType(af.Type)
		s, _ := extractFields(afTyp, af.Offset)
		for _, f := range s {
			name := f.Name
			if _, ok := fieldMap[name]; ok {
				continue
			}
			fieldSlice = append(fieldSlice, f)
			fieldMap[f.Name] = f
		}
	}

	return
}
