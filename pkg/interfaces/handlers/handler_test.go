package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mux "gopkg.in/gorilla/mux.v1"

	"github.com/Yapo/goutils"
)

func MakeMockInputGetter(input HandlerInput, response *goutils.Response) InputGetter {
	return func() (HandlerInput, *goutils.Response) {
		return input, response
	}
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Input(ir InputRequest) HandlerInput {
	args := m.Called(ir)
	return args.Get(0).(HandlerInput)
}

func (m *MockHandler) Execute(getter InputGetter) *goutils.Response {
	args := m.Called(getter)
	_, response := getter()
	if response != nil {
		return response
	}
	return args.Get(0).(*goutils.Response)
}

type MockInputHandler struct {
	mock.Mock
}

func (m *MockInputHandler) Input() (HandlerInput, *goutils.Response) {
	args := m.Called()
	return args.Get(0).(HandlerInput), args.Get(1).(*goutils.Response)
}

func (m *MockInputHandler) NewInputRequest(r *http.Request) InputRequest {
	args := m.Called(r)
	return args.Get(0).(InputRequest)
}

func (m *MockInputHandler) SetInputRequest(ri InputRequest, hi HandlerInput) {
	m.Called(ri, hi)
}

type MockPanicHandler struct {
	mock.Mock
}

func (m *MockPanicHandler) Input(ir InputRequest) HandlerInput {
	args := m.Called(ir)
	return args.Get(0).(HandlerInput)
}
func (m *MockPanicHandler) Execute(getter InputGetter) *goutils.Response {
	m.Called(getter)
	panic("dead")
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogRequestStart(r *http.Request) {
	m.Called(r)
}
func (m *MockLogger) LogRequestEnd(r *http.Request, response *goutils.Response, cacheStatus string) {
	m.Called(r, response, cacheStatus)
}
func (m *MockLogger) LogRequestPanic(r *http.Request, response *goutils.Response, err interface{}) {
	m.Called(r, response, err)
}

type DummyInput struct {
	X int
}

type DummyOutput struct {
	Y string
}

type TestParam struct {
	Param1 string `get:"param1"`
	Param2 string `get:"param2"`
}

type TestParamInt struct {
	Param1 int `get:"param3"`
	Param2 int `get:"param4"`
}

type TestParamStruct struct {
	Param1 TestParam    `get:"param5"`
	Param2 TestParamInt `get:"param6"`
}

type MockCors struct {
	mock.Mock
}

func (mc *MockCors) GetHeaders() map[string]string {
	args := mc.Called()
	return args.Get(0).(map[string]string)
}

type MockCache struct {
	mock.Mock
}

func (cache *MockCache) Validate(w http.ResponseWriter, r *http.Request) bool {
	args := cache.Called()
	return args.Get(0).(bool)
}

type MockRequestCache struct {
	mock.Mock
}

func (cache *MockRequestCache) GetCache(input interface{}) (*goutils.Response, error) {
	args := cache.Called(input)
	return args.Get(0).(*goutils.Response), args.Error(1)
}

func (cache *MockRequestCache) SetCache(input interface{}, response *goutils.Response) error {
	args := cache.Called(input, response)
	return args.Error(0)
}

func TestJsonHandlerFuncOK(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}

	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response, mock.AnythingOfType("string"))

	mCache := MockCache{}
	mCache.On("Validate").Return(false)
	mRequestCache := MockRequestCache{}
	mRequestCache.On("GetCache", mock.AnythingOfType("*handlers.DummyInput")).Return(response, fmt.Errorf(""))
	mRequestCache.On("SetCache", mock.AnythingOfType("*handlers.DummyInput"), response).Return(fmt.Errorf(""))
	mC := MockCors{}
	mC.On("GetHeaders").Return(map[string]string{})
	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mC, &mCache, &mRequestCache)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	l.AssertExpectations(t)
	mC.AssertExpectations(t)
	mCache.AssertExpectations(t)
	mRequestCache.AssertExpectations(t)
}

func TestJsonHandlerFuncOK2(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl?id=1,2", strings.NewReader("{}"))
	r = mux.SetURLVars(r, map[string]string{
		"id": "1, 2",
	})

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response, mock.AnythingOfType("string"))

	mC := MockCors{}
	mC.On("GetHeaders").Return(map[string]string{})

	mCache := MockCache{}
	mCache.On("Validate").Return(false)
	mRequestCache := MockRequestCache{}
	mRequestCache.On("GetCache", mock.AnythingOfType("*handlers.DummyInput")).Return(response, fmt.Errorf(""))
	mRequestCache.On("SetCache", mock.AnythingOfType("*handlers.DummyInput"), response).Return(fmt.Errorf(""))
	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mC, &mCache, &mRequestCache)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	l.AssertExpectations(t)
	mC.AssertExpectations(t)
	mCache.AssertExpectations(t)
	mRequestCache.AssertExpectations(t)
}

