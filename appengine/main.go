package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	//	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	//	"unicode"

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
		mymessage = "Alláh-u-Abhá! I am the Bahá'í bot. Please excuse me while I am not complete yet"
	case "/echo":
		mymessage = regexp.MustCompile(`(["\\])`).ReplaceAllString(text, `\$1`)
	case "/badi":
		mymessage = badi.Default(badi.Badi{Time: time.Now(), Timezone: "Europe/Amsterdam", Latitude: 52.0882573, Longitude: 5.6173006})
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
