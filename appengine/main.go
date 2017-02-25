package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/lapingvino/bahaibot/badi"
	"strconv"
	"time"
)

const URL = "https://api.telegram.org/bot" + TOKEN + "/"

type Incoming struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID        int64
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string
			Type      string
		}
		Chat struct {
			ID        int64
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string
			Type      string
		}
		Date     int64
		Text     string
		Entities []struct {
			Type   string
			Offset int64
			Length int64
		}
	}
}

func iksigi(in string) string {
	conversion := []struct {
		from string
		to   string
	}{
		{"cx", "ĉ"},
		{"gx", "ĝ"},
		{"hx", "ĥ"},
		{"jx", "ĵ"},
		{"sx", "ŝ"},
		{"ux", "ŭ"},
	}
	for _, c := range conversion {
		in = strings.Replace(in, c.from, c.to, -1)
	}
	return in
}

func telegram(w http.ResponseWriter, r *http.Request) {
	var mymessage string

	c := appengine.NewContext(r)
	client := &http.Client{
		Transport: &urlfetch.Transport{
			Context: c,
		},
	}
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%v", err)
	}
	r.Body.Close()
	var Output Incoming
	err = json.Unmarshal(request, &Output)
	if err != nil {
		log.Errorf(c, "%v", err)
	}
	log.Debugf(c, "%v", Output)
	command := regexp.MustCompile("/[^ @]*").FindString(Output.Message.Text)
	text := regexp.MustCompile("^/[^ ]* ").ReplaceAllString(Output.Message.Text, "")
	switch command {
	case "/start":
		mymessage = "Alláh-u-Abhá! I am the Bahá'í bot. Use the /badi command to get the current time and date according to the Bahá'í calendar."
	case "/echo":
		mymessage = regexp.MustCompile(`(["\\])`).ReplaceAllString(text, `\$1`)
	case "/badi":
		options := strings.Split(text, " ")
		var Localconf badi.Badi
		k := datastore.NewKey(c, "location", "", Output.Message.Chat.ID, nil)
		if err := datastore.Get(c, k, &Localconf); err != nil && len(options) < 3 {
			mymessage = "First configure your location: /badi Timezone/Code Latitude Longitude\n"
			Localconf = badi.Badi{Time: time.Now(), Timezone: "Asia/Tehran", Latitude: 35.715298, Longitude: 51.404343}
		} else {
			Localconf.Time = time.Now()
			if len(options) >= 3 {
				Localconf.Timezone = options[0]
				Localconf.Latitude, err = strconv.ParseFloat(options[1], 64)
				if err != nil {
					mymessage = "Latitude conversion failed"
				}
				Localconf.Longitude, err = strconv.ParseFloat(options[2], 64)
				if err != nil {
					mymessage = "Longitude conversion failed"
				}
			}
		}

		if _, err := datastore.Put(c, k, &Localconf); err != nil {
			mymessage = "Saving location configuration failed\n"
		}
		mymessage += badi.Default(Localconf)
	}
	client.Post(URL+"sendMessage", "application/json", strings.NewReader(fmt.Sprintf("{\"chat_id\": %v, \"text\": \"%v\"}", Output.Message.Chat.ID, mymessage)))
}

func api(w http.ResponseWriter, r *http.Request) {
	options := strings.Split(r.URL.Path, "/")
	selected := ""
	length := len(options)
	pos := 0
	for i := range options {
		if options[i] == "api" && i+1 < length {
			selected = options[i+1]
			pos = i + 1
			break
		}
	}
	switch selected {
	case "badi":
		b := badi.Badi{}
		if length-pos >= 4 {
			lat, err := strconv.ParseFloat(options[pos+3], 64)
			if err != nil {
				fmt.Fprintf(w, "Lat parse error: %v", err)
			}
			long, err := strconv.ParseFloat(options[pos+4], 64)
			if err != nil {
				fmt.Fprintf(w, "Long parse error: %v", err)
			}
			b = badi.Badi{
				Time:      time.Now(),
				Timezone:  options[pos+1] + "/" + options[pos+2],
				Latitude:  lat,
				Longitude: long,
			}
		}
		fmt.Fprintln(w, badi.Default(b))
	}
}

func init() {
	http.HandleFunc("/"+SECRETLINK, telegram)
	http.HandleFunc("/api/", api)
}
