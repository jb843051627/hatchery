package model

import "time"

const (
	SensorTypeTemperature = "temperature"
	SensorTypeHumidity     = "humidity"
	SensorTypeOxygen      = "oxygen"
	SensorTypeRotation    = "rotation"
)

type SensorReading struct {
	ID          int64     `json:"id"`
	IncubatorID int64     `json:"incubator_id"`
	SensorType  string    `json:"sensor_type"`
	Value       float64   `json:"value"`
	RecordedAt  time.Time `json:"recorded_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func ValidSensorType(s string) bool {
	switch s {
	case SensorTypeTemperature, SensorTypeHumidity, SensorTypeOxygen, SensorTypeRotation:
		return true
	}
	return false
}
