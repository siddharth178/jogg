package models

import "time"

type Weekly struct {
	Week         time.Time `json:"week"`
	Days         int       `json:"days"`
	AvgDistance  float32   `json:"avg_distance"`
	AvgKMs       float32   `json:"avg_kms"`
	AvgHours     float32   `json:"avg_hours"`
	AvgSpeedKMPH float32   `json:"avg_speed_kmph"`
}
