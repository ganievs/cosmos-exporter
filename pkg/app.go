package pkg

import (
	"cosmos-exporter/pkg/config"
	"cosmos-exporter/pkg/queriers/network_stats"
	"cosmos-exporter/pkg/queriers/node_stats"
	"cosmos-exporter/pkg/query_info"
	"cosmos-exporter/pkg/tendermint"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type App struct {
	Logger   zerolog.Logger
	Queriers []query_info.Querier
}

func NewApp(
	logger *zerolog.Logger,
	config *config.Config,
) *App {

	appLogger := logger.With().Str("component", "app").Logger()
	tendermintRPC := tendermint.NewTendermintRPC(config, logger)

	statsQueriers := []query_info.Querier{
		network_stats.NewNetworkStatsQuerier(logger, tendermintRPC),
		node_stats.NewNodeStatsQuerier(logger, tendermintRPC),
	}

	return &App{
		Logger:   appLogger,
		Queriers: statsQueriers,
	}
}

func (a *App) HandleRequest(w http.ResponseWriter, r *http.Request) {
	requestStart := time.Now()

	sublogger := a.Logger.With().
		Str("request-id", uuid.New().String()).
		Logger()

	var wg sync.WaitGroup
	var mu sync.Mutex
	allResults := map[string][]prometheus.Collector{}
	allQueries := map[string][]query_info.QueryInfo{}

	for _, querier := range a.Queriers {
		wg.Add(1)
		go func(querier query_info.Querier) {
			querierResults, queriesInfo := querier.Get()
			mu.Lock()
			allResults[querier.Name()] = querierResults
			allQueries[querier.Name()] = queriesInfo
			mu.Unlock()
			wg.Done()
		}(querier)
	}

	wg.Wait()

	allResults["query_infos"] = query_info.GetQueryInfoMetrics(allQueries)

	registry := prometheus.NewRegistry()

	for _, querierResults := range allResults {
		for _, result := range querierResults {
			registry.MustRegister(result)
		}
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

	sublogger.Info().
		Str("method", http.MethodGet).
		Str("endpoint", "/metrics").
		Float64("request-time", time.Since(requestStart).Seconds()).
		Msg("Request processed")
}
