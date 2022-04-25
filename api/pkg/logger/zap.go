package logger

import (
	"fmt"
	"os"

	"getsturdy.com/api/pkg/logger/configuration"
	"getsturdy.com/api/pkg/metrics/zapprometheus"
	"getsturdy.com/api/pkg/version"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

var (
	options = []zap.Option{
		zap.Hooks(
			zapprometheus.Hook,
		),
	}

	consoleDebugging = zapcore.Lock(os.Stdout)
	consoleErrors    = zapcore.Lock(os.Stderr)

	// EncoderConfig is the zapcore EncoderConfig used to configure the
	// encoder for production.
	encoderConfigByEnv = map[bool]func() zapcore.EncoderConfig{
		true:  zap.NewProductionEncoderConfig,
		false: zap.NewDevelopmentEncoderConfig,
	}

	// encoderByTerminal is a map of whether the terminal is a TTY to a function
	// that returns a zapcore.EncoderConfig.
	encoderByTerminal = map[bool]func(zapcore.EncoderConfig) zapcore.Encoder{
		true:  zapcore.NewConsoleEncoder,
		false: zapcore.NewJSONEncoder,
	}
)

func New(cfg *configuration.Configuration, sentryClient *sentry.Client) (*zap.Logger, error) {
	encoderConfig := encoderConfigByEnv[cfg.Production]()
	isTeminal := term.IsTerminal(int(os.Stdout.Fd()))
	encoder := encoderByTerminal[isTeminal](encoderConfig)

	loglevel := zapcore.InfoLevel
	switch cfg.Level {
	case "INFO":
		loglevel = zapcore.InfoLevel
	case "WARN":
		loglevel = zapcore.WarnLevel
	case "ERROR":
		loglevel = zapcore.ErrorLevel
	default:
		return nil, fmt.Errorf("unexpected log level: %s (must be one of INFO, WARN, or ERROR)", cfg.Level)
	}

	atomicLevel := zap.NewAtomicLevelAt(loglevel)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return atomicLevel.Enabled(lvl) && lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return atomicLevel.Enabled(lvl) && lvl < zapcore.ErrorLevel
	})

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
		zapcore.NewCore(encoder, consoleErrors, highPriority),
	}
	if cfg.Production {
		core, err := zapsentry.NewCore(zapsentry.Configuration{
			Level:             zapcore.ErrorLevel,
			EnableBreadcrumbs: true,
			BreadcrumbLevel:   zapcore.InfoLevel,
		}, zapsentry.NewSentryClientFromClient(sentryClient))
		if err != nil {
			return nil, fmt.Errorf("failed to create zapsentry core: %w", err)
		}
		cores = append(cores, core)
	}

	return zap.New(zapcore.NewTee(cores...), options...).With(
		zap.String("version", version.Version),
		zap.Stringer("environment", version.Type),
	), nil
}
