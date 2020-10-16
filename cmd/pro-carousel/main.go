package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/handlers"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/repository"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/usecases"
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
	regions, errorRegions := infrastructure.NewRconf(
		conf.EtcdConf.Host,
		conf.EtcdConf.RegionPath,
		conf.EtcdConf.Prefix,
		logger,
	)
	if errorRegions != nil {
		logger.Error("error loading regions from etcd")
		panic(errorRegions)
	}
	fileTools := infrastructure.NewFileTools(conf.ElasticSearchConf.QueryTemplates, ".tmpl")
	queryTemplates := fileTools.LoadTemplatesFromFolder()
	if len(queryTemplates) < 1 {
		errStr := "query templates not found"
		logger.Error(errStr)
		panic(fmt.Errorf(errStr))
	}
	logger.Info("Loaded templates %+v", queryTemplates)
	// interactor loggers
	getSuggestionsLogger := loggers.MakeGetSuggestionsLogger(logger)

	elasticHandler := infrastructure.NewElasticHandlerHandler(
		conf.ElasticSearchConf.MaxIdleConns,
		conf.ElasticSearchConf.MaxIdleConnsPerHost,
		conf.ElasticSearchConf.MaxConnsPerHost,
		conf.ElasticSearchConf.IdleConnTimeout,
		conf.ElasticSearchConf.BatchSize,
		conf.ElasticSearchConf.SearchTimeout,
		conf.ElasticSearchConf.Host+":"+conf.ElasticSearchConf.Port,
		logger,
	)
	// Repos
	adsRepository := repository.NewAdsRepository(
		elasticHandler,
		regions,
		queryTemplates,
		conf.AdConf.ImageServerURL,
		conf.ElasticSearchConf.Index,
		conf.ElasticSearchConf.SearchResultSize,
		conf.ElasticSearchConf.SearchResultPage,
	)
	// Interactors
	getSuggestions := usecases.GetSuggestions{
		SuggestionsRepo:   adsRepository,
		MinDisplayedAds:   conf.AdConf.MinDisplayedAds,
		RequestedAdsQty:   conf.AdConf.DefaultRequestedAdsQty,
		MaxDisplayedAds:   conf.AdConf.MaxDisplayedAds,
		SuggestionsParams: conf.AdConf.SuggestionsParams,
		Logger:            getSuggestionsLogger,
	}
	// HealthHandler
	var healthHandler handlers.HealthHandler
	getProSuggestionsHandler := handlers.GetProSuggestionsHandler{
		Interactor:          &getSuggestions,
		CurrencySymbol:      conf.AdConf.CurrencySymbol,
		UnitOfAccountSymbol: conf.AdConf.UnitOfAccountSymbol,
	}

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
					{
						Name:         "Get suggestions for a specific ad",
						Method:       "GET",
						Pattern:      "/suggestions/{listID:.*}/pro",
						Handler:      &getProSuggestionsHandler,
						RequestCache: "10s",
					}},
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
