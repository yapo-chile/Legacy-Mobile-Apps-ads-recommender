package handlers

import (
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandlerInput(t *testing.T) {
	var h HealthHandler
	mMockInputRequest := MockInputRequest{}

	input := h.Input(&mMockInputRequest)
	var expected *healthHandlerInput
	assert.IsType(t, expected, input)
}

func TestHealthHandlerRun(t *testing.T) {
	var h HealthHandler
	var input HandlerInput
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: healthRequestOutput{"OK"},
	}

	assert.Equal(t, expected, r)
}
