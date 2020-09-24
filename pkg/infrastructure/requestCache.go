package infrastructure

import (
	"crypto/md5" //nolint: gosec
	"encoding/json"
	"fmt"

	"github.com/Yapo/goutils"
	"github.com/anevsky/cachego/memory"
)

// RequestCache holds the cache itself and the variables that control
// its behaviour
type RequestCache struct {
	enabled  bool
	cache    memory.CACHE
	cacheTTL int
}

// getHash will print the data interface as string with all its fields
// and then will md5 it to generate a hash usable as a unique key for
// that data interface
func (rc *RequestCache) getHash(data interface{}) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%v", data))) //nolint: gosec
	return fmt.Sprintf("%x", sum)
}

// GetCache will get a hash for the input data and then will look in memory
// for the hash key to retrieve the data, if found it will unmarshal the data
// into a *gouitls.Response that will be return as the response
func (rc *RequestCache) GetCache(input interface{}) (*goutils.Response, error) {
	var response goutils.Response
	if rc.enabled {
		hash := rc.getHash(input)
		stringResponse, err := rc.cache.Get(hash)
		if err == nil {
			err = json.Unmarshal([]byte(stringResponse.(string)), &response)
		}
		return &response, err
	}
	return &response, fmt.Errorf("cache disabled")
}

// SetCache will get a hash for the input data and then will store in memory
// a json string representation of a goutils.Response associated to the hash key
func (rc *RequestCache) SetCache(input interface{}, response *goutils.Response) error {
	if rc.enabled {
		hash := rc.getHash(input)
		if ok, _ := rc.cache.HasKey(hash); ok {
			return fmt.Errorf("cache already set and still valid")
		}
		stringResponse, err := json.Marshal(response)
		if err != nil {
			return err
		}

		if err := rc.cache.SetString(hash, string(stringResponse)); err != nil {
			return err
		}
		return rc.cache.SetTTL(hash, rc.cacheTTL)
	}
	return fmt.Errorf("cache disabled")
}

// NewRequestCacheHandler will create a new request cache handler that will hold
// data in memory for ttl time in miliseconds
func NewRequestCacheHandler(ttl int) *RequestCache {
	return &RequestCache{
		cache:    memory.Alloc(),
		cacheTTL: ttl,
		enabled:  true,
	}
}
