package infrastructure

import (
	"github.com/Yapo/logger"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/loggers"
)

// yapoLogger struct that implements the Logger interface using the Yapo/logger library
type yapoLogger struct {
	metrics EventCollector
}

// MakeYapoLogger creates and sets up a yapo flavored Logger
func MakeYapoLogger(config *LoggerConf, metrics EventCollector) (loggers.Logger, error) {
	log := yapoLogger{
		metrics: metrics,
	}
	err := log.init(config)

	return log, err
}

// Init initialize the logger
func (y *yapoLogger) init(config *LoggerConf) error {
	loggerConf := logger.LogConfig{
		Syslog: logger.SyslogConfig{
			Enabled:  config.SyslogEnabled,
			Identity: config.SyslogIdentity,
		},
		Stdlog: logger.StdlogConfig{
			Enabled: config.StdlogEnabled,
		},
	}
	if err := logger.Init(loggerConf); err != nil {
		return err
	}

	logger.SetLogLevel(config.LogLevel)

	return nil
}

// Debug logs a message at DEBUG level
func (y yapoLogger) Debug(format string, params ...interface{}) {
	logger.Debug(format, params...)
}

// Info logs a message at INFO level.
// Info events are automatically exported to prometheus.
func (y yapoLogger) Info(format string, params ...interface{}) {
	y.metrics.CollectEvent(getEntityName(), getEventName(), getEventType())
	logger.Info(format, params...)
}

// Success logs a message as Success event.
// Success events are automatically exported to prometheus.
func (y yapoLogger) Success(format string, params ...interface{}) {
	y.metrics.CollectEvent(getEntityName(), getEventName(), getEventType())
	logger.Info(format, params...)
}

// Warn logs a message at WARNING level.
// warning events are automatically exported to prometheus.
func (y yapoLogger) Warn(format string, params ...interface{}) {
	y.metrics.CollectEvent(getEntityName(), getEventName(), getEventType())
	logger.Warn(format, params...)
}

// Error logs a message at ERROR level.
// Error events are automatically exported to prometheus.
func (y yapoLogger) Error(format string, params ...interface{}) {
	y.metrics.CollectEvent(getEntityName(), getEventName(), getEventType())
	logger.Error(format, params...)
}

// LogCrit logs a message at CRITICAL level.
// Critical events are automatically exported to prometheus.
func (y yapoLogger) Crit(format string, params ...interface{}) {
	y.metrics.CollectEvent(getEntityName(), getEventName(), getEventType())
	logger.Crit(format, params...)
}
