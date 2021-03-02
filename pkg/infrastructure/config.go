package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// LoggerConf holds configuration for logging
// LogLevel definition:
//   0 - Debug
//   1 - Info
//   2 - Warning
//   3 - Error
//   4 - Critic
type LoggerConf struct {
	SyslogIdentity string `env:"SYSLOG_IDENTITY"`
	SyslogEnabled  bool   `env:"SYSLOG_ENABLED" envDefault:"false"`
	StdlogEnabled  bool   `env:"STDLOG_ENABLED" envDefault:"true"`
	LogLevel       int    `env:"LOG_LEVEL" envDefault:"2"`
}

// PrometheusConf holds configuration to report to Prometheus
type PrometheusConf struct {
	Port    string `env:"PORT" envDefault:"8877"`
	Enabled bool   `env:"ENABLED" envDefault:"false"`
}

// RuntimeConfig config to start the app
type RuntimeConfig struct {
	Host      string `env:"HOST" envDefault:"0.0.0.0"`
	Port      int    `env:"PORT" envDefault:"8080"`
	Profiling bool   `env:"PROFILING" envDefault:"true"`
}

// CircuitBreakerConf holds all configurations for circuit breaker
type CircuitBreakerConf struct {
	Name               string  `env:"NAME" envDefault:"HTTP_SEND"`
	ConsecutiveFailure uint32  `env:"CONSECUTIVE_FAILURE" envDefault:"10"`
	FailureRatio       float64 `env:"FAILURE_RATIO" envDefault:"0.5"`
	Timeout            int     `env:"TIMEOUT" envDefault:"30"`
	Interval           int     `env:"INTERVAL" envDefault:"30"`
}

// ProCarouselClientConf holds configuration regarding to our http client (pro-carousel itself in this case)
type ProCarouselClientConf struct {
	TimeOut            int    `env:"TIMEOUT" envDefault:"30"`
	GetHealthcheckPath string `env:"HEALTH_PATH" envDefault:"/get/healthcheck"`
}

// CorsConf holds cors headers
type CorsConf struct {
	Enabled bool   `env:"ENABLED" envDefault:"false"`
	Origin  string `env:"ORIGIN" envDefault:"*"`
	Methods string `env:"METHODS" envDefault:"GET, OPTIONS"`
	Headers string `env:"HEADERS" envDefault:"Accept,Content-Type,Content-Length,If-None-Match,Accept-Encoding,User-Agent"`
}

// EtcdConf configure how to read configuration from remote Etcd service
type EtcdConf struct {
	Host       string `env:"HOST" envDefault:"http://etcd-server.yapo.cl"`
	LastUpdate string `env:"LAST_UPDATE" envDefault:"/last_update"`
	Prefix     string `env:"PREFIX" envDefault:"/v2/keys"`
	RegionPath string `env:"REGION_PATH" envDefault:"/public/location/regions.json"`
	Categories string `env:"CATEGORIES" envDefault:"/public/categories.json"`
}

// AdConf configure how to get ads and how to fill some fields
type AdConf struct {
	ImageServerURL         string                              `env:"IMAGE_SERVER_URL" envDefault:"https://img.yapo.cl/%s/%s/%s.jpg"` //nolint:lll
	CurrencySymbol         string                              `env:"CURRENCY_SYMBOL" envDefault:"$"`
	UnitOfAccountSymbol    string                              `env:"UNIT_OF_ACCOUNT_SYMBOL" envDefault:"UF"`
	MinDisplayedAds        int                                 `env:"MIN_DISPLAYED_ADS" envDefault:"2"`
	MaxDisplayedAds        int                                 `env:"MAX_DISPLAYED_ADS" envDefault:"10"`
	DefaultRequestedAdsQty int                                 `env:"DEFAULT_DISPLAYED_ADS_QTY" envDefault:"10"`
	SuggestionsParams      map[string]map[string][]interface{} `env:"SUGGESTIONS_PARAMS"`
	ContactPath            string                              `env:"CONTACT_PATH" envDefault:"http://ad-contact/contact/phones"` //nolint:lll
}

// ResourcesConf resources path settings
type ResourcesConf struct {
	SuggestionsParams string `env:"SUGGESTIONS_PARAMS" envDefault:"resources/suggestion_params.json"`
}

// ElasticSearchConf configuration for the elastic search client
type ElasticSearchConf struct {
	Index               string        `env:"INDEX_ALIAS" envDefault:"ads_dev09"`
	Host                string        `env:"HOST" envDefault:"http://elastic"`
	Port                string        `env:"PORT" envDefault:"9200"`
	MaxIdleConns        int           `env:"MAX_IDLE_CONNECTIONS" envDefault:"10"`
	MaxIdleConnsPerHost int           `env:"MAX_IDLE_CONNECTIONS_PER_HOST" envDefault:"10"`
	MaxConnsPerHost     int           `env:"MAX_CONNECTIONS_PER_HOST" envDefault:"10"`
	IdleConnTimeout     int           `env:"IDLE_CONNECTIONS_TIMEOUT" envDefault:"3"`
	BatchSize           int           `env:"BATCH_SIZE" envDefault:"10000"`
	SearchResultSize    int           `env:"SEARCH_RESULT_SIZE" envDefault:"10"`
	SearchResultPage    int           `env:"SEARCH_RESULT_PAGE" envDefault:"0"`
	SearchTimeout       time.Duration `env:"SEARCH_TIMEOUT" envDefault:"3s"`
	QueryTemplates      string        `env:"QUERY_TEMPLATES" envDefault:"resources/queries/"`
}

