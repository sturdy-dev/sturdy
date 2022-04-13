package configuration

type Configuration struct {
	Production bool   `long:"production" description:"Production mode"`
	Level      string `long:"level" default:"WARN" description:"Log level (INFO, WARN, ERROR)"`
}
