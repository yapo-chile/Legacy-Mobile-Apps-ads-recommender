package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/handlers"
)

func main() { //nolint: funlen
	var shutdownSequence = infrastructure.NewShutdownSequence()
	var conf infrastructure.Config
	fmt.Printf("Etag:%d\n", conf.InBrowserCacheConf.InitEtag())
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)

	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	} else {
		fmt.Printf("Config: \n%+v\n", conf)
	}

	fmt.Printf("Setting up Prometheus\n")

	prometheus := infrastructure.MakePrometheusExporter(
		conf.PrometheusConf.Port,
		conf.PrometheusConf.Enabled,
	)

	fmt.Printf("Setting up logger\n")

	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf,
		prometheus.NewEventsCollector(
			"pro-carousel_service_events_total",
			"events tracker counter for pro-carousel service",
		),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(2) //nolint: gomnd
	}

	shutdownSequence.Push(prometheus)
	logger.Info("Initializing resources")

	// HealthHandler
	var healthHandler handlers.HealthHandler

	useBrowserCache := infrastructure.InBrowserCache{
		MaxAge:  conf.InBrowserCacheConf.MaxAge,
		Etag:    conf.InBrowserCacheConf.Etag,
		Enabled: conf.InBrowserCacheConf.Enabled,
	}
	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:         logger,
		Cors:           conf.CorsConf,
		InBrowserCache: useBrowserCache,
		WrapperFuncs:   []infrastructure.WrapperFunc{prometheus.TrackHandlerFunc},
		WithProfiling:  conf.Runtime.Profiling,
		Routes: infrastructure.Routes{
			{
				// This is the base path, all routes will start with this prefix
				Prefix: "",
				Groups: []infrastructure.Route{
					{
						Name:         "Check service health",
						Method:       "GET",
						Pattern:      "/healthcheck",
						Handler:      &healthHandler,
						RequestCache: "10s",
					},
				},
			},
		},
	}

	router := maker.NewRouter()

	server := infrastructure.NewHTTPServer(
		fmt.Sprintf("%s:%d", conf.Runtime.Host, conf.Runtime.Port),
		router,
		logger,
	)
	shutdownSequence.Push(server)
	logger.Info("Starting request serving")

	go server.ListenAndServe()
	shutdownSequence.Wait()
	logger.Info("Server exited normally")
}
