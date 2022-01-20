package api

import "context"

type Config struct {
	GitListenAddr       string
	HTTPPProfListenAddr string
	MetricsListenAddr   string
	HTTPAddr            string
}

type API interface {
	Start(context.Context, *Config) error
}
