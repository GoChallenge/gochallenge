package submissions_test

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestPostMultipart(t *testing.T) {
	ss := mock.NewSubmissions()
	cs := mock.NewChallenges()
	us := mock.NewUsers()

	c0 := &model.Challenge{
		ID:     1,
		Status: model.Open,
	}
	cs.Save(c0)

	// Creating a pre-existing submission, to make sure
	// an ID for the new one will be assigned correctly
	ss.Add(&model.Submission{
		ID: "1",
	})

	u0, err := model.NewUser()
	require.NoError(t, err)
	us.Save(u0)

	a := api.New(api.Config{
		Challenges:  &cs,
		Submissions: &ss,
		Users:       &us,
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
	res, err := postAsUser(t, u0.APIKey, ts.URL+path, bnd, buf)
	defer res.Body.Close()

	require.NoError(t, err)
	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, "200 OK", res.Status)

	var sx model.Submission
	err = json.Unmarshal(b, &sx)
	require.NoError(t, err)
	require.Equal(t, "2", sx.ID)

	sl, err := ss.Find(sx.ID)
	testSubmissionData(t, sl, map[string]string{
		"test.txt": "test\x0a",
	})
	require.Equal(t, u0, sl.User)
	require.Equal(t, c0, sl.Challenge)
}

func testSubmissionData(t *testing.T, sx *model.Submission, ex map[string]string) {
	var b []byte

	// Test that the data can be unzipped
	z, err := zip.NewReader(bytes.NewReader(*sx.Data), int64(len(*sx.Data)))
	require.NoError(t, err, "zip reader init failed")
	files := map[string]string{}

	// Load all files into a map
	for _, f := range z.File {
		zf, err := f.Open()
		if err == nil {
			b, err = ioutil.ReadAll(zf)
			files[f.Name] = string(b)
		}
		require.NoError(t, err)
	}

	// And verify that the map is the same as the expected one
	require.Equal(t, ex, files)
}

func TestPostToWrongChallenge(t *testing.T) {
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

	path := fmt.Sprintf("/v1/challenges/%d/submissions", 123)
	buf := strings.NewReader("somedata")
	res, err := postAsUser(t, u0.APIKey, ts.URL+path, "xxx", buf)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestPostWithInvalidKey(t *testing.T) {
	cs := mock.NewChallenges()
	c0 := &model.Challenge{
		ID:     1,
		Status: model.Open,
	}
	cs.Save(c0)
	us := mock.NewUsers()
	a := api.New(api.Config{
		Challenges: &cs,
		Users:      &us,
	})
	ts := httptest.NewServer(a)

	path := fmt.Sprintf("/v1/challenges/%d/submissions", c0.ID)
	buf := strings.NewReader("somedata")
	res, err := postAsUser(t, "deadbeef", ts.URL+path, "xxx", buf)
	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func postAsUser(t *testing.T, ukey string, uri string, bnd string,
	body io.Reader) (*http.Response, error) {

	r, err := http.NewRequest("POST", uri, body)
	require.NoError(t, err)
	r.Header.Set("Content-Type", "multipart/related; boundary="+bnd)
	r.Header.Set(auth.HTTPHeader, ukey)

	client := &http.Client{}
	return client.Do(r)
}
