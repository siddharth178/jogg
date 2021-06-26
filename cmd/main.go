package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"gitlab.com/toptal/sidd/jogg/pkg/config"
	"gitlab.com/toptal/sidd/jogg/pkg/datastore"
	"gitlab.com/toptal/sidd/jogg/pkg/logging"
	"gitlab.com/toptal/sidd/jogg/pkg/weather"
	"gitlab.com/toptal/sidd/jogg/pkg/web"
)

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("config file not specified")
	}

	ctx := context.Background()

	cfg, err := config.LoadConfig(args[0])
	if err != nil {
		return err
	}
	logger, err := logging.GetLogger(cfg.LogConfig)
	if err != nil {
		return err
	}
	logger.Infof("running with config: %+v", *cfg)

	ds, err := datastore.NewDS(ctx, cfg.DBConfig, logger)
	if err != nil {
		return err
	}
	logger.Info("db connection done")

	var ws weather.WeatherService
	if cfg.WeatherConfig.Mock {
		logger.Info("using mock weather service")
		ws = &weather.MockWeatherService{Response: &weather.DefaultWeather}
	} else {
		logger.Info("using open weather service")
		w, err := weather.NewWeatherService(logger, cfg.WeatherConfig)
		if err != nil {
			return err
		}
		ws = w
	}

	server := web.NewServer(ctx, logger, cfg.HTTPConfig, ds, ws)
	return server.Start()
}
