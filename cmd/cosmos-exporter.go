package main

import (
	"cosmos-exporter/pkg"
	"cosmos-exporter/pkg/config"
	"cosmos-exporter/pkg/logger"
	"fmt"
	"net/http"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	fmt.Println(config)
	log := logger.GetLogger(config)
	app := pkg.NewApp(log, config)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		app.HandleRequest(w, r)
	})

	log.Info().Str("addr", config.ListenAddress).Msg("Listening")
	err = http.ListenAndServe(config.ListenAddress, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start application")
	}
}
