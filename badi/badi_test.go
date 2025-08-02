package badi

import (
	"strings"
	"testing"
	"time"
)

func TestKnownDays(t *testing.T) {
	testdata := []struct {
		Input     string
		Location  string
		OutputHas []string
		OutputNot []string
	}{
		{"19-3-2017_18:00", "Europe/Amsterdam/52.0882573/5.6173006",
			[]string{"`Alá'", "19"},
			[]string{"Bahá", "Ayyám-i-Há"},
		},
		{"19-3-2017_17:00", "Europe/Amsterdam/52.0882573/5.6173006",
			[]string{"`Alá'", "19"},
			[]string{"Bahá", "Ayyám-i-Há"},
		},
		{"19-3-2017_18:55", "Europe/Amsterdam/52.0882573/5.6173006",
			[]string{"Bahá", "1"},
			[]string{"`Alá'", "Ayyám-i-Há", "19"},
		},
		{"19-3-2017_18:00", "Europe/Amsterdam/52.0882573/5.6173006",
			[]string{"`Alá'", "19"},
			[]string{"Bahá", "Ayyám-i-Há"},
		},
		{"27-2-2017_23:05", "Europe/Amsterdam/52.0882573/5.6173006",
			[]string{"4 Ayyám-i-Há"},
			[]string{"5 Ayyám-i-Há", "`Alá'", "Mulk"},
		},
		{"28-2-2017_00:05", "Europe/Amsterdam/52.0882573/5.6173006",
			[]string{"4 Ayyám-i-Há"},
			[]string{"5 Ayyám-i-Há", "`Alá'", "Mulk"},
		},
	}

	for _, elem := range testdata {
		for _, can := range elem.OutputHas {
			if !strings.Contains(Convert(elem.Input, strings.Split(elem.Location, "/")), can) {
				t.Errorf("Doesn't contain required element %v:\n%#v\n", can, elem)
			}
		}
		for _, cannot := range elem.OutputNot {
			if strings.Contains(Convert(elem.Input, strings.Split(elem.Location, "/")), cannot) {
				t.Errorf("Contains forbidden element %v:\n%#v\n", cannot, elem)
			}
		}
	}
}

func TestYearDay(t *testing.T) {
	beforeNawruz := time.Date(1989, time.January, 25, 0, 0, 0, 0, time.UTC)
	afterNawruz := time.Date(1989, time.November, 23, 0, 0, 0, 0, time.UTC)

	// Create Badi struct with Tehran coordinates and timezone
	tehranBefore := Badi{
		Time:      beforeNawruz,
		Timezone:  "Asia/Tehran",
		Latitude:  35.6892,
		Longitude: 51.3890,
	}
	tehranAfter := Badi{
		Time:      afterNawruz,
		Timezone:  "Asia/Tehran",
		Latitude:  35.6892,
		Longitude: 51.3890,
	}

	bn := tehranBefore.DaysInYear()
	an := tehranAfter.DaysInYear()
	if bn < 365 || bn > 366 || an < 365 || an > 366 {
		t.Errorf("Badí DaysInYear() fails: bn: %v, an: %v", bn, an)
	}
}
