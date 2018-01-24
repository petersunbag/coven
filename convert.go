package struct_converter

import (
	"fmt"
	"reflect"
)

var createdConverters = make(map[string]*converter)

type converter struct {
	dstTyp          reflect.Type
	srcTyp          reflect.Type
	fieldConverters []*fieldConverter
}

func New(src interface{}, dst interface{}) (c *converter) {
	srcTyp := dereferencedType(reflect.TypeOf(src))
	dstTyp := dereferencedType(reflect.TypeOf(dst))

	key := srcTyp.String() + "-" + dstTyp.String()
	if c, ok := createdConverters[key]; ok {
		return c
	}

	if srcTyp.Kind() != reflect.Struct {
		panic("source is not a struct!")
	}

	if dstTyp.Kind() != reflect.Struct {
		panic("target is not a struct!")
	}

	dFieldIndex := fieldIndex(dstTyp, []int{})
	fCvts := make([]*fieldConverter, 0, len(dFieldIndex))
	for _, index := range dFieldIndex {
		df := dstTyp.FieldByIndex(index)
		df.Index = index
		if sf, ok := srcTyp.FieldByName(df.Name); ok {
			if fCvt, ok := newFieldConverter(sf, df); ok {
				fCvts = append(fCvts, fCvt)
			}
		}
	}

	if len(fCvts) > 0 {
		c = &converter{
			dstTyp,
			srcTyp,
			fCvts,
		}
		createdConverters[key] = c
	}

	return
}

func (c *converter) Convert(src interface{}, dst interface{}) {
	dv := dereferencedValue(dst)
	if !dv.CanSet() {
		panic(fmt.Sprintf("target should be a pointer. [actual:%v]", dv.Type()))
	}
	if dv.Type() != c.dstTyp {
		panic(fmt.Sprintf("invalid target type. [expected:%v] [actual:%v]", c.dstTyp, dv.Type()))
	}

	sv := dereferencedValue(src)
	if sv.Type() != c.srcTyp {
		panic(fmt.Sprintf("invalid source type. [expected:%v] [actual:%v]", c.srcTyp, sv.Type()))
	}

	for _, fCvt := range c.fieldConverters {
		sf := sv.FieldByIndex(fCvt.sIndex)
		df := dv.FieldByIndex(fCvt.dIndex)
		fCvt.convert(sf, df)
	}
}

type fieldConverter struct {
	sTyp        reflect.Type
	dTyp        reflect.Type
	sDereferTyp reflect.Type
	dDereferTyp reflect.Type
	sReferDeep  int
	dReferDeep  int
	sIndex      []int
	dIndex      []int
	sName       string
	dName       string
	structCvt   *converter
	//dstOffset uintptr  //todo use unsafe.pointer instead of reflect.Value
	//srcOffset uintptr
}

func newFieldConverter(sf, df reflect.StructField) (f *fieldConverter, ok bool) {
	f = &fieldConverter{
		sTyp:        sf.Type,
		dTyp:        df.Type,
		sDereferTyp: sf.Type,
		dDereferTyp: df.Type,
		sReferDeep:  0,
		dReferDeep:  0,
		sIndex:      sf.Index,
		dIndex:      df.Index,
		sName:       sf.Name,
		dName:       df.Name,
		structCvt:   nil,
	}

	for k := f.sDereferTyp.Kind(); k == reflect.Ptr; k = f.sDereferTyp.Kind() {
		f.sDereferTyp = f.sDereferTyp.Elem()
		f.sReferDeep += 1
	}

	for k := f.dDereferTyp.Kind(); k == reflect.Ptr; k = f.dDereferTyp.Kind() {
		f.dDereferTyp = f.dDereferTyp.Elem()
		f.dReferDeep += 1
	}

	if f.sDereferTyp.ConvertibleTo(f.dDereferTyp) {
		ok = true
	} else if f.sDereferTyp.Kind() == reflect.Struct && f.dDereferTyp.Kind() == reflect.Struct {
		f.structCvt = New(reflect.New(f.sDereferTyp).Interface(), reflect.New(f.dDereferTyp).Interface())
		if f.structCvt != nil {
			ok = true
		}
	}

	return
}

func (f *fieldConverter) convert(sv, dv reflect.Value) {
	if sv.Kind() == reflect.Ptr && sv.IsNil() {
		sv = reflect.New(f.sDereferTyp).Elem()
	} else {
		for d := 0; d < f.sReferDeep; d++ {
			sv = sv.Elem()
		}
	}

	var v reflect.Value
	if f.structCvt == nil {
		v = sv.Convert(f.dDereferTyp)
	} else {
		v = reflect.New(f.dDereferTyp)
		f.structCvt.Convert(sv.Interface(), v.Interface())
		v = v.Elem()
	}

	for t, d := f.dDereferTyp, 0; d < f.dReferDeep; d++ {
		tmp := reflect.New(t).Elem()
		tmp.Set(v)
		v = tmp.Addr()
		t = reflect.PtrTo(t)
	}

	dv.Set(v)
}

// DFS
func fieldIndex(t reflect.Type, prefixIndex []int) (indices [][]int) {
	t = dereferencedType(t)
	fName := make(map[string]struct{})
	var nextLevel [][][]int
	for i, n := 0, t.NumField(); i < n; i++ {
		indices = append(indices, append(prefixIndex, i))
		f := t.Field(i)
		fName[f.Name] = struct{}{}
		if f.Anonymous {
			nextLevel = append(nextLevel, fieldIndex(f.Type, []int{i}))
		}
	}

	for _, nextIndices := range nextLevel {
		for _, index := range nextIndices {
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

func dereferencedType(t reflect.Type) reflect.Type {
	for k := t.Kind(); k == reflect.Ptr || k == reflect.Interface; k = t.Kind() {
		t = t.Elem()
	}

	return t
}

func dereferencedValue(value interface{}) reflect.Value {
	v := reflect.ValueOf(value)

	for k := v.Kind(); k == reflect.Ptr || k == reflect.Interface; k = v.Kind() {
		v = v.Elem()
	}

	return v
}
