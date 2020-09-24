package infrastructure

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// InBrowserCache implement BrowserCache and have the vars
type InBrowserCache struct {
	// MaxAge is used to know how much time the response is valid at
	// browser level
	MaxAge time.Duration
	// Etag contains the identifier of current running version
	Etag int64
	// Enable allows use or ignore the feature
	Enabled bool
}

// Validate set the cache headers and validate if request has changed.
func (ibc *InBrowserCache) Validate(w http.ResponseWriter, r *http.Request) bool {
	if ibc.Enabled {
		key := strconv.FormatInt(ibc.Etag, 10)
		seconds := fmt.Sprintf("%.0f", ibc.MaxAge.Seconds())
		e := `"` + key + `"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age="+seconds)
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				return true
			}
		}
	}
	return false
}

// NewBrowserCache setup the vars of BrowserCache and return the pointer
func NewBrowserCache(enabled bool, etag int64, defaultAge time.Duration, maxAge time.Duration) *InBrowserCache {
	bCache := &InBrowserCache{
		Enabled: enabled,
		Etag:    etag,
		MaxAge:  defaultAge,
	}
	if maxAge > 0 {
		bCache.MaxAge = maxAge
	}
	return bCache
}
