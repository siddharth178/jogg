package models

import "time"

type Activity struct {
	Id                 int       `json:"id"`
	UserId             int       `json:"user_id"`
	Ts                 time.Time `json:"ts"`
	Loc                string    `json:"loc"`
	Distance           int       `json:"distance"`
	Seconds            int       `json:"seconds"`
	WeatherCondition   string    `json:"weather_condition"`
	WeatherDescription string    `json:"weather_description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
