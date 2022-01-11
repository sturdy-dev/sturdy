package digutils_test

import (
	"mash/pkg/digutils"

	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestResolveCycleFor(t *testing.T) {
	type MyInterface interface{}

	resolver, invoker := digutils.ResolveCycleFor(new(MyInterface))

	provide, ok := resolver.(func() *MyInterface)
	assert.True(t, ok, "resolver has a wrong type")
	if assert.NotNil(t, provide()) {
		assert.Nil(t, *provide())
	}

	invoke, ok := invoker.(func(*MyInterface, MyInterface))
	assert.True(t, ok, "invoker has a wrong type")

	c := dig.New()
	assert.NoError(t, c.Provide(func() MyInterface { return "implementation" }))
	assert.NoError(t, c.Provide(provide))
	assert.NoError(t, c.Invoke(invoke))

	c.Invoke(func(i *MyInterface) {
		if assert.NotNil(t, i) {
			assert.Equal(t, "implementation", *i)
		}
	})
}
