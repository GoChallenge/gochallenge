package api

import (
	"net/http"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/morhekil/mw"
)

// Config of the API setup
type Config struct {
	Challenges  model.Challenges
	Submissions model.Submissions
}

type app struct {
	cfg Config
}

func (a *app) server() *httprouter.Router {
	r := httprouter.New()
	r.GET("/v1/challenge/:id", a.getChallenge)
	r.GET("/code/:id", a.getChallenge)

	return r
}

// New server created and configured as an instance of martini server
func New(cfg Config) http.Handler {
	hs := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	a := app{cfg}

	return alice.New(
		mw.Recover,
		mw.Logger,
		mw.Headers(hs),
	).Then(a.server())
}
