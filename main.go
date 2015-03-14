package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
)

func main() {
	cs := mock.NewChallenges()
	cs.Add(model.Challenge{
		ID:     mock.CurrentID - 1,
		Name:   "Go Challenge 0 - That Did Not Exist",
		URL:    "http://golang-challenge.com/go-challenge0/",
		Import: "gc.falsum.me/code/challenge-000",
		Git:    "https://github.com/morhekil/gc-1-drum_machine.git",
		Status: model.Closed,
		Author: model.Author{
			Name: "Gordon Freeman",
		},
		Start: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2015, 2, 14, 0, 0, 0, 0, time.UTC),
	})
	cs.Add(model.Challenge{
		ID:     mock.CurrentID,
		Name:   "Go Challenge 1 - Drum Machine",
		URL:    "http://golang-challenge.com/go-challenge1/",
		Import: "gc.falsum.me/code/challenge-001",
		Git:    "https://github.com/morhekil/gc-1-drum_machine.git",
		Status: model.Open,
		Author: model.Author{
			Name: "Matt Aimonetti",
		},
		Start: time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	})

	a := api.New(api.Config{
		Challenges: cs,
	})
	log.Fatal(http.ListenAndServe(":8081", a))
}
