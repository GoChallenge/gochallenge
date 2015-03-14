package api

import (
	"net/http"

	"github.com/gochallenge/gochallenge/api/challenges"
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

func server(cfg Config) *httprouter.Router {
	r := httprouter.New()
	r.GET("/v1/challenges", challenges.List(cfg.Challenges))
	r.GET("/v1/challenges/:id", challenges.Get(cfg.Challenges))
	r.GET("/code/:id", challenges.Get(cfg.Challenges))

	return r
}

// New server created and configured as an instance of martini server
func New(cfg Config) http.Handler {
	hs := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	return alice.New(
		mw.Recover,
		mw.Logger,
		mw.Headers(hs),
	).Then(server(cfg))
}
