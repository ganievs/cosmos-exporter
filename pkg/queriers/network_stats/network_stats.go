package network_stats

import (
	"cosmos-exporter/pkg/constants"
	"cosmos-exporter/pkg/query_info"
	"cosmos-exporter/pkg/tendermint"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type NetworkStatsQuerier struct {
	TendermintRPC *tendermint.TendermintRPC
	Logger        zerolog.Logger
}

func NewNetworkStatsQuerier(logger *zerolog.Logger, tendermintRPC *tendermint.TendermintRPC) *NetworkStatsQuerier {
	return &NetworkStatsQuerier{
		Logger:        logger.With().Str("component", "tendermint_rpc").Logger(),
		TendermintRPC: tendermintRPC,
	}
}

func (n *NetworkStatsQuerier) Name() string {
	return "gaia-stats-querier"
}

func (n *NetworkStatsQuerier) Get() ([]prometheus.Collector, []query_info.QueryInfo) {
	queryInfo := query_info.QueryInfo{
		Module:  "tendermint",
		Action:  "network_information",
		Success: false,
	}

	netInfo, err := n.TendermintRPC.NetInfo()
	if err != nil {
		n.Logger.Error().Err(err).Msg("Could not fetch network information")
		return []prometheus.Collector{}, []query_info.QueryInfo{queryInfo}
	}

	peersCountGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: constants.MetricsPrefix + "peers_count",
			Help: "Count of peers",
		},
		[]string{},
	)

	peersCountGauge.
		With(prometheus.Labels{}).
		Set(float64(len(netInfo.Peers)))

	queryInfo.Success = true

	return []prometheus.Collector{
		peersCountGauge,
	}, []query_info.QueryInfo{queryInfo}
}
