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
	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func getList(ts *httptest.Server, cid model.ChallengeID,
	u *model.User) (*http.Response, error) {
	path := fmt.Sprintf("/v1/challenges/%d/submissions", cid)
	req, err := http.NewRequest("GET", ts.URL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(auth.HTTPHeader, u.APIKey)

	hc := &http.Client{}
	return hc.Do(req)
}

func TestList(t *testing.T) {
	cs := mock.NewChallenges()
	ss := mock.NewSubmissions()
	us := mock.NewUsers()
	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
		Users:       &us,
	})
	ts := httptest.NewServer(a)

	c0 := &model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Open,
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(c0)

	u0, err := model.NewUser()
	us.Save(u0)
	require.NoError(t, err)
	s0 := model.Submission{
		ID:        "0000-abcd",
		Challenge: c0,
		User:      u0,
		Type:      model.LvlNormal,
	}
	ss.Add(&s0)

	u1, err := model.NewUser()
	us.Save(u0)
	require.NoError(t, err)
	s1 := model.Submission{
		ID:        "0000-fedc",
		User:      u1,
		Challenge: c0,
		Type:      model.LvlFun,
	}
	ss.Add(&s1)

	res, err := getList(ts, c0.ID, u0)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, "200 OK", res.Status)

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	b01, _ := json.Marshal([]model.Submission{s0})
	require.Equal(t, string(b01)+"\n", string(b))
}

func TestListEmpty(t *testing.T) {
	cs := mock.NewChallenges()
	ss := mock.NewSubmissions()
	us := mock.NewUsers()
	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
		Users:       &us,
	})
	ts := httptest.NewServer(a)

	u0, err := model.NewUser()
	require.NoError(t, err)
	us.Save(u0)
	c0 := &model.Challenge{
		ID:     123,
		Name:   "The Challenge",
		Status: model.Open,
		Start:  time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 3, 14, 0, 0, 0, 0, time.UTC),
	}
	cs.Save(c0)

	res, err := getList(ts, c0.ID, u0)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, res.Header.Get("Content-Type"), "application/json")

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "[]\n", string(b))
}

func TestListMissing(t *testing.T) {
	cs := mock.NewChallenges()
	us := mock.NewUsers()
	a := api.New(api.Config{
		Challenges: &cs,
		Users:      &us,
	})
	ts := httptest.NewServer(a)

	u0, err := model.NewUser()
	require.NoError(t, err)
	us.Save(u0)

	res, err := getList(ts, 123, u0)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
