package di

import (
	"fmt"
	"reflect"

	"go.uber.org/dig"
)

type Module func(*Container)

type Container struct {
	container *dig.Container
}

// Register teaches container how to create an instance of the given type, for instance:
//
//     container.Register(
//         // tell the container how to create an instance of *MyType
//         func() *MyType { return &MyType{} },
//     )
//
//     container.Register(
//         // tell the container how to create an instance of *MyOtherType from *MyType
//         func(*MyType) *MyOtherType { ... },
//         // second argument registers *MyOtherType as MyOtherTypeInterface too.
//         new(MyOtherTypeInterface),
//     )
//
// Register also registers resolvers for pointers to the type, so circular dependencies are
// possible:
//
//     container.Register(func(*MyOtherType) MyType { ... })
//     container.Register(func(*MyType) MyOtherTypeType { ... })
//
func (c *Container) Register(provider interface{}, as ...interface{}) {
	c.container.Provide(provider) // register provider for the type itself

	oo := []dig.ProvideOption{}
	asTypes := []reflect.Type{}
	for _, as := range as {
		oo = append(oo, dig.As(as))
		asTypePtr := reflect.TypeOf(as)
		if asTypePtr.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("as must be a pointer, got %T", as))
		}
		asTypes = append(asTypes, reflect.TypeOf(as))
	}

	c.container.Provide(provider, oo...) // register provider for the _as_ implementations

	// register cycle resolvers for the type
	outType := reflect.TypeOf(provider).Out(0)
	outTypePtr := reflect.PtrTo(outType)
	provideFn := providerFor(outTypePtr)
	invokeFn := invokerFor(outTypePtr)

	c.container.Provide(provideFn)
	c.container.Invoke(invokeFn)

	for _, asType := range asTypes {
		cycleInvoker := invokerFor(asType)
		cycleProvider := providerFor(asType)
		c.container.Invoke(cycleInvoker)
		c.container.Provide(cycleProvider)
	}
}

// Import can be used to combine multiple modules into one.
func (c *Container) Import(module Module) {
	module(c)
}

// Init retrieves an instance of the dest from the container.
func Init(dest interface{}, module Module) error {
	c := &Container{
		container: dig.New(),
	}
	module(c)

	destPtrValue := reflect.ValueOf(dest)
	destTypePtr := destPtrValue.Type()
	if destTypePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	destType := destTypePtr.Elem()
	invokeFnType := reflect.FuncOf([]reflect.Type{destType}, nil, false)
	invokeFn := reflect.MakeFunc(invokeFnType, func(args []reflect.Value) []reflect.Value {
		destPtrValue.Elem().Set(args[0])
		return []reflect.Value{}
	})

	return c.container.Invoke(invokeFn.Interface())
}
