package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterWithProfiling(t *testing.T) {
	profiling := []bool{false, true}

	for _, with := range profiling {
		maker := RouterMaker{WithProfiling: with}
		router := maker.NewRouter()
		req := httptest.NewRequest("GET", "/debug/pprof/trace", strings.NewReader(""))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		doesMatch := resp.Code == http.StatusOK
		assert.Equal(t, with, doesMatch)
	}
}
