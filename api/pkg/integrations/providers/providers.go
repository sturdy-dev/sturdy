package providers

import (
	"fmt"
)

// Known errors.
var (
	ErrNotFound = fmt.Errorf("provider not found")
)

type Providers map[ProviderName]Provider

func (p Providers) Get(name ProviderName) (res Provider, err error) {
	if provider, ok := p[name]; ok {
		return provider, nil
	}
	err = ErrNotFound
	return
}
