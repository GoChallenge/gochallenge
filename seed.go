package main

import (
	"time"

	"github.com/gochallenge/gochallenge/model"
)

func seedChallenges(cs model.Challenges) {
	c0 := &model.Challenge{
		Name:   "Go Challenge 1 - Drum Machine",
		URL:    "http://golang-challenge.com/go-challenge1/",
		Import: "gc.falsum.me/code/challenge-001",
		Git:    "https://github.com/morhekil/gc-1-drum_machine.git",
		Status: model.Open,
		Author: model.Author{Name: "Matt Aimonetti"},
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(c0)

	c1 := &model.Challenge{
		Name:   "Go Challenge 2 - NaCl Crypto",
		URL:    "http://golang-challenge.com/go-challenge2/",
		Import: "gc.falsum.me/code/challenge-002",
		Git:    "https://github.com/morhekil/gc-2-nacl.git",
		Status: model.Open,
		Author: model.Author{Name: "Guillaume J. Charmes"},
		Start:  time.Date(2015, 4, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 4, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(c1)
}
