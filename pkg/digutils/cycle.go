package digutils

import (
	"reflect"
)

// given an interface , returns dig.Provider and dig.Invoker to resolve a cycle dependency
//
// for example, calling:
//    provider, invoker := ResolveCycleFor(new(MyInterface))
// will return an equivalent of:
//    provider := func() *MyInterface {
//        return new(MyInterface)
//    }
//    invoker := func(p *MyInterface, t MyInterface) {
//        *p = t
//    }
//
// that allows to resolve a cycle dependency like so:
//     c := dig.New()
//
//     // this is a provider for MyInterface
//     _ = c.Provide(func() MyInterface { ... })
//
//     // this is a provider that requires MyInterface. Note that we have to use
//     // a pointer to MyInterface here
//     _ = c.Provide(func(i *MyInterface) *Result {
//         return &Result{MyInterface: i} // at this point, *MyInterface is not yet resolved
//     })
//
//     provider, invoker := ResolveCycleFor(new(MyInterface))
//     _ = c.Provide(provider)
//     _ = c.Invoke(invoker)
//
//     _ = c.Invoke(func(r Result) {
//         fmt.Println(r.MyInterface) // now, MyInterface is resolved
//     })
func ResolveCycleFor(i interface{}) (provider interface{}, invoker interface{}) {
	typ := reflect.TypeOf(i)
	return providerFor(typ), invokerFor(typ)
}

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
