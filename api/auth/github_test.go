package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestGetGithub(t *testing.T) {
	gh := mock.NewGithub()
	a := api.New(api.Config{
		Github: &gh,
	})
	ts := httptest.NewServer(a)

	// A bit of a get to make GET request without following
	// a redirect
	trn := &http.Transport{}
	req, _ := http.NewRequest("GET", ts.URL+"/v1/auth/github", nil)
	res, err := trn.RoundTrip(req)

	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, "302 Found", res.Status)

	l := res.Header.Get("Location")
	require.NotEmpty(t, l)
}

func TestVerifyNewGithubUser(t *testing.T) {
	gh := mock.NewGithub()
	gu := &model.GithubUser{
		ID:        12134,
		Name:      "Jane Doe",
		Email:     "jd@mailinator.com",
		AvatarURL: "http://localhost/avatar.png",
	}
	gh.SetUser(gu)

	us := mock.NewUsers()
	a := api.New(api.Config{
		Github: &gh,
		Users:  &us,
	})
	ts := httptest.NewServer(a)

	trn := &http.Transport{}
	req, _ := http.NewRequest("GET",
		ts.URL+"/v1/auth/github_verify?code=123456", nil)
	res, err := trn.RoundTrip(req)
	defer res.Body.Close()

	// successful verification should redirect to the authenticated page
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, res.StatusCode)

	// for a new user, we create a record and populate it with
	// some basic details from Github profile
	u, err := us.FindByGithubID(gu.ID)
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, gu.ID, u.GithubID)
	require.Equal(t, gu.Name, u.Name)
	require.Equal(t, gu.Email, u.Email)
	require.Equal(t, gu.AvatarURL, u.AvatarURL)
	require.NotEmpty(t, u.APIKey)

	loc := res.Header.Get("Location")
	require.Equal(t, "/#api_key="+u.APIKey, loc)
}

func TestVerifyGithubFailed(t *testing.T) {
	gh := mock.NewGithub()
	// we don't setup github user here, so mock API calls should error
	us := mock.NewUsers()
	a := api.New(api.Config{
		Github: &gh,
		Users:  &us,
	})
	ts := httptest.NewServer(a)

	trn := &http.Transport{}
	req, _ := http.NewRequest("GET",
		ts.URL+"/v1/auth/github_verify?code=123456", nil)
	res, err := trn.RoundTrip(req)
	defer res.Body.Close()

	// verification errors are reported back to the client page
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, res.StatusCode)

	loc := res.Header.Get("Location")
	require.Equal(t, "/#error=Error+communicating+with+Github+API", loc)
}
