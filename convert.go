package struct_converter

import (
	"reflect"
)

type field struct {
	sTyp reflect.Type
	dTyp reflect.Type
	sDereferTyp reflect.Type
	dDereferTyp reflect.Type
	sReferDeep	int
	dReferDeep	int
	//dstOffset uintptr  //todo 使用代码写死各种指针类型转换时使用
	//srcOffset uintptr
}

func newField(st, dt reflect.Type) (f *field, ok bool) {
	var sReferDeep, dReferDeep int
	sDereferTyp, dDereferTyp := st, dt
	for k := sDereferTyp.Kind(); k==reflect.Ptr; k=sDereferTyp.Kind() {
		sDereferTyp = sDereferTyp.Elem()
		sReferDeep +=1
	}

	for k := dDereferTyp.Kind(); k==reflect.Ptr; k=dDereferTyp.Kind() {
		dDereferTyp = dDereferTyp.Elem()
		dReferDeep +=1
	}

	f = &field{
		sTyp:st,
		dTyp:dt,
		sDereferTyp: sDereferTyp,
		dDereferTyp: dDereferTyp,
		sReferDeep:sReferDeep,
		dReferDeep:dReferDeep,
	}
	if sDereferTyp.ConvertibleTo(dDereferTyp) {
		ok = true
	}
	return
}

func (f *field) convert(sv, dv reflect.Value) {
	if sv.Type() != f.sTyp {
		panic("source field wrong type!")
	}

	if dv.Type() != f.dTyp {
		panic("target field wrong type!")
	}

	if sv.Kind() == reflect.Ptr && sv.IsNil() {
		sv = reflect.New(f.sDereferTyp).Elem()
	} else {
		for d:=0; d<f.sReferDeep; d++ {
			sv = sv.Elem()
		}
	}
	v := sv.Convert(f.dDereferTyp)

	for t,tmp,d:=f.dDereferTyp, reflect.New(f.dDereferTyp).Elem(), 0; d<f.dReferDeep; d++ {
		tmp.Set(v)
		v = tmp.Addr()
		t = reflect.PtrTo(t)
		tmp = reflect.New(t).Elem()
	}

	dv.Set(v)
}

type Converter struct {
	dstTyp reflect.Type
	srcTyp reflect.Type
	fields map[string]*field
}

func New(src interface{}, dst interface{}) *Converter {
	srcTyp := dereferencedType(reflect.TypeOf(src))
	dstTyp := dereferencedType(reflect.TypeOf(dst))

	if srcTyp.Kind() != reflect.Struct {
		panic("source is not a struct!")
	}

	if dstTyp.Kind() != reflect.Struct {
		panic("target is not a struct!")
	}

	fields := make(map[string]*field)
	for i, n :=0, dstTyp.NumField();i<n;i++{
		df := dstTyp.Field(i)
		name := df.Name
		if sf, ok := srcTyp.FieldByName(name); ok {
			if field, ok := newField(sf.Type, df.Type); ok {
				fields[name] = field
			}
		}
	}

	return &Converter{
		dstTyp,
		srcTyp,
		fields,
	}
}

func (c *Converter) Convert(src interface{}, dst interface{}) {
	dv := dereferencedValue(dst)
	if !dv.CanSet() {
		panic("target should be a pointer!")
	}
	sv := dereferencedValue(src)

	if sv.Type() != c.srcTyp {
		panic("source struct wrong type!")
	}

	if dv.Type() != c.dstTyp {
		panic("target struct wrong type!")
	}

	for name, field := range c.fields {
		df := dv.FieldByName(name)
		sf := sv.FieldByName(name)
		field.convert(sf, df)
	}
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