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
		ID:   mock.CurrentID,
		Name: "The Challenge",
	})

	a := api.New(api.Config{
		Challenges: cs,
	})
	log.Fatal(http.ListenAndServe(":8081", a))
}
