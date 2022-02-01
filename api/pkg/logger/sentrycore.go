package logger

import (
	"getsturdy.com/api/pkg/version"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Core = (*sentryCore)(nil)

// sentryCore send event to sentry.
type sentryCore struct {
	zapcore.LevelEnabler
	sentryClient *raven.Client

	fields []zapcore.Field
}

func level(lvl zapcore.Level) raven.Severity {
	switch lvl {
	case zapcore.DebugLevel:
		return raven.DEBUG
	case zapcore.InfoLevel:
		return raven.INFO
	case zapcore.WarnLevel:
		return raven.WARNING
	case zapcore.ErrorLevel:
		return raven.ERROR
	default:
		return raven.FATAL
	}
}

func (c *sentryCore) clone() *sentryCore {
	v := *c
	return &v
}

func (c *sentryCore) With(fields []zapcore.Field) zapcore.Core {
	ret := c.clone()
	ret.fields = append(c.fields, fields...)
	return ret
}

func (c *sentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		ce.AddCore(ent, c)
	}
	return ce
}

func (c *sentryCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	event := &raven.Packet{}
	event.Message = entry.Message
	event.Timestamp = raven.Timestamp(entry.Time)
	event.Logger = entry.LoggerName
	event.Level = level(entry.Level)
	enc := zapcore.NewMapObjectEncoder()

	var fieldErr error
	for _, i := range append(c.fields, fields...) {
		if i.Type == zapcore.ErrorType && fieldErr == nil {
			if err, ok := i.Interface.(error); ok && err != nil {
				fieldErr = err
				continue
			}
		}
		i.AddTo(enc)
	}
	event.Extra = enc.Fields
	event.Release = version.Version
	event.Environment = version.Type.String()
	trace := raven.NewStacktrace(3, 3, []string{"getsturdy.com/api"})
	if fieldErr != nil {
		event.Interfaces = append(event.Interfaces, &raven.Exception{
			Type:       event.Message,
			Value:      fieldErr.Error(),
			Stacktrace: trace,
		})
	} else if trace != nil {
		event.Interfaces = append(event.Interfaces, &raven.Exception{
			Value:      event.Message,
			Stacktrace: trace,
		})
	}
	event.ServerName = "cloud"

	c.sentryClient.Capture(event, map[string]string{})
	return nil
}

func (c *sentryCore) Sync() error {
	return nil
}
