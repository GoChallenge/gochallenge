package submissions_test

import (
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

func TestDownloadSubmission(t *testing.T) {
	ss := mock.NewSubmissions()
	us := mock.NewUsers()
	a := api.New(api.Config{
		Submissions: &ss,
		Users:       &us,
	})
	ts := httptest.NewServer(a)

	u0, err := model.NewUser()
	require.NoError(t, err)
	us.Save(u0)

	d0 := []byte("badc0ffee")
	s0 := &model.Submission{
		ID:   "0000-abcd",
		Type: model.LvlNormal,
		Data: &d0,
		User: u0,
	}
	ss.Add(s0)

	res, err := get(ts, fmt.Sprintf("/submissions/%s/download", s0.ID), u0)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	require.Equal(t, "application/zip", res.Header.Get("Content-Type"))

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, string(d0), string(b))
}

func TestDownloadSubmissionMissing(t *testing.T) {
	ss := mock.NewSubmissions()
	us := mock.NewUsers()
	a := api.New(api.Config{
		Submissions: &ss,
		Users:       &us,
	})
	ts := httptest.NewServer(a)

	u0, err := model.NewUser()
	require.NoError(t, err)
	us.Save(u0)

	res, err := get(ts, "/submissions/123/download", u0)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestDownloadNoAuth(t *testing.T) {
	ss := mock.NewSubmissions()
	a := api.New(api.Config{
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/submissions/%s/download", "123")
	res, err := http.Get(ts.URL + path)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}
