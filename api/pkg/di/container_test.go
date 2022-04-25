package di

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer_provides(t *testing.T) {
	type out struct {
		Out
		i int64
		f float64
	}

	c1 := func(c *Container) {
		c.Register(func() string { return "m1" })
	}
	c2 := func(c *Container) {
		c.Register(func() (int, float32) { return 1, .1 })
		c.Register(func() out { return out{i: 1, f: .1} })
	}
	c3 := func(c *Container) {
		c.Import(c1)
		c.Import(c2)
	}

	assert.Len(t, Init(c1).provides(), 1)
	assert.Len(t, Init(c2).provides(), 4)
	assert.Len(t, Init(c3).provides(), 5)
}

func TestContainer_requires(t *testing.T) {
	type out struct {
		Out
		i int64
		f float64
	}

	type in struct {
		In
		int64
		float64
	}

	c1 := func(c *Container) {
		c.Register(func(byte) string { return "m1" })
	}
	c2 := func(c *Container) {
		c.Register(func(float32) (int, float32) { return 1, .1 })
		c.Register(func(in) out { return out{i: 1, f: .1} })
	}
	c3 := func(c *Container) {
		c.Import(c1)
		c.Import(c2)
	}

	assert.Len(t, Init(c1).requires(), 1)
	assert.Len(t, Init(c2).requires(), 3)
	assert.Len(t, Init(c3).requires(), 4)
}

func TestContainer_IsValid_missing_provider(t *testing.T) {
	m1 := func(i int) string { return fmt.Sprint(i) }
	c := Init(func(c *Container) {
		c.Register(m1)
	})
	assert.Equal(t, "\n\tmissing provider for int", c.IsValid().Error())
}

func TestContainer_IsValid_imports(t *testing.T) {
	c1 := func(c *Container) {
		c.Register(func() int { return 1 })
	}
	c2 := func(c *Container) {
		c.Import(c1)
		c.Register(func(i int) string { return fmt.Sprint(i) })
	}
	assert.NoError(t, Init(c2).IsValid())
}

func TestContainer_To_simple(t *testing.T) {
	c1 := func(c *Container) {
		c.Register(func() int { return 1 })
		c.Register(func() string { return "string" })
	}
	var i int
	var s string
	if assert.NoError(t, Init(c1).To(&i, &s)) {
		assert.Equal(t, 1, i)
		assert.Equal(t, "string", s)
	}
}

func TestContainer_To_Import(t *testing.T) {
	c1 := func(c *Container) {
		c.Register(func() int { return 1 })
	}
	c2 := func(c *Container) {
		c.Import(c1)
		c.Register(func(i int) string { return fmt.Sprint(i) })
	}
	var i int
	if assert.NoError(t, Init(c2).To(&i)) {
		assert.Equal(t, 1, i)
	}
}

func TestContainer_Decorate(t *testing.T) {
	c1 := func(c *Container) {
		c.Register(func() int { return 1 })
	}
	c2 := func(c *Container) {
		c.Import(c1)
		c.Register(func(i int) string { return fmt.Sprint(i) })
	}
	c3 := func(c *Container) {
		c.Import(c2)
		c.Decorate(func(i string) string { return i + "!" })
	}
	var i string
	if assert.NoError(t, Init(c3).To(&i)) {
		assert.Equal(t, "1!", i)
	}
}

func TestContainer_ImportWithForce(t *testing.T) {
	c1 := func(c *Container) {
		c.Register(func() int { return 1 })
	}
	c2 := func(c *Container) {
		c.Register(func() int { return 2 })
	}
	c3 := func(c *Container) {
		c.Import(c1)
		c.ImportWithForce(c2)
	}
	var i int
	if assert.NoError(t, Init(c3).To(&i)) {
		assert.Equal(t, 2, i)
	}
}

func TestContainer_RegisterWithForce(t *testing.T) {
	c := func(c *Container) {
		c.Register(func() int { return 1 })
		c.RegisterWithForce(func() int { return 2 })
	}
	var i int
	if assert.NoError(t, Init(c).To(&i)) {
		assert.Equal(t, 2, i)
	}
}

func TestContainer_ImportRegisterWithForce(t *testing.T) {
	c1 := func(c *Container) {
		c.Register(func() int { return 1 })
	}
	c2 := func(c *Container) {
		c.RegisterWithForce(func() int { return 2 })
	}
	c3 := func(c *Container) {
		c.Import(c1)
		c.ImportWithForce(c2)
		c.RegisterWithForce(func() int { return 3 })
	}
	var i int
	if assert.NoError(t, Init(c3).To(&i)) {
		assert.Equal(t, 3, i)
	}
}

func TestContainer_ImportDuplicateAlternativeImplementations(t *testing.T) {
	type Something interface{}

	type SomethingV1 struct{}
	type SomethingV2 struct{}

	c1 := func(c *Container) {
		c.Register(func() Something { return &SomethingV1{} })
	}

	c2 := func(c *Container) {
		c.Register(func() Something { return &SomethingV2{} })
	}

	c3 := func(c *Container) {
		c.Import(c1)
		c.Import(c2)
	}

	var something Something
	err := Init(c3).To(&something)
	assert.Error(t, err)
}
