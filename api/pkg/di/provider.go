package di

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/dig"
)

type Out = dig.Out

type In = dig.In

type Module func(*Container)

type Container struct {
	container *dig.Container
	invokes   []interface{}
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
	if len(as) == 0 {
		if err := c.container.Provide(provider); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	} else {
		if err := c.container.Provide(provider, dig.As(as...)); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}

	errTyp := reflect.TypeOf(new(error)).Elem()
	cycleForTypes := []reflect.Type{}
	providerType := reflect.TypeOf(provider)
	for i := 0; i < providerType.NumOut(); i++ {
		outType := providerType.Out(i)
		if outType == errTyp {
			continue
		}
		if outType.Kind() == reflect.Ptr {
			continue
		}
		if dig.IsOut(outType) {
			continue
		}
		cycleForTypes = append(cycleForTypes, outType)
	}

	for _, a := range as {
		asTypePtr := reflect.TypeOf(a)
		asType := asTypePtr.Elem()
		if asType.Kind() != reflect.Ptr {
			cycleForTypes = append(cycleForTypes, asType)
		}
	}

	invokes := make([]interface{}, 0, len(cycleForTypes))
	for _, cycleType := range cycleForTypes {
		cycleValPtr := reflect.New(cycleType)

		provideNullFuncType := reflect.FuncOf(nil, []reflect.Type{cycleValPtr.Type()}, false)
		provideNullFun := reflect.MakeFunc(provideNullFuncType, func(args []reflect.Value) []reflect.Value {
			return []reflect.Value{cycleValPtr}
		}).Interface()
		if err := c.container.Provide(provideNullFun); err != nil {
			panic(fmt.Sprintf("failed to provide %s: %+v", fullType(cycleType), err))
		}

		setNullFuncType := reflect.FuncOf([]reflect.Type{cycleType}, nil, false)
		setNullFunc := reflect.MakeFunc(setNullFuncType, func(args []reflect.Value) []reflect.Value {
			cycleValPtr.Elem().Set(args[0])
			return []reflect.Value{}
		}).Interface()

		invokes = append(invokes, setNullFunc)
	}
	c.invokes = append(invokes, c.invokes...)
}

func fullType(in reflect.Type) string {
	return in.PkgPath() + "." + in.String()
}

// Import can be used to combine multiple modules into one.
func (c *Container) Import(module Module) {
	module(c)
}

func isDigError(err error) bool {
	return strings.HasPrefix(fmt.Sprintf("%T", dig.RootCause(err)), "dig.")
}

// Init retrieves an instance of the dest from the container.
func Init(dest interface{}, module Module) error {
	c := &Container{
		container: dig.New(),
	}

	module(c)

	for _, invoke := range c.invokes {
		if err := c.container.Invoke(invoke); err != nil {
			if isDigError(err) {
				return err
			}
			return dig.RootCause(err)
		}
	}

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

	if err := c.container.Invoke(invokeFn.Interface()); err != nil {
		return err
	}

	return nil
}
