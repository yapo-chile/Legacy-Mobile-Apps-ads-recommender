package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var t = time.Now()
var today = fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
var ufPath = ""
var defaultValue = 30000

func TestNewIndicatorsRepository(t *testing.T) {
	mHTTPCachedHandler := new(MockHTTPCachedHandler)
	indicatorsRepository := &indicatorsRepository{
		HTTPCachedHandler: mHTTPCachedHandler,
		UFPath:            ufPath,
		DefaultValue:      float64(defaultValue),
	}
	repository := NewIndicatorsRepository(mHTTPCachedHandler, ufPath, defaultValue)
	assert.Equal(t, indicatorsRepository, repository)
	mHTTPCachedHandler.AssertExpectations(t)
}

func TestGetUFOK(t *testing.T) {
	expectedResult := 29095.61
	// nolint: misspell
	response := `{
			"version":"1.6.0",
			"autor":"mindicador.cl",
			"codigo":"uf",
			"nombre":"Unidad de fomento (UF)",
			"unidad_medida":"Pesos",
			"serie":[{
				"fecha":"2021-01-21T03:00:00.000Z",
				"valor":29095.61
			}]
		}`

	mHTTPCachedHandler := new(MockHTTPCachedHandler)
	mHTTPRequest := new(mockRequest)
	mHTTPCachedHandler.On("NewRequest").Return(mHTTPRequest, nil)
	mHTTPRequest.On("SetPath", ufPath+today).Return(mHTTPRequest)
	mHTTPRequest.On("SetMethod", "GET").Return(mHTTPRequest)
	mHTTPCachedHandler.On("Send", mHTTPRequest).Return(response, nil)
	indicatorsRepository := &indicatorsRepository{
		HTTPCachedHandler: mHTTPCachedHandler,
		UFPath:            ufPath,
	}
	result, err := indicatorsRepository.GetUF()
	assert.Equal(t, result, expectedResult)
	assert.NoError(t, err)
	mHTTPCachedHandler.AssertExpectations(t)
	mHTTPRequest.AssertExpectations(t)
}

func TestGetUFError(t *testing.T) {
	var expectedResult float64
	// nolint: misspell
	response := `{
			"version":"1.6.0",
			"autor":"mindicador.cl",
			"codigo":"uf",
			"nombre":"Unidad de fomento (UF)",
			"unidad_medida":"Pesos",
			"serie":[{
				"fecha":"2021-01-21T03:00:00.000Z",
				"valor":29095.61
			}]
		}`

	mHTTPCachedHandler := new(MockHTTPCachedHandler)
	mHTTPRequest := new(mockRequest)
	mHTTPCachedHandler.On("NewRequest").Return(mHTTPRequest, nil)
	mHTTPRequest.On("SetPath", ufPath+today).Return(mHTTPRequest)
	mHTTPRequest.On("SetMethod", "GET").Return(mHTTPRequest)
	mHTTPCachedHandler.On("Send", mHTTPRequest).Return(response, fmt.Errorf(""))
	indicatorsRepository := &indicatorsRepository{
		HTTPCachedHandler: mHTTPCachedHandler,
		UFPath:            ufPath,
	}
	result, err := indicatorsRepository.GetUF()
	assert.Equal(t, result, expectedResult)
	assert.Error(t, err)
	mHTTPCachedHandler.AssertExpectations(t)
	mHTTPRequest.AssertExpectations(t)
}

func TestGetUFSetEmpty(t *testing.T) {
	expectedResult := 0.0
	// nolint: misspell
	response := `{
			"version":"1.6.0",
			"autor":"mindicador.cl",
			"codigo":"uf",
			"nombre":"Unidad de fomento (UF)",
			"unidad_medida":"Pesos",
			"serie":[]
		}`

	mHTTPCachedHandler := new(MockHTTPCachedHandler)
	mHTTPRequest := new(mockRequest)
	mHTTPCachedHandler.On("NewRequest").Return(mHTTPRequest, nil)
	mHTTPRequest.On("SetPath", ufPath+today).Return(mHTTPRequest)
	mHTTPRequest.On("SetMethod", "GET").Return(mHTTPRequest)
	mHTTPCachedHandler.On("Send", mHTTPRequest).Return(response, nil)
	indicatorsRepository := &indicatorsRepository{
		HTTPCachedHandler: mHTTPCachedHandler,
		UFPath:            ufPath,
	}
	result, err := indicatorsRepository.GetUF()
	assert.Equal(t, result, expectedResult)
	assert.Error(t, err)
	mHTTPCachedHandler.AssertExpectations(t)
	mHTTPRequest.AssertExpectations(t)
}
