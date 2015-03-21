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
	a := api.New(api.Config{
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	s0 := model.Submission{
		ID:   "0000-abcd",
		Type: model.LvlNormal,
		Data: []byte("badc0ffee"),
	}
	ss.Add(&s0)

	path := fmt.Sprintf("/v1/submissions/%s/download", s0.ID)
	res, err := http.Get(ts.URL + path)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	require.Equal(t, "application/zip", res.Header.Get("Content-Type"))

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, string(s0.Data), string(b))
}

func TestDownloadSubmissionMissing(t *testing.T) {
	ss := mock.NewSubmissions()
	a := api.New(api.Config{
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/submissions/%s/download", "123")
	res, err := http.Get(ts.URL + path)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
