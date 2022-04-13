package configuration

type Configuration struct {
	Region string `long:"region" description:"AWS region to use" default:"eu-north-1"`
}
