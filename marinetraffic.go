package main

import (
	"fmt"
	"time"
)

type SelfReport struct {
	Mmsi      string
	Lat       float64
	Lon       float64
	Speed     float64
	Course    float64
	Timestamp time.Time
}

const MARINETRAFFIC_EMAIL = "report@marinetraffic.com"

func (r *SelfReport) ToString() string {
	body := "________________\n"
	body += fmt.Sprintf("MMSI=%s\n", r.Mmsi)
	body += fmt.Sprintf("LAT=%.5f\n", r.Lat)
	body += fmt.Sprintf("LON=%.5f\n", r.Lon)
	body += fmt.Sprintf("SPEED=%.1f\n", r.Speed)
	body += fmt.Sprintf("COURSE=%.0f\n", r.Course)
	body += fmt.Sprintf("TIMESTAMP=%s\n", r.Timestamp)
	body += "________________\n"

	return body
}
