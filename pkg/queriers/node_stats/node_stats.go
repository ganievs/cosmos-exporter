package node_stats

import (
	"cosmos-exporter/pkg/constants"
	"cosmos-exporter/pkg/query_info"
	"cosmos-exporter/pkg/tendermint"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type NodeStatsQuerier struct {
	TendermintRPC *tendermint.TendermintRPC
	Logger        zerolog.Logger
}

func NewNodeStatsQuerier(logger *zerolog.Logger, tendermintRPC *tendermint.TendermintRPC) *NodeStatsQuerier {
	return &NodeStatsQuerier{
		Logger:        logger.With().Str("component", "tendermint_rpc").Logger(),
		TendermintRPC: tendermintRPC,
	}
}

func (n *NodeStatsQuerier) Name() string {
	return "gaia-stats-querier"
}

func (n *NodeStatsQuerier) Get() ([]prometheus.Collector, []query_info.QueryInfo) {
	queryInfo := query_info.QueryInfo{
		Module:  "tendermint",
		Action:  "node_status",
		Success: false,
	}

	status, err := n.TendermintRPC.GetStatus()
	if err != nil {
		n.Logger.Error().Err(err).Msg("Could not fetch node status")
		return []prometheus.Collector{}, []query_info.QueryInfo{queryInfo}
	}

	timeSinceLatestBlockGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: constants.MetricsPrefix + "time_since_latest_block",
			Help: "Time since latest block, in seconds",
		},
		[]string{},
	)

	latestBlockHeightGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: constants.MetricsPrefix + "latest_block_height",
			Help: "Latest block height",
		},
		[]string{},
	)

	timeSinceLatestBlockGauge.
		With(prometheus.Labels{}).
		Set(time.Since(status.SyncInfo.LatestBlockTime).Seconds())

	latestBlockHeightGauge.
		With(prometheus.Labels{}).
		Set(float64(status.SyncInfo.LatestBlockHeight))

	queryInfo.Success = true

	return []prometheus.Collector{
		timeSinceLatestBlockGauge,
		latestBlockHeightGauge,
	}, []query_info.QueryInfo{queryInfo}
}
