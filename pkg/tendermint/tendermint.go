package tendermint

import (
	"context"
	"cosmos-exporter/pkg/config"

	"github.com/rs/zerolog"

	rpcClient "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type TendermintRPC struct {
	Logger       zerolog.Logger
	Client       *rpcClient.HTTP
	BlocksBehind int64
}

func NewTendermintRPC(config *config.Config, logger *zerolog.Logger) *TendermintRPC {
	client, err := rpcClient.New(config.TendermintRPC, "websocket")
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot instantiate Tendermint client")
	}

	return &TendermintRPC{
		Logger:       logger.With().Str("component", "tendermint_rpc").Logger(),
		Client:       client,
		BlocksBehind: 1000,
	}
}

func (t *TendermintRPC) GetStatus() (*coretypes.ResultStatus, error) {
	return t.Client.Status(context.Background())
}

func (t *TendermintRPC) NetInfo() (*coretypes.ResultNetInfo, error) {
	return t.Client.NetInfo(context.Background())
}
