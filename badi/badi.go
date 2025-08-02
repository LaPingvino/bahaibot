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
		"\n\U0001F305 " + s.Sunrise().Format("15:04") +
		" \U0001F3DC " + s.Sunnoon().Format("15:04") +
		" \U0001F307 " + s.Sunset().Format("15:04") +
		""
}

func Convert(input string, location []string) string {
	b := Badi{Time: time.Now(), Timezone: "UTC", Latitude: 0, Longitude: 0}

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
			loc = time.UTC
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
	md := nawruz.Marchday[s.Time.Year()]
	// Use sunset time on Nawruz day in Tehran (standard reference)
	nr := time.Date(s.Time.Year(), time.March, md, 18, 0, 0, 0, time.UTC)
	return nr
}

// BadiYearStart returns the actual start of the Badí year (sunset on the day before Nawruz)
func (s Badi) BadiYearStart() time.Time {
	var r sunrise.Sunrise
	md := nawruz.Marchday[s.Time.Year()]
	nawruzDay := time.Date(s.Time.Year(), time.March, md, 12, 0, 0, 0, time.UTC)
	// Get sunset on the day before Nawruz
	dayBefore := nawruzDay.AddDate(0, 0, -1)
	r.Around(s.Latitude, s.Longitude, dayBefore)
	return r.Sunset()
}
func nawRuz(year int) time.Time {
	md := nawruz.Marchday[year]
	nr := time.Date(year, time.March, md, 12, 0, 0, 0, time.UTC)
	return nr
}

func (s Badi) Year() int {
	y := s.Time.Year() - 1844
	if s.Time.After(s.BadiYearStart()) {
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
		return 20 - (s.DaysInYear() - s.YearDay()) // `Alá'
	} else {
		return s.YearDay() - 19*18 + 1 // Ayyám-i-Há
	}
}

func daysInYear(year int) int {
	diy := 365
	if time.Date(year, time.March, 0, 0, 0, 0, 0, time.UTC).Day() > 28 {
		diy++
	}
	return diy
}

func (s Badi) DaysInYear() int {
	if s.Time.After(s.Nawruz()) {
		return (daysInYear(s.Time.Year()) - s.Nawruz().YearDay()) + nawRuz(s.Time.Year()+1).YearDay()
	} else {
		return (daysInYear(s.Time.Year()-1) - nawRuz(s.Time.Year()-1).YearDay()) + s.Nawruz().YearDay()
	}
}

func (s Badi) YearDay() int {
	var yd int
	if s.Time.After(s.BadiYearStart()) {
		// We're in the new Badí year
		yearStart := s.BadiYearStart()
		daysSinceStart := int(s.Time.Sub(yearStart).Hours() / 24)
		yd = daysSinceStart + 1
		// Only increment for sunset if we're not on the year start day
		if s.Time.After(s.Sunset()) && s.Time.UTC().Truncate(24*time.Hour) != yearStart.UTC().Truncate(24*time.Hour) {
			yd++
		}
	} else {
		// We're in the previous Badí year
		yd = (daysInYear(s.Time.Year()-1) - nawRuz(s.Time.Year()-1).YearDay()) + s.Time.YearDay()
		if s.Time.After(s.Sunset()) {
			yd++
		}
	}
	return yd
}

func (s Badi) Sunrise() time.Time {
	t := s.Time
	if s.Time.Hour() < 2 {
		t = s.Time.Add(2 * time.Hour)
	}
	var r sunrise.Sunrise
	r.Around(s.Latitude, s.Longitude, t)
	return r.Sunrise()
}

func (s Badi) Sunnoon() time.Time {
	return s.Sunrise().Add(s.Sunset().Sub(s.Sunrise()) / 2)
}

func (s Badi) Sunset() time.Time {
	t := s.Time
	if s.Time.Hour() < 2 {
		t = s.Time.Add(2 * time.Hour)
	}
	var r sunrise.Sunrise
	r.Around(s.Latitude, s.Longitude, t)
	return r.Sunset()
}
