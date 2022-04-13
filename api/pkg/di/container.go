package di

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/dig"
)

type Out = dig.Out

type In = dig.In

type Module func(*Container)

type provider struct {
	ID   string
	Type reflect.Type
	Func any
	As   []any
}

func newProvider(p any, as ...any) *provider {
	id := getFuncName(p)
	return &provider{
		ID:   id,
		Type: reflect.TypeOf(p),
		Func: p,
		As:   as,
	}
}

func (p *provider) provides() map[reflect.Type]bool {
	result := map[reflect.Type]bool{}
	if len(p.As) == 0 {
		for i := 0; i < p.Type.NumOut(); i++ {
			outType := p.Type.Out(i)
			if outType == errTyp {
				continue
			}
			addProvides(outType, result)
		}
	} else {
		for _, a := range p.As {
			asTypePtr := reflect.TypeOf(a)
			asType := asTypePtr.Elem()
			if asType.Kind() != reflect.Ptr {
				addProvides(asType, result)
			}
		}
	}
	return result
}

type Container struct {
	name string

	imports            []*Container
	importsWithForce   []*Container
	providers          map[string]*provider
	providersWithForce map[string]*provider
	decorators         map[string]any

	// caches
	_provides    map[reflect.Type]bool
	providesOnce sync.Once

	_requires    map[reflect.Type]bool
	requiresOnce sync.Once

	_allProviders    map[string]*provider
	allProvidersOnce sync.Once

	_allDecorators    map[string]any
	allDecoratorsOnce sync.Once

	isValid     error
	isValidOnce sync.Once
}

var errTyp = reflect.TypeOf(new(error)).Elem()

func addRequires(in reflect.Type, to map[reflect.Type]bool) {
	if in.Name() == "digSentinel" {
		return
	}

	if !dig.IsIn(in) {
		to[in] = true
		return
	}

	for i := 0; i < in.NumField(); i++ {
		field := in.Field(i)
		addRequires(field.Type, to)
	}
}

func (c *Container) requires() map[reflect.Type]bool {
	c.requiresOnce.Do(func() {
		result := map[reflect.Type]bool{}
		for _, provider := range c.allProviders() {
			providerType := reflect.TypeOf(provider.Func)
			for i := 0; i < providerType.NumIn(); i++ {
				inType := providerType.In(i)
				addRequires(inType, result)
			}
		}
		for _, decorator := range c.allDecorators() {
			decoratorType := reflect.TypeOf(decorator)
			for i := 0; i < decoratorType.NumIn(); i++ {
				inType := decoratorType.In(i)
				addRequires(inType, result)
			}
		}
		c._requires = result
	})
	return c._requires
}

func addProvides(in reflect.Type, to map[reflect.Type]bool) {
	if in.Name() == "digSentinel" {
		return
	}

	if !dig.IsOut(in) {
		to[in] = true
		return
	}

	for i := 0; i < in.NumField(); i++ {
		field := in.Field(i)
		addProvides(field.Type, to)
	}
}

func (c *Container) provides() map[reflect.Type]bool {
	c.providesOnce.Do(func() {
		result := map[reflect.Type]bool{}
		for _, provider := range c.allProviders() {
			result = merge(result, provider.provides())
		}
		c._provides = result
	})
	return c._provides
}

func (c *Container) IsValid() error {
	c.isValidOnce.Do(func() {
		for _, i := range c.imports {
			if err := i.IsValid(); err != nil {
				c.isValid = fmt.Errorf("%s: %w", c.name, err)
				return
			}
		}
		errors := []error{}
		provides := c.provides()

		for r := range c.requires() {
			if !provides[r] {
				errors = append(errors, fmt.Errorf("missing provider for %s", getTypeName(r)))
			}
		}
		c.isValid = joinErrors(errors)
	})
	return c.isValid
}

func merge(a, b map[reflect.Type]bool) map[reflect.Type]bool {
	r := map[reflect.Type]bool{}
	for k, v := range a {
		r[k] = v
	}
	for k, v := range b {
		r[k] = v
	}
	return r
}

func joinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	if len(errors) == 1 {
		return errors[0]
	}
	errStrings := make([]string, len(errors))
	for i, err := range errors {
		errStrings[i] = err.Error()
	}
	return fmt.Errorf(strings.Join(append([]string{""}, errStrings...), "\n	- "))
}

// Import imports a module into the container.
//
// All content of the imported module will be available
// in the container.
func (c *Container) Import(m Module) {
	c.imports = append(c.imports, Init(m))
}

// ImportWithForce is like import, but will overwrite any existing providers.
func (c *Container) ImportWithForce(m Module) {
	c.importsWithForce = append(c.importsWithForce, Init(m))
}

