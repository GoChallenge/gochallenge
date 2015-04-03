package submissions_test

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
	t.SkipNow()

	cs := mock.NewChallenges()
	ss := mock.NewSubmissions()
	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	c0 := model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Open,
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(&c0)

	s0 := model.Submission{
		ID:        "0000-abcd",
		Challenge: &c0,
		Type:      model.LvlNormal,
	}
	ss.Add(&s0)
	s1 := model.Submission{
		ID:        "0000-fedc",
		Challenge: &c0,
		Type:      model.LvlFun,
	}
	ss.Add(&s1)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", c0.ID)
	res, err := http.Get(ts.URL + path)
	defer res.Body.Close()

	require.NoError(t, err, "GET /v1/.../submissions should not error")
	require.Equal(t, "200 OK", res.Status,
		fmt.Sprintf("GET /v1/.../submissions error code %s", res.Status))

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "GET /v1/.../submissions should read the body")

	var sx []model.Submission
	err = json.Unmarshal(b, &sx)
	require.NoError(t, err, "GET /v1/.../submissions unmarshaling failed")
	require.Equal(t, []model.Submission{s0, s1}, sx,
		"GET /v1/.../submissions unmarshalled incorrectly")
}

func TestListEmpty(t *testing.T) {
	cs := mock.NewChallenges()
	ss := mock.NewSubmissions()
	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	c0 := model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Open,
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(&c0)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", c0.ID)
	res, err := http.Get(ts.URL + path)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, res.Header.Get("Content-Type"), "application/json")

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, "[]\n", string(b))
}

func TestListMissing(t *testing.T) {
	cs := mock.NewChallenges()
	a := api.New(api.Config{
		Challenges: &cs,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", 123)
	res, err := http.Get(ts.URL + path)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
