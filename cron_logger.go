package metrics

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

type CronLogger struct {
	l *slog.Logger
}

func (c *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error", err)

	c.l.Error(msg,
		c.kvs(keysAndValues...)...,
	)
}

func (c *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.l.Info(msg, c.kvs(keysAndValues...)...)
}

func (c *CronLogger) kvs(keysAndValues ...interface{}) []interface{} {
	if len(keysAndValues)%2 != 0 {
		c.l.Warn("keysAndValues should be in pairs", slog.Any("keysAndValues", keysAndValues))
		return nil
	}

	args := make([]interface{}, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			c.l.Warn("key should be a string", slog.Any("key", keysAndValues[i]))
			continue
		}
		args = append(args, slog.Any(key, keysAndValues[i+1]))
	}

	return args
}

func NewCronLogger(l *slog.Logger) *CronLogger {
	return &CronLogger{l: l}
}

var _ cron.Logger = &CronLogger{}
