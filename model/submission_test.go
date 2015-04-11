package model_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestSubmissionMarshal(t *testing.T) {
	cs := mock.NewChallenges()
	us := mock.NewUsers()
	c := &model.Challenge{
		ID: 10,
	}
	cs.Save(c)
	u := &model.User{
		ID:   5,
		Name: "Jane Doe",
	}
	us.Save(u)

	s := model.Submission{
		ID:        "1234-abcde",
		Type:      model.LvlAnonymous,
		Challenge: c,
		User:      u,
		Created:   time.Date(2015, 3, 1, 10, 0, 0, 0, time.UTC),
	}
	js := strings.Replace(`
{
"id":"1234-abcde",
"user_id":5,
"challenge_id":10,
"type":"anonymous",
"created":"2015-03-01T10:00:00Z"
}
`, "\n", "", -1)

	b, err := json.Marshal(s)
	require.NoError(t, err, "Submission JSON marshalling failed")
	require.Equal(t, js, string(b), "Submission JSON is incorrect")

	sx := model.Submission{}
	err = sx.Unmarshal(b, &cs, &us)

	require.NoError(t, err, "Submission JSON unmarshalling failed")
	require.Equal(t, s, sx, "Submission JSON unmarshalled incorrectly")
}
