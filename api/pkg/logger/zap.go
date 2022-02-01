package logger

import (
	"flag"

	"getsturdy.com/api/pkg/metrics/zapprometheus"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	productionLogger = flag.Bool("production-logger", false, "")
)

func New(sentryClient *raven.Client) (*zap.Logger, error) {
	options := []zap.Option{
		zap.Hooks(
			zapprometheus.Hook,
		),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, &sentryCore{
				LevelEnabler: zapcore.ErrorLevel,
				sentryClient: sentryClient,
			})
		}),
	}

	if *productionLogger {
		return zap.NewProduction(options...)
	}
	return zap.NewDevelopment(options...)
}
