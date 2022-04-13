package configuration

import "getsturdy.com/api/pkg/configuration/flags"

type Configuration struct {
	Addr flags.Addr `long:"addr" description:"listen address" default:"127.0.0.1:3002"`
}
