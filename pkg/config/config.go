package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddress string `mapstructure:"LISTEN_ADDRESS" default:":9000"`
	LogLevel      string `mapstructure:"LOG_LEVEL" default:"info"`
	JSONOutput    bool   `mapstructure:"JSON" default:"false"`
	TendermintRPC string `mapstructure:"TENDERMINT_RPC"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configStruct Config

	err = viper.Unmarshal(&configStruct)
	if err != nil {
		return nil, err
	}

	return &configStruct, nil
}
