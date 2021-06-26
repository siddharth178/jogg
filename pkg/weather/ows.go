package weather

import (
	"encoding/json"
	stderrors "errors"
	"io"

	"github.com/pkg/errors"
)

/*

Sample API call and its output

Location: Pune, 18.51957, 73.85535

curl -v "https://api.openweathermap.org/data/2.5/onecall/timemachine?lat=18.51957&lon=73.85535&dt=1622956069&appid=5d01143eb3cc9a0a92ca918fbe52efb0" | jq .
{
  "lat": 18.5196,
  "lon": 73.8554,
  "timezone": "Asia/Kolkata",
  "timezone_offset": 19800,
  "current": {
    "dt": 1622956069,
    "sunrise": 1622939244,
    "sunset": 1622986762,
    "temp": 299.04,
    "feels_like": 300.49,
    "pressure": 1010,
    "humidity": 67,
    "dew_point": 292.45,
    "uvi": 14.35,
    "clouds": 99,
    "wind_speed": 2.73,
    "wind_deg": 275,
    "wind_gust": 3.45,
    "weather": [
      {
        "id": 804,
        "main": "Clouds",
        "description": "overcast clouds",
        "icon": "04d"
      }
    ]
  },
  "hourly": []
}

*/

var (
	ErrWeatherNotFound = stderrors.New("weather not found in response")
)

type OwsTimemachineCurrentWeather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type OwsTimemachineCurrent struct {
	Weather []OwsTimemachineCurrentWeather `json:"weather"`
}

type OwsHistoricalWeather struct {
	Current OwsTimemachineCurrent `json:"current"`
}

func parseWeatherResponse(r io.Reader) (*OwsTimemachineCurrentWeather, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	w := OwsHistoricalWeather{}
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, errors.WithStack(err)
	}

	if len(w.Current.Weather) == 0 {
		return nil, errors.WithStack(ErrWeatherNotFound)
	}

	return &w.Current.Weather[0], nil
}
