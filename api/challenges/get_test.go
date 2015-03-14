package challenges_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func getBody(t *testing.T, path string, c0 model.Challenge) []byte {
	cs := mock.NewChallenges()
	a := api.New(api.Config{
		Challenges: &cs,
	})
	ts := httptest.NewServer(a)
	cs.Add(c0)

	res, err := http.Get(ts.URL + path)

	require.NoError(t, err, "GET "+path+" should not error")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "GET "+path+" should read the body")

	return body
}

func testGettingChallenge(t *testing.T, path string, c0 model.Challenge) {
	body := getBody(t, path, c0)

	c1 := model.Challenge{}
	err := json.Unmarshal(body, &c1)

	require.NoError(t, err, "GET "+path+" unmarshal errored")
	require.Equal(t, c0, c1, "GET "+path+" unmarshalled incorrectly")
}

func TestGetChallenge(t *testing.T) {
	c0 := model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Open,
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	path := fmt.Sprintf("/v1/challenges/%d", c0.ID)

	testGettingChallenge(t, path, c0)
}

func TestGetCurrentChallenge(t *testing.T) {
	c0 := model.Challenge{
		ID:    mock.CurrentID,
		Name:  "The Current Challenge",
		Start: time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	path := "/v1/challenges/current"

	testGettingChallenge(t, path, c0)
}

func TestGoGetChallenge(t *testing.T) {
	c0 := model.Challenge{
		ID:     mock.CurrentID,
		Name:   "The Current Challenge",
		Import: "gochallenge.org/gochallenge-x",
		Git:    "https://github.com/author/challengex",
	}
	path := fmt.Sprintf("/code/challenge-%03d?go-get=1", mock.CurrentID)

	body := getBody(t, path, c0)

	meta := fmt.Sprintf(`<meta name="go-import" content="%s git %s">`,
		c0.Import, c0.Git)
	require.Contains(t, string(body), meta,
		"go get response did not return correct meta tag")
}
