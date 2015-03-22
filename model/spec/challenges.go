package spec

import (
	"testing"
	"time"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

// MustBehaveLikeChallenges tests behaviour of the given challenges
// repo, to make sure it conforms to the expected API
func MustBehaveLikeChallenges(t *testing.T, cs model.Challenges) {
	const chalID = 123
	const chalUnknownID = 1

	c1 := model.Challenge{
		ID:    chalID,
		Name:  "The Test Challenge",
		Start: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2044, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, cs.Add(c1))

	// Search for First should return the challenge data
	c, err := cs.Find(c1.ID)
	require.NoError(t, err, "existing challenge lookup should not error")
	require.Equal(t, c, c1, "existing challenge should be returned")

	// Current challenge should return an error, as it doesn't exist
	c, err = cs.Find(chalUnknownID)
	require.Equal(t, model.ErrNotFound, err)

	// Current challenge should return the correct one, after it has
	// been added
	c0 := model.Challenge{
		ID: mock.CurrentID,
	}
	cs.Add(c0)
	c, err = cs.Find(c0.ID)
	require.NoError(t, err, "current challenge lookup should not error")
	require.Equal(t, c, c0, "current challenge should be returned")

	// All should return all added challenges
	cx, err := cs.All()
	require.NoError(t, err, "all challenges returned an error")
	require.Equal(t, []model.Challenge{c1, c0}, cx,
		"all challenges not returned correctly")
}
