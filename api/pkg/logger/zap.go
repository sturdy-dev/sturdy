package logger

import (
	"getsturdy.com/api/pkg/metrics/zapprometheus"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Configuration struct {
	Production bool `long:"production" description:"Production mode"`
}

func New(cfg *Configuration, sentryClient *raven.Client) (*zap.Logger, error) {
	options := []zap.Option{
		zap.Hooks(
			zapprometheus.Hook,
		),
	}

	if cfg.Production {
		options = append(options, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, &sentryCore{
				LevelEnabler: zapcore.ErrorLevel,
				sentryClient: sentryClient,
			})
		}))
		return zap.NewProduction(options...)
	}
	return zap.NewDevelopment(options...)
}
