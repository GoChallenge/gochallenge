package model

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

const githubOAuthURL = "https://github.com/login/oauth/authorize"
const githubOTokenURL = "https://github.com/login/oauth/access_token"
const githubAPIUser = "https://api.github.com/user"
const githubScope = "user:email"

// GithubAPI is an interface defining available methods of a Github API client
// implementation
type GithubAPI interface {
	AuthURL(string) string
	NewClientWithToken(string) (*http.Client, error)
	User(*http.Client) (GithubUser, error)
}

// NewGithub return new configured Github client
func NewGithub() Github {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		Scopes:       strings.Split(githubScope, ","),
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubOAuthURL,
			TokenURL: githubOTokenURL,
		},
	}
	return Github{
		config: conf,
	}
}

// Github is an implementation of Github API talking to the Github server
type Github struct {
	config *oauth2.Config
}

// AuthURL generates authentication URL to redirect user to,
// using provided string as the state
func (gh *Github) AuthURL(s string) string {
	return gh.config.AuthCodeURL(s)
}

// NewClientWithToken returns an http.Client that can be used
// for authenticated communications with Github API
func (gh *Github) NewClientWithToken(t string) (*http.Client, error) {
	var (
		tok *oauth2.Token
		err error
	)
	gc := gh.config
	if tok, err = gc.Exchange(oauth2.NoContext, t); err != nil || !tok.Valid() {
		return nil, ErrGithubAPIError
	}

	return gc.Client(oauth2.NoContext, tok), nil
}

// User details for the user that is currently authenticated
func (gh Github) User(hc *http.Client) (GithubUser, error) {
	var gu GithubUser

	res, err := hc.Get(githubAPIUser)
	defer res.Body.Close()

	if err != nil {
		return gu, err
	}

	err = json.NewDecoder(res.Body).Decode(&gu)
	return gu, err
}

// GithubUser represents a user of Github.
type GithubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

// Populate given user record with their Github account details
func (gu *GithubUser) Populate(u *User) {
	choose := func(s1 string, s2 string) string {
		if s1 == "" {
			return s2
		}
		return s1
	}
	u.GithubID = gu.ID
	u.Name = choose(u.Name, gu.Name)
	u.Email = choose(u.Email, gu.Email)
	u.AvatarURL = choose(u.AvatarURL, gu.AvatarURL)
	u.GithubURL = choose(u.GithubURL, gu.HTMLURL)
	u.GithubLogin = choose(u.GithubLogin, gu.Login)
}
