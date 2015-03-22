package github_test

import (
	"net/url"
	"testing"

	"github.com/gochallenge/gochallenge/github"
	"github.com/stretchr/testify/require"
)

func TestAuthURL(t *testing.T) {
	gh := github.NewClient()
	loc := gh.AuthURL("c0ffee")
	require.Contains(t, loc, "https://github.com/login/oauth/authorize")
	require.Contains(t, loc, "c0ffee")

	u, _ := url.ParseRequestURI(loc)
	q := u.Query()
	require.Equal(t, "user:email", q.Get("scope"))
	require.NotEmpty(t, q.Get("state"))
}
