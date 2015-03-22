package model

import "net/http"

// GithubAPI is an interface defining available methods of a Github API client
// implementation
type GithubAPI interface {
	AuthURL(string) string
	NewClientWithToken(string) (*http.Client, error)
	User(*http.Client) (*GithubUser, error)
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
