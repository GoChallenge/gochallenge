package auth

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

const githubOAuthURL = "https://github.com/login/oauth/authorize"
const githubOTokenURL = "https://github.com/login/oauth/access_token"
const githubAPIUser = "https://api.github.com/user"
const githubScope = "user:email"

// GithubInit initiates github authentication workflow
func GithubInit() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		_ httprouter.Params) {
		url := githubClient().AuthCodeURL(state())
		http.RedirectHandler(url, http.StatusFound).ServeHTTP(w, r)
	}
}

// GithubVerify verifies github callback information, and inits
// user record
func GithubVerify(us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		_ httprouter.Params) {
		var (
			res *http.Response
			u   *model.User
		)

		r.ParseForm()
		hc, err := githubClientWithToken(r.FormValue("code"))
		if err == nil {
			res, err = hc.Get(githubAPIUser)
		}

		gu := new(model.GitHubUser)
		if err == nil {
			err = json.NewDecoder(res.Body).Decode(gu)
		}

		if u, err = us.FindByID(gu.ID); err != nil {
			// Add new user.
			u = gu.ToUser()
			err = us.Add(u)
		} else {
			// Update user from GitHub.
			u.Name = gu.Name
			u.Email = gu.Email
			u.AvatarURL = gu.AvatarURL
			err = us.Update(u)
		}

		if err == nil {
			// @TODO(Akeda)
			//
			// Make url dynamic based on the `url_redirect` query string before
			// authorising with GitHub.
			//
			// Before that, make sure HTTP router has web handler in addition
			// to API handlers
			url := fmt.Sprintf("http://localhost:8080/#api_key=%s", u.APIKey)
			http.RedirectHandler(url, http.StatusFound).ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
}

func githubClientWithToken(c string) (*http.Client, error) {
	var (
		tok *oauth2.Token
		err error
	)
	gh := githubClient()
	if tok, err = gh.Exchange(oauth2.NoContext, c); err != nil {
		return nil, err
	}
	if !tok.Valid() {
		return nil, fmt.Errorf("INVALID TOKEN: %+v", tok)
	}

	return gh.Client(oauth2.NoContext, tok), nil
}

func githubClient() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		Scopes:       strings.Split(githubScope, ","),
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubOAuthURL,
			TokenURL: githubOTokenURL,
		},
	}
}

func state() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", md5.Sum(b))
}
