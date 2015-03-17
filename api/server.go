package api

import (
	"net/http"

	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/api/challenges"
	"github.com/gochallenge/gochallenge/api/submissions"
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
	r.GET("/v1/challenges/:id/submissions",
		submissions.List(cfg.Challenges, cfg.Submissions))
	r.GET("/v1/submissions/:id/download", submissions.Download(cfg.Submissions))
	r.POST("/v1/challenges/:id/submissions",
		submissions.Post(cfg.Challenges, cfg.Submissions))
	r.GET("/v1/auth/github", auth.GithubInit())
	r.GET("/v1/auth/github_verify", auth.GithubVerify())
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
