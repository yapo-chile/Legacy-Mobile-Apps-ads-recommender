package infrastructure

import (
	"crypto/md5" //nolint: gosec
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anevsky/cachego/memory"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/interfaces/loggers"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/interfaces/repository"
)

type httpCachedHandler struct {
	logger   loggers.Logger
	cache    memory.CACHE
	cacheTTL int
}

// NewHTTPCachedHandler will create a new instance of a custom http cached request handler
func NewHTTPCachedHandler(logger loggers.Logger, ttl int) repository.HTTPHandler {
	return &httpCachedHandler{
		logger:   logger,
		cache:    memory.Alloc(),
		cacheTTL: ttl,
	}
}

func (h *httpCachedHandler) getHash(data interface{}) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%v", data))) //nolint: gosec
	return fmt.Sprintf("%x", sum)
}

func (h *httpCachedHandler) setCache(hash, response string) error {
	if err := h.cache.SetString(hash, response); err != nil {
		return err
	}
	return h.cache.SetTTL(hash, h.cacheTTL)
}

func (h *httpCachedHandler) getCache(hash string) (string, error) {
	v, err := h.cache.Get(hash)
	return v.(string), err
}

// Send will execute the sending of a http request
// a custom http client has been made to add a request timeout of 10 seconds
func (h *httpCachedHandler) Send(req repository.HTTPRequest) (interface{}, error) {
	requestHash := h.getHash(req.(*request).innerRequest)
	if response, err := h.getCache(requestHash); err == nil {
		h.logger.Debug("Http - %s - HTTP request retrieved from cache(%s): %+v", req.GetMethod(), requestHash, req.GetPath())
		return response, nil
	}
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
	if err := h.setCache(requestHash, string(response)); err != nil {
		h.logger.Error(
			"Http - %s - HTTP request Error setting cache(%s) for request: %+v, err: %+v",
			req.GetMethod(),
			requestHash,
			req.GetPath(),
			err,
		)
	}
	return string(response), nil
}

// NewRequest returns an initialized struct that can be used to make a http request
func (*httpCachedHandler) NewRequest() repository.HTTPRequest {
	return &request{
		innerRequest: http.Request{
			Header: make(http.Header),
		},
		timeOut: time.Duration(10),
	}
}
