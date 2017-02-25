package badi

import (
	"github.com/keep94/sunrise"
	"github.com/lapingvino/bahaibot/nawruz"
	"strconv"
	"time"
)

type Badi struct {
	Time      time.Time
	Timezone  string
	Latitude  float64
	Longitude float64
}

var TEHRAN, _ = time.LoadLocation("Asia/Tehran")

var MONTHS = []string{
	"Ayyám-i-Há",
	"Bahá",
	"Jalál",
	"Jamál",
	"`Azamat",
	"Núr",
	"Rahmat",
	"Kalimát",
	"Kamál",
	"Asmá'",
	"`Izzat",
	"Mashíyyat",
	"`Ilm",
	"Qudrat",
	"Qawl",
	"Masá'il",
	"Sharaf",
	"Sultán",
	"Mulk",
	"`Alá'",
}

func Default(s Badi) string {
	var evening string
	if s.Time.After(s.Sunset()) {
		evening = "\U0001F319"
	} else {
		evening = "\u2600"
	}
	location, err := time.LoadLocation(s.Timezone)
	if err != nil {
		return "Location cannot be set, correct timezone"
	}
	s.Time = s.Time.In(location)
	return s.Time.Format("15:04") + " " + evening + "\n" +
		strconv.Itoa(s.Day()) + " " + MONTHS[s.Month()] + " " + strconv.Itoa(s.Year()) +
		//		" (" + strconv.Itoa(s.YearDay()) + ")" +
		"\n\U0001F305 " + s.Sunrise().Format("15:04") +
		" \U0001F3DC " + s.Sunnoon().Format("15:04") +
		" \U0001F307 " + s.Sunset().Format("15:04")
}

func (s Badi) Nawruz() time.Time {
	var r sunrise.Sunrise
	md := nawruz.Marchday[s.Time.Year()]
	nr := time.Date(s.Time.Year(), time.March, md, 0, 0, 0, 0, TEHRAN)
	r.Around(35.696111, 51.423056, nr)
	return r.Sunset()
}

func (s Badi) Year() int {
	y := s.Time.Year() - 1844
	if s.Time.After(s.Nawruz()) {
		y += 1
	}
	return y
}

func (s Badi) Month() int {
	if s.YearDay() <= 19*18 {
		return s.YearDay() / 19 // First 18 months
	}
	if s.Nawruz().Sub(s.Time) <= 19*24*time.Hour {
		return 19 // `Alá'
	} else {
		return 0 // Ayyám-i-Há
	}
}

func (s Badi) Day() int {
	if s.YearDay() <= 19*18 {
		yd := s.YearDay() % 19 // First 18 months
		if yd == 0 {
			return 19
		} else {
			return yd
		}
	}
	if s.Nawruz().Sub(s.Time) <= 19*24*time.Hour {
		return 19 - int(s.Nawruz().Sub(s.Time).Hours()/24) // `Alá'
	} else {
		return s.YearDay() - 19*18 // Ayyám-i-Há
	}
}

func (s Badi) YearDay() int {
	yd := s.Time.YearDay() - s.Nawruz().YearDay()
	if yd < 1 {
		yd = time.Date(s.Time.Year(), time.December, 31, 0, 0, 0, 0, TEHRAN).YearDay() + yd
	}
	if s.Time.After(s.Sunset()) {
		yd++
	}
	return yd
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
