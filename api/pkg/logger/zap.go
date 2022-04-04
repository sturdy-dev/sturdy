package logger

import (
	"fmt"
	"os"

	"getsturdy.com/api/pkg/metrics/zapprometheus"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

type Configuration struct {
	Production bool   `long:"production" description:"Production mode"`
	Level      string `long:"level" default:"WARN" description:"Log level (INFO, WARN, ERROR)"`
}

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

func New(cfg *Configuration, sentryClient *raven.Client) (*zap.Logger, error) {
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
		cores = append(cores, &sentryCore{LevelEnabler: highPriority, sentryClient: sentryClient})
	}

	return zap.New(zapcore.NewTee(cores...), options...), nil
}
