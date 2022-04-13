package configuration

type Configuration struct {
	PublicAPIHostname string `long:"public-api-hostname" description:"Public API hostname. Used to fetch codebases from CI"`
}
