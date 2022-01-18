package zapprometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
)

var (
	logsTotalCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sturdy_logs_total",
		Help: "Number of logs",
	}, []string{"level", "loggerName", "caller", "message"})
)

func Hook(entry zapcore.Entry) error {
	logsTotalCounter.WithLabelValues(entry.Level.String(),
		entry.LoggerName,
		entry.Caller.TrimmedPath(),
		entry.Message).Inc()
	return nil
}
