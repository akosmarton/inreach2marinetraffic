package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	mapshareID := os.Getenv("MAPSHARE_ID")
	mapsharePassword := os.Getenv("MAPSHARE_PASSWORD")
	mapshareInterval, _ := strconv.ParseInt(os.Getenv("MAPSHARE_INTERVAL"), 10, 64)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	mmsi := os.Getenv("MMSI")

	if mapshareID == "" {
		log.Fatal("MAPSHARE_ID is empty")
	}
	if smtpHost == "" {
		log.Fatal("SMTP_HOST is empty")
	}
	if smtpPort == "" {
		log.Fatal("SMTP_PORT is empty")
	}
	if smtpUser == "" {
		log.Fatal("SMTP_USER is empty")
	}
	if smtpPassword == "" {
		log.Fatal("SMTP_PASSWORD is empty")
	}
	if emailAddress == "" {
		log.Fatal("EMAIL_ADDRESS is empty")
	}
	if mmsi == "" {
		log.Fatal("MMSI is empty")
	}
	if mapshareInterval == 0 {
		mapshareInterval = 60
	}

	smtpClient := SmtpClient{
		Host:     smtpHost,
		Port:     smtpPort,
		User:     smtpUser,
		Password: smtpPassword,
	}

	d1 := time.Now().UTC()
	d2 := time.Now().UTC()
	for {
		time.Sleep(time.Second * (time.Duration)(mapshareInterval))
		d2 = time.Now().UTC()
		url := fmt.Sprintf("https://share.garmin.com/feed/Share/%s?d1=%s&d2=%s", mapshareID, d1.Format("2006-01-02T15:04:05Z"), d2.Format("2006-01-02T15:04:05Z"))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			continue
		}
		if mapsharePassword != "" {
			req.SetBasicAuth("", mapsharePassword)
		}
		req.Header.Set("cache-control", "no-cache")

		log.Println(req.URL.RequestURI())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Println(resp.StatusCode, resp.Status)
			continue
		}

		k := KML{}
		if err := xml.NewDecoder(resp.Body).Decode(&k); err != nil {
			log.Println(err)
			continue
		}

		log.Printf("%s  %d  %d\n", url, resp.StatusCode, len(k.Placemark))

		im := InreachMessage{}
		for _, pm := range k.Placemark {
			if pm.Timestamp.IsZero() {
				continue
			}

			im.Timestamp = pm.Timestamp.UTC()
			im.Description = pm.Description

			for _, v := range pm.Data {
				switch v.Name {
				case "Latitude":
					fmt.Sscanf(v.Value, "%f", &im.Latitude)
				case "Longitude":
					fmt.Sscanf(v.Value, "%f", &im.Longitude)
				case "Elevation":
					fmt.Sscanf(v.Value, "%f m from MSL", &im.Elevation)
				case "Device Type":
					im.DeviceType = v.Value
				case "Course":
					fmt.Sscanf(v.Value, "%f °", &im.Course)
				case "Velocity":
					fmt.Sscanf(v.Value, "%f km/h", &im.Velocity)
				}
			}

			r := &SelfReport{
				Mmsi:      mmsi,
				Lat:       im.Latitude,
				Lon:       im.Longitude,
				Course:    im.Course,
				Speed:     im.Velocity / 1.852,
				Timestamp: im.Timestamp.UTC(),
			}

			m := &Mail{
				From:    emailAddress,
				To:      MARINETRAFFIC_EMAIL,
				Subject: "self report",
				Body:    r.ToString(),
			}

			if err := smtpClient.Send(m); err != nil {
				log.Println(err)
				continue
			}
			if im.Timestamp.After(d1) {
				d1 = r.Timestamp.Add(time.Second)
			}
			log.Printf("%s", m.Body)
		}
	}
}