func TestJsonHandlerFuncParseError(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	getter := mock.AnythingOfType("handlers.InputGetter")
	response := &goutils.Response{
		Code: 400,
		Body: struct{ ErrorMessage string }{ErrorMessage: "unexpected EOF"},
	}
	h.On("Execute", getter)
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/someurl", strings.NewReader("{"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, mock.AnythingOfType("*goutils.Response"), mock.AnythingOfType("string"))

	mC := MockCors{}
	mC.On("GetHeaders").Return(map[string]string{})

	mCache := MockCache{}
	mCache.On("Validate").Return(false)
	mRequestCache := MockRequestCache{}
	mRequestCache.On("GetCache", mock.AnythingOfType("*handlers.DummyInput")).Return(response, fmt.Errorf(""))
	mRequestCache.On("SetCache", mock.AnythingOfType("*handlers.DummyInput"), response).Return(fmt.Errorf(""))
	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mC, &mCache, &mRequestCache)
	fn(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"ErrorMessage":"unexpected EOF"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	l.AssertExpectations(t)
	mC.AssertExpectations(t)
	mCache.AssertExpectations(t)
	mRequestCache.AssertExpectations(t)
}

func TestJsonHandlerFuncPanic(t *testing.T) {
	h := MockPanicHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	getter := mock.AnythingOfType("handlers.InputGetter")
	input := &DummyInput{}
	h.On("Execute", getter)
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{"))

	l.On("LogRequestStart", r)
	l.On("LogRequestPanic", r, mock.AnythingOfType("*goutils.Response"), "dead")

	mC := MockCors{}
	mC.On("GetHeaders").Return(map[string]string{})

	mCache := MockCache{}
	mCache.On("Validate").Return(false)
	mRequestCache := MockRequestCache{}
	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mC, &mCache, &mRequestCache)
	fn(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "null"+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	l.AssertExpectations(t)
	mC.AssertExpectations(t)
	mCache.AssertExpectations(t)
	mRequestCache.AssertExpectations(t)
}
func TestJsonHandlerFuncHeaders(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl?id=1,2", strings.NewReader("{}"))
	r = mux.SetURLVars(r, map[string]string{
		"id": "1, 2",
	})

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response, mock.AnythingOfType("string"))

	mC := MockCors{}
	headers := map[string]string{
		"Origin":  "myorigin",
		"Methods": "mistherious",
	}
	mC.On("GetHeaders").Return(headers)

	mCache := MockCache{}
	mCache.On("Validate").Return(false)
	mRequestCache := MockRequestCache{}
	mRequestCache.On("GetCache", mock.AnythingOfType("*handlers.DummyInput")).Return(response, fmt.Errorf(""))
	mRequestCache.On("SetCache", mock.AnythingOfType("*handlers.DummyInput"), response).Return(fmt.Errorf(""))
	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mC, &mCache, &mRequestCache)
	fn(w, r)

	expectedHeaders := http.Header{
		"Access-Control-Allow-Methods": []string{"mistherious"},
		"Access-Control-Allow-Origin":  []string{"myorigin"},
		"Content-Type":                 []string{"application/json"}}

	assert.Equal(t, expectedHeaders, w.HeaderMap) //nolint: staticcheck
	assert.Equal(t, 42, w.Code)
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	l.AssertExpectations(t)
	mC.AssertExpectations(t)
	mCache.AssertExpectations(t)
	mRequestCache.AssertExpectations(t)
}
func TestJsonHandlerFuncCache(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 304,
	}
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl?id=1,2", strings.NewReader("{}"))
	r = mux.SetURLVars(r, map[string]string{
		"id": "1, 2",
	})

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response, mock.AnythingOfType("string"))

	mC := MockCors{}
	mC.On("GetHeaders").Return(map[string]string{})
	mCache := MockCache{}
	mCache.On("Validate").Return(true)
	mRequestCache := MockRequestCache{}
	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mC, &mCache, &mRequestCache)
	r.Header.Add("If-None-Match", "\"123\"")
	fn(w, r)

	expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}

	assert.Equal(t, expectedHeaders, w.HeaderMap) //nolint: staticcheck
	assert.Equal(t, 304, w.Code)
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	l.AssertExpectations(t)
	mC.AssertExpectations(t)
	mCache.AssertExpectations(t)
	mRequestCache.AssertExpectations(t)
}
