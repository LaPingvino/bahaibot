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

func Tehran(t time.Time) Badi {
	return Badi{Time: t, Timezone: "Asia/Tehran", Latitude: 35.696111, Longitude: 51.423056}
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
	if s.Month() < 0 || s.Month() > 19 {
		return "Month set to " + strconv.Itoa(s.Month()) + ", aborting"
	}
	s.Time = s.Time.In(location)
	return s.Time.Format("15:04") + " " + evening + "\n" +
		strconv.Itoa(s.Day()) + " " + MONTHS[s.Month()] + " " + strconv.Itoa(s.Year()) +
		//		strconv.Itoa(s.Day()) + " " + strconv.Itoa(s.Month()) + " " + strconv.Itoa(s.Year()) +
		"\n\U0001F305 " + s.Sunrise().Format("15:04") +
		" \U0001F3DC " + s.Sunnoon().Format("15:04") +
		" \U0001F307 " + s.Sunset().Format("15:04") +
		//		"\n(" + strconv.Itoa(s.YearDay()) + " " + s.Time.String() + ")" +
		""
}

func Convert(input string, location []string) string {
	b := Badi{Time: time.Now(), Timezone: "Asia/Tehran", Latitude: 35.715298, Longitude: 51.404343}

	if len(location) >= 4 {
		lat, err := strconv.ParseFloat(location[2], 64)
		if err != nil {
			return "Lat parse error: " + err.Error()
		}
		long, err := strconv.ParseFloat(location[3], 64)
		if err != nil {
			return "Long parse error: " + err.Error()
		}
		b.Timezone = location[0] + "/" + location[1]
		b.Latitude = lat
		b.Longitude = long
		loc, err := time.LoadLocation(b.Timezone)
		if err != nil {
			loc = TEHRAN
		}
		if t, err := time.ParseInLocation("2-1-2006", input, loc); err == nil {
			b.Time = t
		}

		if t, err := time.ParseInLocation("2-1-2006_15:04", input, loc); err == nil {
			b.Time = t
		}
	}
	return Default(b)
}

func (s Badi) Nawruz() time.Time {
	var r sunrise.Sunrise
	nr := nawRuz(s.Time.Year())
	r.Around(s.Latitude, s.Longitude, nr)
	return r.Sunset()
}

func nawRuz(year int) time.Time {
	md := nawruz.Marchday[year]
	nr := time.Date(year, time.March, md-1, 12, 0, 0, 0, TEHRAN)
	return nr
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
		m := s.YearDay() / 19 // First 18 months
		if s.YearDay()%19 != 0 {
			m++
		}
		return m
	}
	if s.DaysInYear()-s.YearDay() < 19 {
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
	if s.DaysInYear()-s.YearDay() < 19 {
		return 19 - (s.DaysInYear() - s.YearDay()) // `Alá'
	} else {
		return s.YearDay() - 19*18 // Ayyám-i-Há
	}
}

func daysInYear(year int) int {
	diy := 365
	if time.Date(year, time.March, 0, 0, 0, 0, 0, TEHRAN).Day() > 28 {
		diy++
	}
	return diy
}

func (s Badi) DaysInYear() int {
	println(daysInYear(s.Time.Year()))
	println(daysInYear(s.Time.Year() - 1))
	println(s.Nawruz().YearDay())
	println(nawRuz(s.Time.Year() - 1).YearDay())
	println(nawRuz(s.Time.Year() + 1).YearDay())
	println(s.Nawruz().YearDay())
	if s.Time.After(s.Nawruz()) {
		return (daysInYear(s.Time.Year()) - s.Nawruz().YearDay()) + nawRuz(s.Time.Year()+1).YearDay()
	} else {
		return (daysInYear(s.Time.Year()-1) - nawRuz(s.Time.Year()-1).YearDay()) + s.Nawruz().YearDay()
	}
}

func (s Badi) YearDay() int {
	var yd int
	if s.Time.After(s.Nawruz()) {
		yd = s.Time.YearDay() - s.Nawruz().YearDay()
		if !s.Time.Before(s.Sunset()) {
			yd++
		}
	} else {
		yd = (daysInYear(s.Time.Year()-1) - nawRuz(s.Time.Year()-1).YearDay()) + s.Time.YearDay() - 1
		if !s.Time.Before(s.Sunset()) {
			yd++
		}
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
