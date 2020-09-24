package loggers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Yapo/goutils"
)

// There are no return values to assert on, as logger only cause side effects
// to communicate with the outside world. These tests only ensure that the
// loggers don't panic

func TestJSONHandlerLogger(t *testing.T) {
	m := &loggerMock{t: t}
	r := httptest.NewRequest("GET", "/test", strings.NewReader(""))
	l := MakeJSONHandlerLogger(m)
	l.LogRequestStart(r)
	l.LogRequestEnd(r, &goutils.Response{}, "")
	l.LogRequestPanic(r, &goutils.Response{}, nil)
}
