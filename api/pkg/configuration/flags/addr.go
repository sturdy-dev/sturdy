package flags

import "net"

type Addr struct {
	net.Addr
}

func (a *Addr) UnmarshalFlag(s string) error {
	addr, err := net.ResolveTCPAddr("tcp", s)
	if err != nil {
		return err
	}
	a.Addr = addr
	return nil
}

func (a Addr) MarshalFlag() string {
	return a.String()
}
