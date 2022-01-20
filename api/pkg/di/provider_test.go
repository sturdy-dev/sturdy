package di_test

import (
	"fmt"
	"testing"

	"getsturdy.com/api/pkg/di"

	"github.com/stretchr/testify/assert"
)

func TestModule_simple(t *testing.T) {
	provideString := func() string {
		return "test"
	}
	module := func(c *di.Container) {
		c.Register(provideString)
	}
	var target string
	assert.NoError(t, di.Init(&target, module))
	assert.Equal(t, "test", target)
}

func TestModule_multipleProviders(t *testing.T) {
	provideInt := func() int {
		return 1
	}
	provideString := func(i int) string {
		return fmt.Sprintf("%d", i)
	}
	module := func(c *di.Container) {
		c.Register(provideInt)
		c.Register(provideString)
	}
	var target string
	assert.NoError(t, di.Init(&target, module))
	assert.Equal(t, "1", target)
}

func TestModule_multipleModules(t *testing.T) {
	provideInt := func() int {
		return 1
	}
	intModule := func(c *di.Container) {
		c.Register(provideInt)
	}

	provideString := func(i int) string {
		return fmt.Sprintf("%d", i)
	}
	stringModule := func(c *di.Container) {
		c.Register(provideString)
		c.Import(intModule)
	}

	var target string
	assert.NoError(t, di.Init(&target, stringModule))
	assert.Equal(t, "1", target)
}

type b struct {
	v  string
	a  *a
	ia *IA
}

func (*b) B() {}

type IB interface {
	B()
}

type a struct {
	v  string
	b  *b
	ib *IB
}

func (*a) A() {}

type IA interface {
	A()
}

func TestModule_cycleStructs(t *testing.T) {
	moduleA := func(c *di.Container) {
		c.Register(func(b *b) a {
			return a{v: "a", b: b}
		})
	}
	moduleB := func(c *di.Container) {
		c.Register(func(a *a) b {
			return b{v: "b", a: a}
		})
	}

	both := func(c *di.Container) {
		c.Import(moduleA)
		c.Import(moduleB)
	}

	var targetA a
	if assert.NoError(t, di.Init(&targetA, both)) {
		assert.NotNil(t, targetA.b)
	}

	var targetB b
	if assert.NoError(t, di.Init(&targetB, both)) {
		assert.NotNil(t, targetB.a)
	}
}

func TestModule_interface(t *testing.T) {
	provideA := func() *a {
		return &a{v: "a"}
	}
	moduleA := func(c *di.Container) {
		c.Register(provideA, new(IA))
	}
	var target IA
	assert.NoError(t, di.Init(&target, moduleA))
	if assert.NotNil(t, target) {
		assert.Equal(t, "a", target.(*a).v)
	}
}

func TestModule_cycleInterfaces(t *testing.T) {
	moduleA := func(c *di.Container) {
		c.Register(func(ib *IB) IA {
			return &a{v: "a", ib: ib}
		})
	}
	moduleB := func(c *di.Container) {
		c.Register(func(ia *IA) *b {
			return &b{v: "b", ia: ia}
		}, new(IB))
	}

	type result struct {
		ia *IA
		ib *IB
	}
	both := func(c *di.Container) {
		c.Import(moduleA)
		c.Import(moduleB)
		c.Register(func(ib *IB, ia *IA) *result {
			return &result{ia: ia, ib: ib}
		})
	}

	other := func(c *di.Container) {
		c.Import(both)
	}

	var target *result
	assert.NoError(t, di.Init(&target, other))

	var targetA IA
	var targetB IB
	assert.NoError(t, di.Init(&targetA, both))
	assert.NoError(t, di.Init(&targetB, both))

	vb, ok := (*targetA.(*a).ib).(*b)
	assert.True(t, ok)
	assert.NotNil(t, vb)

	va, ok := (*targetB.(*b).ia).(*a)
	assert.True(t, ok)
	assert.NotNil(t, va)
}
