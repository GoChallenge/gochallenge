package submissions_test

import (
	"archive/zip"
	"bytes"
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
		ID: "1",
	}
	ss.Add(&s0)

	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", c0.ID)
	bnd := "c0ffee"

	// the data here is a zipped file "test.txt", containing word "test"
	data := fmt.Sprintf(`--%[1]s
Content-Type: application/json; charset=UTF-8

{"type":"normal"}
--%[1]s
Content-Type: application/zip
Content-Transfer-Encoding: base64

UEsDBAoAAAAAAONob0bGNbk7BQAAAAUAAAAIABwAdGVzdC50eHRVVAkAA0rpBFVO6QRVdXgLAAEE9QEA
AAQUAAAAdGVzdApQSwECHgMKAAAAAADjaG9GxjW5OwUAAAAFAAAACAAYAAAAAAABAAAApIEAAAAAdGVz
dC50eHRVVAUAA0rpBFV1eAsAAQT1AQAABBQAAABQSwUGAAAAAAEAAQBOAAAARwAAAAAA
--%[1]s--
`, bnd)
	buf := strings.NewReader(data)
	req, err := http.NewRequest("POST", ts.URL+path, buf)
	req.Header.Add("Content-Type", "multipart/related; boundary="+bnd)

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
	require.Equal(t, "2", sx.ID,
		"GET /v1/.../submission unmarshalled incorrectly")

	sl, err := ss.Find(sx.ID)
	testSubmissionData(t, sl, map[string]string{
		"test.txt": "test\x0a",
	})
}

func testSubmissionData(t *testing.T, sx *model.Submission, ex map[string]string) {
	var b []byte

	// Test that the data can be unzipped
	z, err := zip.NewReader(bytes.NewReader(sx.Data), int64(len(sx.Data)))
	require.NoError(t, err, "zip reader init failed")
	files := map[string]string{}

	// Load all files into a map
	for _, f := range z.File {
		zf, err := f.Open()
		if err == nil {
			b, err = ioutil.ReadAll(zf)
			files[f.Name] = string(b)
		}
		require.NoError(t, err, "reading files from zip data failed")
	}

	// And verify that the map is the same as the expected one
	require.Equal(t, ex, files, "GET /v1/.../submission data not loaded")
}

func TestPostToWrongChallenge(t *testing.T) {
	cs := mock.NewChallenges()
	a := api.New(api.Config{
		Challenges: &cs,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", 123)
	buf := strings.NewReader("somedata")
	res, err := http.Post(ts.URL+path, "multipart/related; boundary=xxx", buf)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
