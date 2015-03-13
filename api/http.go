package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/morhekil/mw"
)

// func render(r service.Response, serr *service.Error, w http.ResponseWriter) {
// 	xr := xpgResponse{Data: r}.add(serr)
// 	s, err := json.Marshal(xr)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ps := fmt.Sprintf(`{"d":%s}`, s)
// 	w.Write([]byte(ps))
// }

// func worker(h service.Handler) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		var res service.Response
// 		d, err := h.Parse(r.Body)

// 		if err == nil {
// 			res, err = h.Handle(d)
// 		}
// 		render(res, err, w)
// 	}
// }

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

	return r
}

func (a *app) getChallenge(w http.ResponseWriter,
	r *http.Request, ps httprouter.Params) {
	var (
		c   model.Challenge
		err error
		js  []byte
	)

	id := ps.ByName("id")
	if id == "current" {
		c, err = a.cfg.Challenges.Current()
	} else {
		c, err = a.cfg.Challenges.Find(id)
	}

	if err == nil {
		js, err = json.Marshal(c)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write(js)
	}
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
