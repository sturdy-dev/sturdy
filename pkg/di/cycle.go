package di

import (
	"reflect"
)

func providerFor(typPtr reflect.Type) interface{} {
	funcTyp := reflect.FuncOf(nil, []reflect.Type{typPtr}, false)
	typ := typPtr.Elem()
	// create a new instance of the type
	v := reflect.New(typ).Elem()
	// create a new instance of a pointer to the type
	vPtr := reflect.New(typPtr).Elem()
	// set pointer to the value
	vPtr.Set(v.Addr())
	return reflect.MakeFunc(funcTyp, func(_ []reflect.Value) []reflect.Value {
		return []reflect.Value{vPtr}
	}).Interface()
}

func invokerFor(typPtr reflect.Type) interface{} {
	typ := typPtr.Elem()
	funcTyp := reflect.FuncOf([]reflect.Type{typPtr, typ}, nil, false)
	return reflect.MakeFunc(funcTyp, func(args []reflect.Value) []reflect.Value {
		// set pointer to the value
		args[0].Elem().Set(args[1])
		return nil
	}).Interface()
}
