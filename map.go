package coven

import (
	"reflect"
	"unsafe"
)

type mapConverter struct {
	*convertType
	dKeyTyp            reflect.Type
	dValTyp            reflect.Type
	keyConverter       *elemConverter
	valConverter       *elemConverter
	sEmptyMapInterface *emptyInterface
	dEmptyMapInterface *emptyInterface
}

func newMapConverter(convertType *convertType) converter {
	sKeyTyp := convertType.srcTyp.Key()
	dKeyTyp := convertType.dstTyp.Key()
	sValTyp := convertType.srcTyp.Elem()
	dValTyp := convertType.dstTyp.Elem()
	if keyConverter, ok := newElemConverter(dKeyTyp, sKeyTyp); ok {
		if valueConverter, ok := newElemConverter(dValTyp, sValTyp); ok {
			sEmpty := reflect.New(convertType.srcTyp).Interface()
			dEmpty := reflect.New(convertType.dstTyp).Interface()
			m := &mapConverter{
				convertType,
				dKeyTyp,
				dValTyp,
				keyConverter,
				valueConverter,
				(*emptyInterface)(unsafe.Pointer(&sEmpty)),
				(*emptyInterface)(unsafe.Pointer(&dEmpty)),
			}
			return m
		}
	}
	return nil
}

//dv and sv must be dereferened value
func (m *mapConverter) convert(dPtr, sPtr unsafe.Pointer) {
	sv := ptrToMapValue(m.sEmptyMapInterface, sPtr)
	dv := ptrToMapValue(m.dEmptyMapInterface, dPtr)
	if dv.IsNil() {
		dv.Set(reflect.MakeMap(m.dstTyp))
	}

	for _, sKey := range sv.MapKeys() {
		sVal := sv.MapIndex(sKey).Interface()
		sValPtr := (*emptyInterface)(unsafe.Pointer(&sVal)).word
		sKeyObj := sKey.Interface()
		sKeyPtr := (*emptyInterface)(unsafe.Pointer(&sKeyObj)).word
		dKey := reflect.New(m.dKeyTyp).Elem()
		dVal := reflect.New(m.dValTyp).Elem()
		m.keyConverter.convert(unsafe.Pointer(dKey.UnsafeAddr()), sKeyPtr)
		m.valConverter.convert(unsafe.Pointer(dVal.UnsafeAddr()), sValPtr)
		dv.SetMapIndex(dKey, dVal)
	}

}

func ptrToMapValue(emptyMapInterface *emptyInterface, ptr unsafe.Pointer) reflect.Value {
	emptyMapInterface.word = ptr
	realInterface := *(*interface{})(unsafe.Pointer(emptyMapInterface))
	return reflect.ValueOf(realInterface).Elem()
}
