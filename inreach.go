package main

import "time"

type InreachMessage struct {
	Description string
	Timestamp   time.Time
	DeviceType  string
	Latitude    float64
	Longitude   float64
	Elevation   float64
	Course      float64
	Velocity    float64
}
