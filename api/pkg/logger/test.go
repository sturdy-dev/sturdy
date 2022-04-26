package logger

import (
	"io"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger is the minimal subset of logging functions that exists both in
// testing.T and testing.B.
type logger interface {
	Log(args ...interface{})
}

// testOutput is a zapcore.SyncWriter that writes all output to l.
type testOutput struct {
	logger
}

// writeSyncer decorates an io.Writer with a no-op Sync() function.
type writeSyncer struct {
	io.Writer
}

func NewTest(t *testing.T) *zap.Logger {
	conf := zap.NewDevelopmentEncoderConfig()
	conf.TimeKey = ""
	enc := zapcore.NewConsoleEncoder(conf)
	core := zapcore.NewCore(enc, writeSyncer{testOutput{t}}, zap.DebugLevel)
	return zap.New(core)
}

// Write logs all messages as logs via o.
func (o testOutput) Write(p []byte) (int, error) {
	msg := strings.TrimSpace(string(p))
	o.Log(msg)
	return len(p), nil
}

// Sync does nothing since all output was written to the writer immediately.
func (ws writeSyncer) Sync() error {
	return nil
}
