package infrastructure

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gorilla/mux.v1"
)

func TestQueryParamsOK(t *testing.T) {
	type input struct {
		ID string `query:"id"`
	}

	result := input{}
	expected := input{"1"}
	r := httptest.NewRequest("GET", "/api/v1?id=1", nil)

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromQuery()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestPathOK(t *testing.T) {
	type input struct {
		ID string `path:"id"`
	}

	result := input{}
	expected := input{"1"}
	r := httptest.NewRequest("GET", "/api/v1/1", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromPath()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestJsonBodyOK(t *testing.T) {
	type input struct {
		ID string `json:"id"`
	}

	result := input{}
	expected := input{"edgar"}
	r := httptest.NewRequest("POST", "/api/v1/", strings.NewReader(`{"id": "edgar"}`))

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromJSONBody()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestHeadersOK(t *testing.T) {
	type input struct {
		ID string `headers:"Id"`
	}

	result := input{}
	expected := input{"edgar"}
	r := httptest.NewRequest("POST", "/api/v1/", nil)
	r.Header.Add("Id", "edgar")

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromHeaders()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestAllSourcesAtSameTimeOK(t *testing.T) {
	type input struct {
		HeaderID string `headers:"Id"`
		BodyID   string `json:"id"`
		PathID   string `path:"id"`
		QueryID  string `query:"id"`
	}

	result := input{}
	expected := input{"edgar", "edgod", "edgugu", "edgarda"}
	r := httptest.NewRequest("POST", "/api/v1/edgugu?id=edgarda", strings.NewReader(`{"id": "edgod"}`))
	r.Header.Add("Id", "edgar")
	r = mux.SetURLVars(r, map[string]string{
		"id": "edgugu",
	})

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromHeaders().FromJSONBody().FromPath().FromQuery()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestOverlapSourcesOK(t *testing.T) {
	type input struct {
		BodyID string `json:"id" path:"id"`
		PathID string `path:"id"`
	}

	result := input{}
	expected := input{"edgugu", "edgugu"}
	r := httptest.NewRequest("POST", "/api/v1/edgugu", strings.NewReader(`{"id": "edgod"}`))
	r = mux.SetURLVars(r, map[string]string{
		"id": "edgugu",
	})

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromJSONBody().FromPath()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestRawBodyOK(t *testing.T) {
	type input struct {
		BodyID []byte `raw:"body"`
	}

	result := input{}
	expected := input{BodyID: []byte(`{"id": "edgod"}`)}
	r := httptest.NewRequest("POST", "/api/v1/edgugu", strings.NewReader(`{"id": "edgod"}`))

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromRawBody()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestEmptySourceOK(t *testing.T) {
	type input struct {
		BodyID string `json:"id"`
		PathID string `path:"id"`
	}

	result := input{}
	expected := input{BodyID: "edgod"}
	r := httptest.NewRequest("POST", "/api/v1/edgugu", strings.NewReader(`{"id": "edgod"}`))

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromJSONBody().FromPath()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}
