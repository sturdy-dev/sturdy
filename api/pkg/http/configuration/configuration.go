package configuration

import "getsturdy.com/api/pkg/configuration/flags"

type Configuration struct {
	Addr             flags.Addr `long:"addr" description:"Address to listen on" default:"localhost:3000"`
	AllowCORSOrigins []string   `long:"allow-cors-origin" description:"Additional origin that is allowed to make CORS requests (can be provided multiple times)"`
}
