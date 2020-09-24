package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Yapo/logger"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/repository"
)

type httpHandler struct {
	logger loggers.Logger
}

// NewHTTPHandler will create a new instance of a custom http request handler
func NewHTTPHandler(logger loggers.Logger) repository.HTTPHandler {
	return &httpHandler{
		logger: logger,
	}
}

// Send will execute the sending of a http request
// a custom http client has been made to add a request timeout of 10 seconds
func (h *httpHandler) Send(req repository.HTTPRequest) (interface{}, error) {
	h.logger.Debug("Http - %s - Sending HTTP request to: %+v", req.GetMethod(), req.GetPath())

	// this makes a custom http client with a timeout in secs for each request
	var httpClient = &http.Client{
		Timeout: time.Second * req.(*request).timeOut,
	}
	resp, err := httpClient.Do(&req.(*request).innerRequest)
	if err != nil {
		h.logger.Error("Http - %s - Error sending HTTP request: %+v", req.GetMethod(), err)
		return "", fmt.Errorf("found error: %+v", err)
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if isErrorCode(resp.StatusCode) {
		h.logger.Error("Http - %s - Received an error response: %+v", req.GetMethod(), err)
		var msg interface{}
		if e := json.Unmarshal(response, &msg); e != nil {
			return "", fmt.Errorf("the error code was %d", resp.StatusCode)
		}
		return "", fmt.Errorf("%s", msg)
	}
	if err != nil {
		h.logger.Error("Http - %s - Error reading response: %+v", req.GetMethod(), err)
	}
	return string(response), nil
}

// request is a custom golang http.Request
type request struct {
	innerRequest http.Request
	body         interface{}
	timeOut      time.Duration
	logger       loggers.Logger
}

// NewRequest returns an initialized struct that can be used to make a http request
func (*httpHandler) NewRequest() repository.HTTPRequest {
	return &request{
		innerRequest: http.Request{
			Header: make(http.Header),
		},
		timeOut: time.Duration(10),
	}
}

// SetMethod sets the http method to be used, like GET, POST, PUT, etc
func (r *request) SetMethod(method string) repository.HTTPRequest {
	r.innerRequest.Method = method
	return r
}

// GetMethod retrieves the actual http method
func (r *request) GetMethod() string {
	return r.innerRequest.Method
}

// SetPath sets the url path that will be requested
func (r *request) SetPath(path string) repository.HTTPRequest {
	url, err := url.Parse(path)
	r.innerRequest.URL = url
	if err != nil {
		r.logger.Error("Http - there was an error setting the request path: %+v", err)
	}
	return r
}

// GetPath retrieves the actual url path
func (r request) GetPath() string {
	return r.innerRequest.URL.String()
}

// SetHeaders will set custom headers to the request
func (r *request) SetHeaders(headers map[string]string) repository.HTTPRequest {
	for header, value := range headers {
		r.innerRequest.Header.Set(header, value)
	}
	return r
}

// GetHeaders will retrieve the custom headers of the request
func (r *request) GetHeaders() map[string][]string {
	return r.innerRequest.Header
}

// SetBody will set a custom body to the request, this body is the json representation of an interface{}
// this method will also set the custom header Content-type to application-json
// and will save the original body
func (r *request) SetBody(body interface{}) repository.HTTPRequest {
	var reader io.Reader
	var length int

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			logger.Error("Http - Error parsing request data to json: %+v", err)
		}
		reader = strings.NewReader(string(jsonBody))
		length = len(jsonBody)
	}
	// if SetBody is called then we add the Content-type header as a default
	r.SetHeaders(map[string]string{"Content-type": "application/json"})
	r.innerRequest.Body = ioutil.NopCloser(reader)
	r.innerRequest.ContentLength = int64(length)

	// this will be useful if we need to call GetBody(...)
	r.body = body
	return r
}

// GetBody retrieves the original interface{} set on this request
// so after calling this method you should be able to assert it to its original type
func (r *request) GetBody() interface{} {
	return r.body
}

// SetQueryParams will set custom query parameters to the request
func (r *request) SetQueryParams(queryParams map[string]string) repository.HTTPRequest {
	q := r.innerRequest.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	r.innerRequest.URL.RawQuery = q.Encode()
	return r
}

// GetQueryParams will retrieve the query parameters of the request
func (r *request) GetQueryParams() map[string][]string {
	return r.innerRequest.URL.Query()
}

// GetTimeout will retrieve the timeout of the request
func (r *request) GetTimeOut() time.Duration {
	return r.timeOut
}

// SetTimeout will set the timeout to the request
func (r *request) SetTimeOut(timeout int) repository.HTTPRequest {
	r.timeOut = time.Duration(timeout)
	return r
}

func isErrorCode(statusCode int) bool {
	return (statusCode >= http.StatusBadRequest &&
		statusCode <= http.StatusNetworkAuthenticationRequired)
}
