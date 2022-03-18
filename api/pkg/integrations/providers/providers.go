package providers

import (
	"fmt"
)

// Known errors.
var (
	ErrNotFound = fmt.Errorf("provider not found")
)

var registry = map[ProviderName]Provider{}

func Register(p Provider) {
	registry[p.ProviderName()] = p
}

func Get[T any](name ProviderName) (res T, err error) {
	if provider, ok := registry[name]; ok {
		if t, ok := provider.(T); ok {
			return t, nil
		}
	}
	err = ErrNotFound
	return
}
