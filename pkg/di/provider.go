package di

import (
	"fmt"
	"reflect"

	"go.uber.org/dig"
)

type Module interface {
	Build(interface{}) error
}

type module struct {
	hooks []Hook
}

func NewModule(hooks ...Hook) *module {
	return &module{
		hooks: hooks,
	}
}

func (m *module) Invoke(fn interface{}) {
	m.hooks = append(m.hooks, func(c *dig.Container) (invoke, error) {
		return func(c *dig.Container) error {
			return c.Invoke(fn)
		}, nil
	})
}

func (m *module) Build(dest interface{}) error {
	container := dig.New()
	var ii []invoke
	for _, hook := range m.hooks {
		i, err := hook(container)
		if err != nil {
			return err
		}
		if i != nil {
			ii = append(ii, i)
		}
	}
	for _, i := range ii {
		if err := i(container); err != nil {
			return err
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

	if err := container.Invoke(invokeFn.Interface()); err != nil {
		return err
	}
	return nil
}

type invoke func(c *dig.Container) error

type Hook func(*dig.Container) (invoke, error)

func Needs(m Module) Hook {
	return func(c *dig.Container) (invoke, error) {
		ii := []invoke{}
		for _, hook := range m.(*module).hooks {
			i, err := hook(c)
			if err != nil {
				return nil, err
			}
			if i != nil {
				ii = append(ii, i)
			}
		}
		return func(c *dig.Container) error {
			for _, i := range ii {
				if err := i(c); err != nil {
					return err
				}
			}
			return nil
		}, nil
	}
}

func Invoke(fn interface{}) Hook {
	return func(c *dig.Container) (invoke, error) {
		return func(c *dig.Container) error {
			return c.Invoke(fn)
		}, nil
	}
}

func ProvidesCycle(provider interface{}, as ...interface{}) Hook {
	outType := reflect.TypeOf(provider).Out(0)
	outTypePtr := reflect.PtrTo(outType)
	provideFn := providerFor(outTypePtr)
	invokeFn := invokerFor(outTypePtr)

	oo := []dig.ProvideOption{}
	for _, as := range as {
		oo = append(oo, dig.As(as))
	}
	return func(c *dig.Container) (invoke, error) {
		if err := c.Provide(provider, oo...); err != nil {
			return nil, err
		}
		if err := c.Provide(provideFn); err != nil {
			return nil, err
		}
		return func(c *dig.Container) error {
			return c.Invoke(invokeFn)
		}, nil
	}
}

func Provides(provider interface{}, as ...interface{}) Hook {
	oo := []dig.ProvideOption{}
	for _, as := range as {
		oo = append(oo, dig.As(as))
	}
	return func(c *dig.Container) (invoke, error) {
		return nil, c.Provide(provider, oo...)
	}
}
