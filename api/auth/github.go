package auth

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

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
func GithubVerify() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		_ httprouter.Params) {
		var (
			res *http.Response
			b   []byte
		)

		r.ParseForm()
		hc, err := githubClientWithToken(r.FormValue("code"))
		if err == nil {
			res, err = hc.Get(githubAPIUser)
		}
		if err == nil {
			b, err = ioutil.ReadAll(res.Body)
		}
		if err == nil {
			w.Write(b)
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
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", md5.Sum(b))
}
