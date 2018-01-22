package struct_converter

import (
	"reflect"
)

type cvtFlag int

const (
	vToV cvtFlag = iota
	pToV
	vToP
	pToP
)

var cvtOp = map[cvtFlag]func(dfv, sfv reflect.Value){
	vToV:cvtVV,
	vToP:cvtVP,
	pToV:cvtPV,
	pToP:cvtPP,
}

func cvtVV(sfv, dfv reflect.Value) {
	v := sfv.Convert(dfv.Type())
	dfv.Set(v)
}

func cvtVP(sfv, dfv reflect.Value) {
	deep := 0
	t := dfv.Type()
	for k := t.Kind(); k==reflect.Ptr; k= t.Kind(){
		t=t.Elem()
		deep +=1
	}
	var v reflect.Value
	for d:=0; d<deep; d++ {
		v = reflect.New(t).Elem()
		v.Set(sfv.Convert(t))
		v = v.Addr()

		t = reflect.PtrTo(t)
	}
	dfv.Set(v)
}

func cvtPV(sfv, dfv reflect.Value) {
	var v reflect.Value
	if sfv.IsNil() {
		v = reflect.New(dfv.Type()).Elem()
	} else {
		v = sfv.Elem().Convert(dfv.Type())
	}
	dfv.Set(v)
}

func cvtPP(sfv, dfv reflect.Value) {
	v := reflect.New(dfv.Type().Elem()).Elem()
	if sfv.IsNil() {
		v.Set(reflect.New(dfv.Type().Elem()).Elem())
	} else {
		v.Set(sfv.Elem().Convert(dfv.Type().Elem()))
	}
	dfv.Set(v.Addr())
}


func canCvtVV(st, dt reflect.Type) bool {
	if dt.Kind() != reflect.Ptr && dt.Kind() != reflect.Ptr && st.ConvertibleTo(dt) {
		return true
	}
	return false
}

func canCvtVP(st, dt reflect.Type) bool {
	if dt.Kind() == reflect.Ptr && st.Kind() != reflect.Ptr && st.ConvertibleTo(dt.Elem()) {
		return true
	}
	return false
}

func canCvtPV(st, dt reflect.Type) bool {
	if st.Kind() == reflect.Ptr && dt.Kind() != reflect.Ptr && st.Elem().ConvertibleTo(dt) {
		return true
	}
	return false
}

func canCvtPP(st, dt reflect.Type) bool {
	if st.Kind() == reflect.Ptr && dt.Kind() == reflect.Ptr && st.Elem().ConvertibleTo(dt.Elem()) {
		return true
	}
	return false
}

func canCvt(st, dt reflect.Type) (int, int, bool) {
	sDeep, dDeep := 0,0
	for k := st.Kind(); k==reflect.Ptr; k=st.Kind() {
		st = st.Elem()
		sDeep +=1
	}
	for k := dt.Kind(); k==reflect.Ptr; k=dt.Kind() {
		dt = dt.Elem()
		dDeep +=1
	}
	if st.ConvertibleTo(dt) {
		return sDeep, dDeep, true
	} else {
		return 0,0,false
	}
}

type field struct {
	name string
	dTyp reflect.Type
	//dstOffset uintptr  //todo 使用代码写死各种指针类型转换时使用
	//srcOffset uintptr
	cvtFlag cvtFlag
}

func NewField(name string, st, dt reflect.Type) (*field, bool) {
	var flag cvtFlag
	if canCvtVV(st, dt) {
		flag = vToV
	}else if canCvtVP(st, dt) {
		flag = vToP
	} else if canCvtPV(st, dt) {
		flag = pToV
	} else if canCvtPP(st, dt) {
		flag = pToP
	} else {
		return nil, false
	}

	f := &field{
		name:name,
		dTyp:dt,
		cvtFlag:flag,
	}
	return f, true
}

func (f *field) convert(sfv, dfv reflect.Value)  {
	cvtOp[f.cvtFlag](sfv, dfv)
}

type Converter struct {
	dstTyp reflect.Type
	srcTyp reflect.Type
	fields []*field
}

func New(src interface{}, dst interface{}) *Converter {
	dstTyp := deferencedType(reflect.TypeOf(dst))
	srcTyp := deferencedType(reflect.TypeOf(src))

	n := dstTyp.NumField()
	fields := make([]*field, 0, n)
	for i:=0;i<n;i++{
		df := dstTyp.Field(i)
		name := df.Name
		if sf, ok := srcTyp.FieldByName(name); ok {
			if field, ok :=  NewField(name, sf.Type, df.Type); ok {
				fields = append(fields, field)
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
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src).Elem()
	for _, field := range c.fields {
		df := dv.FieldByName(field.name)
		sf := sv.FieldByName(field.name)
		field.convert(sf, df)
	}
}

func deferencedType(t reflect.Type) reflect.Type {
	for k := t.Kind(); k == reflect.Ptr || k == reflect.Interface; k = t.Kind() {
		t = t.Elem()
	}

	return t
}
