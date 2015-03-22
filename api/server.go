package api

import (
	"net/http"
	"strings"

	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/api/challenges"
	"github.com/gochallenge/gochallenge/api/submissions"
	"github.com/gochallenge/gochallenge/api/users"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/morhekil/mw"
)

const assetsPath = "web/assets"
const indexPath = "web/index.html"

// Config of the API setup
type Config struct {
	Challenges  model.Challenges
	Submissions model.Submissions
	Users       model.Users
}

func server(cfg Config) *httprouter.Router {
	r := httprouter.New()
	r.GET("/v1/challenges", challenges.List(cfg.Challenges))
	r.GET("/v1/challenges/:id", challenges.Get(cfg.Challenges))
	r.GET("/v1/challenges/:id/submissions",
		submissions.List(cfg.Challenges, cfg.Submissions))
	r.GET("/v1/submissions/:id/download", submissions.Download(cfg.Submissions))
	r.POST("/v1/challenges/:id/submissions",
		submissions.Post(cfg.Challenges, cfg.Submissions, cfg.Users))
	r.GET("/v1/auth/github", auth.GithubInit())
	r.GET("/v1/auth/github_verify", auth.GithubVerify(cfg.Users))
	r.GET("/v1/users/:id", users.Get(cfg.Users))
	r.GET("/v1/user", users.Me(cfg.Users))
	r.GET("/code/:id", challenges.Get(cfg.Challenges))

	r.ServeFiles("/assets/*filepath", http.Dir(assetsPath))
	r.GET("/", indexHTML)
	return r
}

func indexHTML(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, indexPath)
}

// New server created and configured as an instance of martini server
func New(cfg Config) http.Handler {
	return alice.New(
		mw.Recover,
		mw.Logger,
		headerMiddleware,
	).Then(server(cfg))
}

func headerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/"):
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case strings.HasPrefix(r.URL.Path, "/code"):
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		h.ServeHTTP(w, r)
	})
}
