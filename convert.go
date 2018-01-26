package coven

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"
	"unsafe"
)

type convertType struct {
	srcTyp reflect.Type
	dstTyp reflect.Type
}

var (
	createdConvertersMu sync.Mutex
	createdConverters   = make(map[convertType]*converter)
)

// converter stores convertible fields of srcTyp and dstTyp.
// Field type of nested pointer is supported.
// srcTyp and dstTyp are both dereferenced reflect.Type

// All methods in converter are thread-safe.
// We can define a global variable to hold a converter and use it in any goroutine.
type converter struct {
	convertType
	fieldConverters []*fieldConverter
}

// NewConverter analyzes type information of src and dst
// and returns a *converter with convertible fields of the same name.
// Field type of nested pointer is supported.
// It panics if src or dst is not a struct.
func NewConverter(dst interface{}, src interface{}) (c *converter) {
	return newConverter(dst, src, true)
}

func newConverter(dst interface{}, src interface{}, lock bool) (c *converter) {
	srcTyp := dereferencedType(reflect.TypeOf(src))
	dstTyp := dereferencedType(reflect.TypeOf(dst))

	cTyp := convertType{srcTyp, dstTyp}

	if lock {
		createdConvertersMu.Lock()
		defer createdConvertersMu.Unlock()
	}
	if c, ok := createdConverters[cTyp]; ok {
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
			if fCvt, ok := newFieldConverter(df, sf); ok {
				fCvts = append(fCvts, fCvt)
			}
		}
	}

	if len(fCvts) > 0 {
		c = &converter{
			cTyp,
			fCvts,
		}
		createdConverters[cTyp] = c
	}

	return
}

// Convert creates field values converted from src and set them in dst on fields stored in converter.
// Field type of nested pointer is supported.
// dst should be a pointer to a struct, otherwise Convert panics.
// dereferenced src and dst type should match their counterparts in converter.
func (c *converter) Convert(dst interface{}, src interface{}) {
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
		if fCvt.cvtOp != nil {
			dPtr := unsafe.Pointer(dv.UnsafeAddr() + fCvt.dOffset)
			sPtr := unsafe.Pointer(sv.UnsafeAddr() + fCvt.sOffset)
			fCvt.convertByPtr(dPtr, sPtr)
		} else {
			sf := sv.FieldByIndex(fCvt.sIndex)
			df := dv.FieldByIndex(fCvt.dIndex)
			fCvt.convert(df, sf)
		}
	}
}

// fieldConverter stores information of convertible field from one type to another.
// sTyp and dTyp are original types of src and dst.
// sDereferTyp and dDereferTyp are dereferenced types of src and dst.
// sReferDeep and dReferDeep are levels of nested pointer of src and dst.
// If src and dst are different struct, make them a converter.
type fieldConverter struct {
	sOffset     uintptr
	dOffset     uintptr
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
	cvtOp       cvtOp
	structCvt   *converter
}

// newFieldConverter analyzes information of src and dst field
// and returns a *fieldConverter, if they are convertible, ok is true.
// Field type of nested pointer is supported.
func newFieldConverter(df, sf reflect.StructField) (f *fieldConverter, ok bool) {
	f = &fieldConverter{
		sOffset:     sf.Offset,
		dOffset:     df.Offset,
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
		cvtOp:       nil,
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
		sk, dk := f.sDereferTyp.Kind(), f.dDereferTyp.Kind()
		f.cvtOp = cvtOps[convertKind{sk, dk}]
		ok = true
	} else if f.sDereferTyp.Kind() == reflect.Struct && f.dDereferTyp.Kind() == reflect.Struct {
		f.structCvt = newConverter(reflect.New(f.dDereferTyp).Interface(), reflect.New(f.sDereferTyp).Interface(), false)
		if f.structCvt != nil {
			ok = true
		}
	}

	return
}

// convert creates a value converted from src field and set it in dst field.
// The new value is first created as type of dDereferTyp,
// and then pointer nested for dReferDeep times to become a dTyp value.
func (f *fieldConverter) convert(dv, sv reflect.Value) {
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
		f.structCvt.Convert(v.Interface(), sv.Addr().Interface())
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

func (f *fieldConverter) convertByPtr(dPtr, sPtr unsafe.Pointer) {
	if *((**int)(sPtr)) == nil {
		sPtr = newValue(f.sDereferTyp.Kind())
	} else {
		for d := 0; d < f.sReferDeep; d++ {
			sPtr = unsafe.Pointer(*((**int)(sPtr)))
		}
	}

	if f.dReferDeep > 0 {
		v := newValue(f.dDereferTyp.Kind())
		f.cvtOp(sPtr, v)
		for d := 0; d < f.dReferDeep; d++ {
			tmp := v
			v = unsafe.Pointer(&tmp)
		}
		*((*int)(dPtr)) = *(*int)(v)
	} else {
		sPtr := unsafe.Pointer(sPtr)
		dPtr := unsafe.Pointer(dPtr)
		f.cvtOp(sPtr, dPtr)
	}
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
