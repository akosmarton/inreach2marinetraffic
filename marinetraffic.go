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
	body := "________________\r\n"
	body += fmt.Sprintf("MMSI=%s\r\n", r.Mmsi)
	body += fmt.Sprintf("LAT=%.5f\r\n", r.Lat)
	body += fmt.Sprintf("LON=%.5f\r\n", r.Lon)
	body += fmt.Sprintf("SPEED=%.1f\r\n", r.Speed)
	body += fmt.Sprintf("COURSE=%.0f\n", r.Course)
	body += fmt.Sprintf("TIMESTAMP=%s\n", r.Timestamp)
	body += "________________\n"

	return body
}
