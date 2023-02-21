package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddress string `mapstructure:"LISTEN_ADDRESS"`
	LogLevel      string `mapstructure:"LOG_LEVEL"`
	JSONOutput    bool   `mapstructure:"JSON"`
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
