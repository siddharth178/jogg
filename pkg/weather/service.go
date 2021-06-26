package weather

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/toptal/sidd/jogg/pkg/config"
	"go.uber.org/zap"
)

type Weather struct {
	Condition   string
	Description string
}

var (
	DefaultWeather = Weather{
		Condition:   "NA",
		Description: "NA",
	}
)

type WeatherService interface {
	GetWeather(ctx context.Context, loc string, ts time.Time) (*Weather, error)
}

type OpenWeatherService struct {
	logger *zap.SugaredLogger
	cfg    config.WeatherConfig

	lll LocLatLong
}

func NewWeatherService(logger *zap.SugaredLogger, cfg config.WeatherConfig) (*OpenWeatherService, error) {
	lll, err := buildLocLatLong(cfg.FileName)
	if err != nil {
		return nil, err
	}
	logger.Infof("loaded %d locations", len(lll))

	return &OpenWeatherService{
		logger: logger,
		cfg:    cfg,
		lll:    lll,
	}, nil
}

func (ws *OpenWeatherService) GetWeather(ctx context.Context, loc string, ts time.Time) (*Weather, error) {
	ll, ok := ws.lll[loc]
	if !ok {
		ws.logger.Warnw("unknown location, skipping", "loc", loc)
		return &DefaultWeather, nil
	}

	//
	// Call OpenWeatherAPI
	//
	ws.logger.Infow("calling openweather", "lat", ll.Lat, "long", ll.Long, "ts", ts, "ts.Unix", ts.UTC().Unix())
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall/timemachine?lat=%v&lon=%v&dt=%v&appid=%v", ll.Lat, ll.Long, ts.UTC().Unix(), ws.cfg.APIKey)
	ws.logger.Debugf("req url: %v", url)
	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &DefaultWeather, errors.WithStack(err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		ws.logger.Errorw("OWS API call failed", "err", err)
		return &DefaultWeather, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &DefaultWeather, errors.New("non ok response from OWS API")
	}
	w, err := parseWeatherResponse(resp.Body)
	if err != nil {
		return &DefaultWeather, err
	}

	return &Weather{w.Main, w.Description}, nil
}
