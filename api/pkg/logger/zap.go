package logger

import (
	"os"

	"getsturdy.com/api/pkg/metrics/zapprometheus"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

type Configuration struct {
	Production bool `long:"production" description:"Production mode"`
}

var (
	options = []zap.Option{
		zap.Hooks(
			zapprometheus.Hook,
		),
	}

	highPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

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

func New(cfg *Configuration, sentryClient *raven.Client) *zap.Logger {
	encoderConfig := encoderConfigByEnv[cfg.Production]()
	isTeminal := term.IsTerminal(int(os.Stdout.Fd()))
	encoder := encoderByTerminal[isTeminal](encoderConfig)
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
		zapcore.NewCore(encoder, consoleErrors, highPriority),
	}
	if cfg.Production {
		cores = append(cores, &sentryCore{LevelEnabler: highPriority, sentryClient: sentryClient})
	}
	return zap.New(zapcore.NewTee(cores...), options...)
}
