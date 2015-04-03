package challenges_test

import (
	"encoding/json"
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
	cs.Save(&c0)

	c1 := model.Challenge{
		ID:     124,
		Name:   "The Challenge Two",
		Status: model.Open,
		Start:  time.Date(2015, 4, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 4, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(&c1)

	res, err := http.Get(ts.URL + "/v1/challenges")
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, res.Header.Get("Content-Type"), "application/json")

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var cx []model.Challenge
	err = json.Unmarshal(b, &cx)
	require.NoError(t, err)
	require.Equal(t, []model.Challenge{c0, c1}, cx)
}