// GetHeaders return map of cors used
func (cc CorsConf) GetHeaders() map[string]string {
	if !cc.Enabled {
		return map[string]string{}
	}

	return map[string]string{
		"Origin":  cc.Origin,
		"Methods": cc.Methods,
		"Headers": cc.Headers,
	}
}

// IndicatorsConf defines the configuration needed to communicate with indicators api
type IndicatorsConf struct {
	UFPath   string `env:"UF_PATH" envDefault:"https://mindicador.cl/api/uf/"`
	CacheTTL int    `env:"CACHE_TTL" envDefault:"600000"` // time in milliseconds
}

// InBrowserCacheConf Used to handle browser cache
type InBrowserCacheConf struct {
	Enabled bool `env:"ENABLED" envDefault:"false"`
	// Cache max age in secs(use browser cache)
	MaxAge time.Duration `env:"MAX_AGE" envDefault:"720h"`
	Etag   int64
}

// InitEtag use current epoc to config etag
func (chc *InBrowserCacheConf) InitEtag() int64 {
	chc.Etag = time.Now().Unix()
	return chc.Etag
}

// Config holds all configuration for the service
type Config struct {
	PrometheusConf        PrometheusConf        `env:"PROMETHEUS_"`
	LoggerConf            LoggerConf            `env:"LOGGER_"`
	Runtime               RuntimeConfig         `env:"APP_"`
	CircuitBreakerConf    CircuitBreakerConf    `env:"CIRCUIT_BREAKER_"`
	ProCarouselClientConf ProCarouselClientConf `env:"PRO_CAROUSEL_"`
	CorsConf              CorsConf              `env:"CORS_"`
	InBrowserCacheConf    InBrowserCacheConf    `env:"BROWSER_CACHE_"`
	ElasticSearchConf     ElasticSearchConf     `env:"ELASTIC_"`
	EtcdConf              EtcdConf              `env:"ETCD_"`
	AdConf                AdConf                `env:"AD_"`
	ResourcesConf         ResourcesConf         `env:"RESOURCES_"`
	IndicatorsConf        IndicatorsConf        `env:"INDICATORS_"`
}

// LoadFromEnv loads the config data from the environment variables
func LoadFromEnv(data interface{}) {
	load(reflect.ValueOf(data), "", "")
}

// valueFromEnv lookup the best value for a variable on the environment
func valueFromEnv(envTag, envDefault string) string {
	// Maybe it's a secret and <envTag>_FILE points to a file with the value
	// https://rancher.com/docs/rancher/v1.6/en/cattle/secrets/#docker-hub-images
	if fileName, ok := os.LookupEnv(fmt.Sprintf("%s_FILE", envTag)); ok {
		// filepath.Clean() will clean the input path and remove some unnecessary things
		// like multiple separators doble "." and others
		// if for some reason you are having troubles reaching your file, check the
		// output of the Clean function and test if its what you expect
		// you can find more info here: https://golang.org/pkg/path/filepath/#Clean
		b, err := ioutil.ReadFile(filepath.Clean(fileName))
		if err == nil {
			return string(b)
		}

		fmt.Print(err)
	}
	// The value might be set directly on the environment
	if value, ok := os.LookupEnv(envTag); ok {
		return value
	}
	// Nothing to do, return the default
	return envDefault
}

// load the variable defined in the envTag into Value
func load(conf reflect.Value, envTag, envDefault string) { //nolint: gocyclo, gocognit
	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// Only attempt to set writeable variables
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value := valueFromEnv(envTag, envDefault)
			// Print message if config is missing
			if envTag != "" && value == "" && !strings.HasSuffix(envTag, "_") {
				fmt.Printf("Config for %s missing\n", envTag)
			}

			switch reflectedConf.Interface().(type) {
			case int:
				if value, err := strconv.ParseInt(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(int(value)))
				}
			case int64:
				if value, err := strconv.ParseInt(value, 10, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case uint32:
				if value, err := strconv.ParseUint(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(uint32(value)))
				}
			case float64:
				if value, err := strconv.ParseFloat(value, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case []string:
				values := strings.Split(value, ",")
				reflectedConf.Set(reflect.ValueOf(values))
			case string:
				reflectedConf.Set(reflect.ValueOf(value))
			case bool:
				if value, err := strconv.ParseBool(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case time.Time:
				if value, err := time.Parse(time.RFC3339, value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case time.Duration:
				if t, err := time.ParseDuration(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(t))
				}
			}

			if reflectedConf.Kind() == reflect.Struct {
				// Recursively load inner struct fields
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
			}
		}
	}
}
