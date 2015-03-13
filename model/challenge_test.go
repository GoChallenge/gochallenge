package model_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestChallengeMarshal(t *testing.T) {
	c := model.Challenge{
		ID:     10,
		Name:   "The Challenge",
		Import: "http://github.com/gochallenge",
		Status: model.Open,
	}
	s := strings.Replace(`
{
"id":10,
"name":"The Challenge",
"author":{"name":""},
"url":"",
"import":"http://github.com/gochallenge",
"status":"open"
}
`, "\n", "", -1)

	b, err := json.Marshal(c)
	require.NoError(t, err, "Challenge JSON marshalling failed")
	require.Equal(t, s, string(b), "Challenge JSON is incorrect")

	c1 := model.Challenge{}
	err = json.Unmarshal(b, &c1)
	require.NoError(t, err, "Challenge JSON unmarshalling failed")
	require.Equal(t, c, c1, "Challenge JSON unmarshalled incorrectly")
}
