package main

import (
	"log"
	"net/http"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
)

func main() {
	cs := mock.Challenges{}
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
	})

	a := api.New(api.Config{
		Challenges: cs,
	})
	log.Fatal(http.ListenAndServe(":8081", a))
}
