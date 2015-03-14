package submissions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestPostMultipart(t *testing.T) {
	ss := mock.NewSubmissions()
	cs := mock.NewChallenges()

	c0 := model.Challenge{
		ID:     1,
		Status: model.Open,
	}
	cs.Add(c0)

	s0 := model.Submission{
		ID: 1,
	}
	ss.Add(&s0)

	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", c0.ID)
	buf := strings.NewReader("{}")
	req, err := http.NewRequest("POST", ts.URL+path, buf)

	hc := &http.Client{}
	res, err := hc.Do(req)
	defer res.Body.Close()

	require.NoError(t, err, "POST /v1/.../submissions should not error")
	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "POST /v1/.../submissions should read the body")

	require.Equal(t, "200 OK", res.Status,
		fmt.Sprintf("POST /v1/.../submissions returned error %s, %s",
			res.Status, b))

	var sx model.Submission
	err = json.Unmarshal(b, &sx)
	require.NoError(t, err, "POST /v1/.../submissions unmarshaling failed")
	require.Equal(t, 2, sx.ID,
		"GET /v1/.../submissions unmarshalled incorrectly")
}
