package loggers

import (
	"net/http"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/handlers"
)

type jsonHandlerDefaultLogger struct {
	logger Logger
}

func (l *jsonHandlerDefaultLogger) LogRequestStart(r *http.Request) {
	l.logger.Info("< %s %s %s", r.RemoteAddr, r.Method, r.URL)
}

func (l *jsonHandlerDefaultLogger) LogRequestEnd(r *http.Request, response *goutils.Response, cacheStatus string) {
	l.logger.Info("> %s %s %s (%d)%s", r.RemoteAddr, r.Method, r.URL, response.Code, cacheStatus)
}

func (l *jsonHandlerDefaultLogger) LogRequestPanic(r *http.Request, response *goutils.Response, err interface{}) {
	l.logger.Error("> %s %s %s (%d): %s", r.RemoteAddr, r.Method, r.URL, response.Code, err)
}

// MakeJSONHandlerLogger sets up a JsonHandlerLogger instrumented
// via the provided logger
func MakeJSONHandlerLogger(logger Logger) handlers.JSONHandlerLogger {
	return &jsonHandlerDefaultLogger{
		logger: logger,
	}
}
