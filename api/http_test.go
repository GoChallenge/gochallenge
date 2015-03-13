package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func testGettingChallenge(t *testing.T, path string, c0 model.Challenge) {
	cs := mock.Challenges{}
	a := api.New(api.Config{
		Challenges: cs,
	})
	ts := httptest.NewServer(a)
	cs.Add(c0)

	res, err := http.Get(ts.URL + path)

	require.NoError(t, err, "GET "+path+" should not error")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "GET "+path+" should read the body")

	c1 := model.Challenge{}
	err = json.Unmarshal(body, &c1)

	require.NoError(t, err, "GET "+path+" unmarshal errored")
	require.Equal(t, c0, c1, "GET "+path+" unmarshalled incorrectly")
}

func TestGetChallenge(t *testing.T) {
	c0 := model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Open,
	}
	path := fmt.Sprintf("/v1/challenge/%d", c0.ID)

	testGettingChallenge(t, path, c0)
}

func TestGetCurrentChallenge(t *testing.T) {
	c0 := model.Challenge{
		ID:   mock.CurrentID,
		Name: "The Current Challenge",
	}
	path := "/v1/challenge/current"

	testGettingChallenge(t, path, c0)
}
