package main

import (
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
)

var _ = func() bool {
	rand.Seed(time.Now().UnixNano())
	return true
}()

func elektivorton() string {
	bvortoj, err := ioutil.ReadFile("static/eblaj-vortoj.txt")
	if err != nil {
		return "misfunkcio"
	}
	vortoj := strings.Split(string(bvortoj),"\n")
	elekto := rand.Intn(len(vortoj))
	return vortoj[elekto]
}
