package integrations

import (
	"fmt"
)

// Known errors.
var (
	ErrNotFound = fmt.Errorf("provider not found")
)

var registry = map[ProviderType]Provider{}

func Register(name ProviderType, p Provider) {
	registry[name] = p
}

func Get(name ProviderType) (Provider, error) {
	if provider, ok := registry[name]; ok {
		return provider, nil
	}
	return nil, ErrNotFound
}
