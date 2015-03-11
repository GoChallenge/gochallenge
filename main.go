package main

import (
	"log"
	"net/http"

	"github.com/morhekil/gochallenge/api"
	"github.com/morhekil/gochallenge/mock"
	"github.com/morhekil/gochallenge/model"
)

func main() {
	cs := mock.Challenges{}
	cs.Add(model.Challenge{
		ID:   mock.CurrentID,
		Name: "The Challenge",
	})

	a := api.New(api.Config{
		Challenges: cs,
	})
	log.Fatal(http.ListenAndServe(":8081", a))
}
