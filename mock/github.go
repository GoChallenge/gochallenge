package mock

import (
	"net/http"

	"github.com/gochallenge/gochallenge/model"
)

// NewGithub mock
func NewGithub() Github {
	return Github{}
}

// Github implements a mock for Github API client
type Github struct {
	user *model.GithubUser
}

// AuthURL is a fake authentication URL that includes state string
func (gh *Github) AuthURL(s string) string {
	return "http://localhost?state=" + s
}

// NewClientWithToken stubs token exchange, and just returns a plain
// http client to simulate real client's behaviour
func (gh *Github) NewClientWithToken(t string) (*http.Client, error) {
	hc := &http.Client{}
	return hc, nil
}

// User returns currently configured GithubUser to fake authentication
// process, or ErrGithubAPIError if the users hasn't been set
func (gh *Github) User(hc *http.Client) (*model.GithubUser, error) {
	if gh.user == nil {
		return nil, model.ErrGithubAPIError
	}

	return gh.user, nil
}

// SetUser to be considered a currently authenticated user for mock client
func (gh *Github) SetUser(u *model.GithubUser) {
	gh.user = u
}