// Register registers a single provider, for example:
//
//   // this adds a provider for *MyType into the container.
//   c.Register(func() *MyType { return &MyType{} } )
//
//   // this adds a provider for *MyOtherType into the container that depends on *MyType.
//   // provider for *MyType must be registred in the container for successful Init.
//   c.Register(func(*MyType) *MyOtherType { return &MyOtherType{} } )
//
//   // you can use as... parameters to register provider for an interface.
//   // this adds a provider for MyInterface into the container.
//   // IMPORTANT: this will not register a provider for *MyType.
//   c.Register(func() *MyType { return &MyType{} }, new(MyInterface) )
func (c *Container) Register(p any, as ...any) {
	provider := newProvider(p, as...)
	c.providers[provider.ID] = provider
}

// RegisterWithForce is like Register, but will overwrite any existing providers.
func (c *Container) RegisterWithForce(p any, as ...any) {
	provider := newProvider(p, as...)
	c.providersWithForce[provider.ID] = provider
}

func (c *Container) Decorate(decorator any) {
	id := getFuncName(decorator)
	if _, ok := c.decorators[id]; ok {
		return
	}
	c.decorators[id] = decorator
}

func (c *Container) allProviders() map[string]*provider {
	c.allProvidersOnce.Do(func() {
		result := map[string]*provider{}
		// recursively import all subcontainers
		for _, i := range c.imports {
			for id, provider := range i.allProviders() {
				result[id] = provider
			}
		}
		// register all providers
		for id, provider := range c.providers {
			result[id] = provider
		}

		// make a map of provided types
		provided := map[reflect.Type]*provider{}
		for _, provider := range result {
			for t := range provider.provides() {
				provided[t] = provider
			}
		}

		// override with force if needed
		for _, i := range c.importsWithForce {
			for id, provider := range i.allProviders() {
				for t := range provider.provides() {
					if provider, ok := provided[t]; ok {
						delete(result, provider.ID)
					}
					provided[t] = provider
				}
				result[id] = provider
			}
		}

		for id, provider := range c.providersWithForce {
			for t := range provider.provides() {
				if provider, ok := provided[t]; ok {
					delete(result, provider.ID)
				}
				provided[t] = provider
			}
			result[id] = provider
		}

		c._allProviders = result
	})
	return c._allProviders
}

func (c *Container) allDecorators() map[string]any {
	c.allDecoratorsOnce.Do(func() {
		result := map[string]any{}
		for _, i := range c.imports {
			for k, v := range i.allDecorators() {
				result[k] = v
			}
		}
		for k, v := range c.decorators {
			result[k] = v
		}
		c._allDecorators = result
	})
	return c._allDecorators
}

func (c *Container) register(container *dig.Container) error {
	for _, provider := range c.allProviders() {
		if len(provider.As) == 0 {
			if err := container.Provide(provider.Func); err != nil {
				return err
			}
		} else {
			if err := container.Provide(provider.Func, dig.As(provider.As...)); err != nil {
				return err
			}
		}
	}
	for _, decorator := range c.allDecorators() {
		if err := container.Decorate(decorator); err != nil {
			return err
		}
	}
	return nil
}

// To builds the container and fetches dest from it, dest must be a pointer.
func (c *Container) To(dest ...any) error {
	if err := c.IsValid(); err != nil {
		return fmt.Errorf("%s: %w", c.name, err)
	}

	container := dig.New()
	if err := c.register(container); err != nil {
		return err
	}

	destPointers := make([]reflect.Value, 0, len(dest))
	destTypes := make([]reflect.Type, 0, len(dest))
	for _, d := range dest {
		destPtrValue := reflect.ValueOf(d)
		destPointers = append(destPointers, destPtrValue)

		destTypePtr := destPtrValue.Type()
		if destTypePtr.Kind() != reflect.Ptr {
			return fmt.Errorf("%T: must be a pointer", d)
		}
		destTypes = append(destTypes, destTypePtr.Elem())
	}

	invokeFnType := reflect.FuncOf(destTypes, nil, false)
	invokeFn := reflect.MakeFunc(invokeFnType, func(args []reflect.Value) []reflect.Value {
		for i, dp := range destPointers {
			dp.Elem().Set(args[i])
		}
		return []reflect.Value{}
	})

	if err := container.Invoke(invokeFn.Interface()); err != nil {
		return err
	}

	return nil

}

var globalRegistry = map[string]*Container{}

// Init prepares the container contents.
// Only makes sense to use it in conjunction with To:
//
// Init(myModule).To(&dest)
func Init(module Module) *Container {
	name := getFuncName(module)
	if _, ok := globalRegistry[name]; ok {
		return globalRegistry[name]
	}
	c := &Container{
		name:               getFuncName(module),
		providers:          map[string]*provider{},
		providersWithForce: map[string]*provider{},
		decorators:         map[string]any{},
	}
	module(c)
	globalRegistry[name] = c
	return c
}

func getFuncName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func getTypeName(in reflect.Type) string {
	if in.Kind() == reflect.Ptr {
		in = in.Elem()
		if in.PkgPath() == "" {
			return in.Name()
		}
		return fmt.Sprintf("*%s.%s", in.PkgPath(), in.Name())
	} else {
		if in.PkgPath() == "" {
			return in.Name()
		}
		return fmt.Sprintf("%s.%s", in.PkgPath(), in.Name())
	}
}
