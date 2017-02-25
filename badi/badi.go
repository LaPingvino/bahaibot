package badi

import (
	"github.com/keep94/sunrise"
	"time"
)

type Badi struct {
	Time      time.Time
	Timezone  string
	Latitude  float64
	Longitude float64
}

func Default(s Badi) string {
	var evening string
	if s.Time.After(s.Sunset()) {
		evening = "🌙"
	} else {
		evening = "🌞"
	}
	location, err := time.LoadLocation(s.Timezone)
	if err != nil {
		return "Location cannot be set, correct timezone"
	}
	s.Time = s.Time.In(location)
	return s.Time.Format("15:04") + " " + evening + "\n" +
		s.Day() + " " + s.Month() + " " + s.Year() +
		"\n\U0001F305 " + s.Sunrise().Format("15:04") +
		" \U0001F3DC " + s.Sunnoon().Format("15:04") +
		" \U0001F307 " + s.Sunset().Format("15:04")
}

func (s Badi) Year() string {
	return "173"
}

func (s Badi) Month() string {
	return "Ayyám-i-Há"
}

func (s Badi) Day() string {
	return "1"
}

func (s Badi) Sunrise() time.Time {
	var r sunrise.Sunrise
	r.Around(s.Latitude, s.Longitude, s.Time)
	return r.Sunrise()
}

func (s Badi) Sunnoon() time.Time {
	return s.Sunrise().Add(s.Sunset().Sub(s.Sunrise()) / 2)
}

func (s Badi) Sunset() time.Time {
	var r sunrise.Sunrise
	r.Around(s.Latitude, s.Longitude, s.Time)
	return r.Sunset()
}
