package weather

import (
	"context"
	"time"
)

type MockWeatherService struct {
	Loc string
	Ts  time.Time

	Response *Weather
	Err      error
}

func (m *MockWeatherService) GetWeather(_ context.Context, loc string, ts time.Time) (*Weather, error) {
	m.Loc = loc
	m.Ts = ts

	if m.Err != nil {
		return nil, m.Err
	}
	return m.Response, nil
}
