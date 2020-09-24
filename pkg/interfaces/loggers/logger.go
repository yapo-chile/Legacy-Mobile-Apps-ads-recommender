package loggers

// Logger is an interface for logging facilities
type Logger interface {
	Debug(format string, params ...interface{})
	Info(format string, params ...interface{})
	Warn(format string, params ...interface{})
	Error(format string, params ...interface{})
	Crit(format string, params ...interface{})
	Success(format string, params ...interface{})
}
