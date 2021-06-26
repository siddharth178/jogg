package weather

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_parseWeatherResponse_HappyPath(t *testing.T) {
	s := `{
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
	"hourly": [
	  {
		"dt": 1622937600,
		"temp": 297.8,
		"feels_like": 300.51,
		"pressure": 1008,
		"humidity": 76,
		"dew_point": 293.29,
		"clouds": 97,
		"wind_speed": 1.49,
		"wind_deg": 271,
		"wind_gust": 2.32,
		"weather": [
		  {
			"id": 500,
			"main": "Rain",
			"description": "light rain",
			"icon": "10n"
		  }
		],
		"rain": {
		  "1h": 0.11
		}
	  }
	]
  }`

	w, err := parseWeatherResponse(strings.NewReader(s))
	if assert.NoError(t, err) {
		assert.Equal(t, "Clouds", w.Main)
		assert.Equal(t, "overcast clouds", w.Description)
	}
}

func Test_parseWeatherResponse_SadPath(t *testing.T) {
	s := `{
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
	  "weather": []
	},
	"hourly": []
  }`

	_, err := parseWeatherResponse(strings.NewReader(s))
	assert.Equal(t, ErrWeatherNotFound, errors.Cause(err))
}
