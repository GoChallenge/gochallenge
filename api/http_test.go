package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morhekil/gochallenge/api"
	"github.com/morhekil/gochallenge/mock"
	"github.com/morhekil/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestGetChallenge(t *testing.T) {
	cs := mock.Challenges{}
	a := api.New(api.Config{
		Challenges: cs,
	})
	ts := httptest.NewServer(a)

	c0 := model.Challenge{
		ID:   "First",
		Name: "The Challenge",
	}
	cs.Add(c0)

	path := "/v1/challenge/" + c0.ID
	res, err := http.Get(ts.URL + path)

	require.NoError(t, err, "GET "+path+" should not error")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "GET "+path+" should read the body")
	require.Contains(t, string(body),
		c0.Name, "GET "+path+" should return the challenge")
}

func TestGetCurrentChallenge(t *testing.T) {
	cs := mock.Challenges{}
	a := api.New(api.Config{
		Challenges: cs,
	})
	ts := httptest.NewServer(a)

	c0 := model.Challenge{
		ID:   mock.CurrentID,
		Name: "The Current Challenge",
	}
	cs.Add(c0)

	path := "/v1/challenge/current"
	res, err := http.Get(ts.URL + path)

	require.NoError(t, err, "GET "+path+" should not error")
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "GET "+path+" should read the body")
	require.Contains(t, string(body),
		c0.Name, "GET "+path+" should return the challenge")
}
