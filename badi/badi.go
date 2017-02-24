package badi

import "time"

func Default() string {
	return Day(time.Now()) + " " + Month(time.Now()) + " " + Year(time.Now())
}

func Year(t time.Time) string {
	return "173"
}

func Month(t time.Time) string {
	return "Ayyám-i-Há"
}

func Day(t time.Time) string {
	return "1"
}
