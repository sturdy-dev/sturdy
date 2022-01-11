package di_test

import (
	"fmt"
	"testing"

	"mash/pkg/di"

	"github.com/stretchr/testify/assert"
)

func TestModule_simple(t *testing.T) {
	provideString := func() string {
		return "test"
	}
	m := di.NewModule(di.Provides(provideString))
	var target string
	assert.NoError(t, m.Build(&target))
	assert.Equal(t, "test", target)
}

func TestModule_multipleProviders(t *testing.T) {
	provideInt := func() int {
		return 1
	}
	provideString := func(i int) string {
		return fmt.Sprintf("%d", i)
	}
	m := di.NewModule(di.Provides(provideString), di.Provides(provideInt))
	var target string
	assert.NoError(t, m.Build(&target))
	assert.Equal(t, "1", target)
}

func TestModule_multipleModules(t *testing.T) {
	intModule := di.NewModule(di.Provides(func() int {
		return 1
	}))
	stringModule := di.NewModule(di.Provides(func(i int) string {
		return fmt.Sprintf("%d", i)
	}), di.Needs(intModule))

	var target string
	assert.NoError(t, stringModule.Build(&target))
	assert.Equal(t, "1", target)
}

type b struct {
	v string
	A *a
}
type a struct {
	v string
	B *b
}

func TestModule_cycle(t *testing.T) {
	moduleA := di.NewModule(di.ProvidesCycle(func(a *a) b {
		return b{"b", a}
	}))

	moduleB := di.NewModule(di.ProvidesCycle(func(b *b) a {
		return a{"a", b}
	}))

	both := di.NewModule(di.Needs(moduleA), di.Needs(moduleB))

	var targetA a
	if assert.NoError(t, both.Build(&targetA)) {
		assert.NotNil(t, targetA.B)
	}

	var targetB b
	if assert.NoError(t, both.Build(&targetB)) {
		assert.NotNil(t, targetB.A)
	}
}
