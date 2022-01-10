package ip

import (
	"context"
	"net"
)

type ipKeyType struct{}

var ipKey = ipKeyType{}

func NewContext(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, ipKey, &ip)
}

func FromContext(ctx context.Context) (*net.IP, bool) {
	s, ok := ctx.Value(ipKey).(*net.IP)
	return s, ok
}
