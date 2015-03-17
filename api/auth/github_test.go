package auth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gochallenge/gochallenge/api"
	"github.com/stretchr/testify/require"
)

func TestGetGithub(t *testing.T) {
	a := api.New(api.Config{})
	ts := httptest.NewServer(a)

	// A bit of a get to make GET request without following
	// a redirect
	trn := &http.Transport{}
	req, _ := http.NewRequest("GET", ts.URL+"/v1/auth/github", nil)
	res, err := trn.RoundTrip(req)

	require.NoError(t, err, "GET /v1/auth/github should not error")
	defer res.Body.Close()

	require.Equal(t, "302 Found", res.Status,
		"GET /v1/auth/github invalid error code")

	l := res.Header.Get("Location")
	require.Contains(t,
		l, "https://github.com/login/oauth/authorize",
		"GET /v1/auth/github invalid redirect")

	u, _ := url.ParseRequestURI(l)
	q := u.Query()
	require.Equal(t, "user:email", q.Get("scope"),
		"Github auth requested invalid scope")
	require.NotEmpty(t, q.Get("state"))
}
