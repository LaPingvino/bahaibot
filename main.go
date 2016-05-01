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
	"unicode"
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

type Informoj struct {
	Celvorto string
	Diveno   []rune
	Vicoj    int
}

func kontroli(d Informoj, diveno rune) (montri string, ĝusta bool, kompleta bool) {
	kompleta = true
	for _, c := range d.Celvorto {
		if c == diveno {
			ĝusta = true
		}
		litero := "_"
		for _, s := range d.Diveno {
			if c == s {
				litero = string(c)
			}
		}
		if litero == "_" {
			kompleta = false
		}
		montri += " " + litero
	}
	return montri, ĝusta, kompleta
}

func telegram(w http.ResponseWriter, r *http.Request) {
	var mymessage string
	var Diveno Informoj

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
		mymessage = "Saluton! Bonvenon al la pendumula roboto. Uzu la komandon /diveni por komenci ludon kaj diveni literojn. Por pli facile diveni literojn, ankaŭ eblas simple uzi suprenstreko spaco litero, do ekzemple '/ o'. Bonan ludadon!"
	case "/echo":
		mymessage = regexp.MustCompile(`(["\\])`).ReplaceAllString(text, `\$1`)
	case "/maldeculo":
		mymessage = "Kiu? Mi? Eble vi estas, " + Output.Message.From.FirstName + "... \U0001F60F"
	case "/diveni", "/":
		var diveno rune
		k := datastore.NewKey(c, "Diveno", "", Output.Message.Chat.ID, nil)
		if err := datastore.Get(c, k, &Diveno); err != nil || Diveno.Vicoj < 1 {
			mymessage = "Ni komencas novan ludon, sendu literon por diveni.\n"
			Diveno.Celvorto = elektivorton()
			Diveno.Vicoj = 10
			Diveno.Diveno = []rune{}
		}
		diveno = []rune(text)[0]
		if unicode.IsLetter(diveno) {
			Diveno.Diveno = append(Diveno.Diveno, diveno)
		} else {
			Diveno.Vicoj++
		}
		montri, ĝusta, kompleta := kontroli(Diveno, diveno)
		if text != Diveno.Celvorto {
			mymessage += "Vi divenis '" + string(diveno) + "'\n"
			mymessage += montri + "\n"
			if ĝusta {
				mymessage += "Tiu litero enestas!"
			} else {
				mymessage += "Bedaŭrinde tiu litero ne enestas..."
				Diveno.Vicoj--
			}
			mymessage += "\nVi ĝis nun divenis la literojn " + string(Diveno.Diveno)
		}
		if kompleta || Diveno.Celvorto == text {
			mymessage += "\nVi ĝuste divenis " + Diveno.Celvorto + "!"
			Diveno.Vicoj = -1
		} else {
			switch Diveno.Vicoj {
			case 0:
			case 1:
				mymessage += "\nVi ne plu rajtas erari... Sukceson!"
			default:
				mymessage += fmt.Sprintf("\nVi rajtas ankoraŭ maksimume %v-foje erari", Diveno.Vicoj-1)
			}
		}
		if Diveno.Vicoj == 0 {
			mymessage += "\nVi ne sukcesis diveni...\nLa vorto estis " + Diveno.Celvorto
		}

		if _, err := datastore.Put(c, k, &Diveno); err != nil {
			mymessage += "\nVia nova diveno ne sukcese konserviĝis...\nEraro: " + err.Error() + "\n" + fmt.Sprintf("%#v", Diveno)
		}
	}
	client.Post(URL+"sendMessage", "application/json", strings.NewReader(fmt.Sprintf("{\"chat_id\": %v, \"text\": \"%v\"}", Output.Message.Chat.ID, mymessage)))
}

func init() {
	http.HandleFunc("/"+SECRETLINK, telegram)
}
