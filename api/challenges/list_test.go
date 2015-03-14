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

func TestList(t *testing.T) {
	cs := mock.NewChallenges()
	a := api.New(api.Config{
		Challenges: &cs,
	})
	ts := httptest.NewServer(a)

	c0 := model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Closed,
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Add(c0)

	c1 := model.Challenge{
		ID:     124,
		Name:   "The Challenge Two",
		Status: model.Open,
		Start:  time.Date(2015, 4, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 4, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Add(c1)

	res, err := http.Get(ts.URL + "/v1/challenges")
	defer res.Body.Close()

	require.NoError(t, err, "GET /v1/challenges should not error")
	require.Equal(t, "200 OK", res.Status,
		fmt.Sprintf("GET /v1/challenges returned error code %s", res.Status))

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "GET /v1/challenges should read the body")

	var cx []model.Challenge
	err = json.Unmarshal(b, &cx)
	require.NoError(t, err, "GET /v1/challenges unmarshaling failed")
	require.Equal(t, []model.Challenge{c0, c1}, cx,
		"GET /v1/challenges unmarshalled incorrectly")
}
