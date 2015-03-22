package github

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gochallenge/gochallenge/model"
	"golang.org/x/oauth2"
)

const githubOAuthURL = "https://github.com/login/oauth/authorize"
const githubOTokenURL = "https://github.com/login/oauth/access_token"
const githubAPIUser = "https://api.github.com/user"
const githubScope = "user:email"

// NewClient return new configured Github client
func NewClient() github {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		Scopes:       strings.Split(githubScope, ","),
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubOAuthURL,
			TokenURL: githubOTokenURL,
		},
	}
	return github{
		config: conf,
	}
}

// Github is an implementation of Github API talking to the Github server
type github struct {
	config *oauth2.Config
}

// AuthURL generates authentication URL to redirect user to,
// using provided string as the state
func (gh *github) AuthURL(s string) string {
	return gh.config.AuthCodeURL(s)
}

// NewClientWithToken returns an http.Client that can be used
// for authenticated communications with Github API
func (gh *github) NewClientWithToken(t string) (*http.Client, error) {
	var (
		tok *oauth2.Token
		err error
	)
	gc := gh.config
	if tok, err = gc.Exchange(oauth2.NoContext, t); err != nil || !tok.Valid() {
		return nil, model.ErrGithubAPIError
	}

	return gc.Client(oauth2.NoContext, tok), nil
}

// User details for the user that is currently authenticated
func (gh github) User(hc *http.Client) (*model.GithubUser, error) {
	var gu model.GithubUser

	res, err := hc.Get(githubAPIUser)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&gu)
	return &gu, err
}
