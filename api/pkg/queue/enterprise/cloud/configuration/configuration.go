package configuration

type Configuration struct {
	Local    bool   `long:"local" description:"Use in-memory queue instead of SQS"`
	Hostname string `long:"hostname" description:"Hostname of the queue"`
	Prefix   string `long:"prefix" description:"Prefix for queue names"`
}
