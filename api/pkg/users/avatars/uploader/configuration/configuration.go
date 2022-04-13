package configuration

import "getsturdy.com/api/pkg/configuration/flags"

type Configuration struct {
	URL flags.URL `long:"url" description:"Avatars base url" default:"http://127.0.0.1:3000/"`
}
