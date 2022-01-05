package ctxlog

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

func IsError(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	return true
}

// ErrorOrWarn logs some "safe" errors as Warn, and all others as Error
// An example of a "safe" error is context.Cancelled
func ErrorOrWarn(logger *zap.Logger, msg string, err error) {
	l := logger.WithOptions(zap.AddCallerSkip(1))
	if IsError(err) {
		l.Error(msg, zap.Error(err))
	} else {
		l.Warn(msg, zap.Error(err))
	}
}
