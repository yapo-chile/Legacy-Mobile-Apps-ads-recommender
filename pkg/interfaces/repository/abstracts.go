package repository

import (
	"time"
)

// HTTPRequest interface represents the request that is going to be sent via HTTP
type HTTPRequest interface {
	GetMethod() string
	SetMethod(string) HTTPRequest
	GetPath() string
	SetPath(string) HTTPRequest
	GetBody() interface{}
	SetBody(interface{}) HTTPRequest
	GetHeaders() map[string][]string
	SetHeaders(map[string]string) HTTPRequest
	GetQueryParams() map[string][]string
	SetQueryParams(map[string]string) HTTPRequest
	GetTimeOut() time.Duration
	SetTimeOut(int) HTTPRequest
}

// HTTPHandler implements HTTP handler operations
type HTTPHandler interface {
	Send(HTTPRequest) (interface{}, error)
	NewRequest() HTTPRequest
}

// ElasticSearchHandler defines the methods that are available for a elastic search handler
type ElasticSearchHandler interface {
	Info() (interface{}, error)
	Create(index string) error
	PutMapping(mapping []byte, index string) error
	Search(index, query string, size, from int) (string, error)
}

// DataMapping allows get specific configuration params from etcd
type DataMapping interface {
	Get(string) string
}
