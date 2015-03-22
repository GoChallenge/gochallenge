package spec

import (
	"testing"
	"time"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MustBehaveLikeChallenges tests behaviour of the given challenges
// repo, to make sure it conforms to the expected API
func MustBehaveLikeChallenges(t *testing.T, cs model.Challenges,
	cur *model.Challenge) {
	const chalID = 123
	const chalUnknownID = 1

	c1 := model.Challenge{
		ID:    chalID,
		Name:  "The Test Challenge",
		Start: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	// just a sanity check to make sure the first challenge is
	// not the "current" one, to not trip over it in later tests
	assert.False(t, c1.Current(), "test challenge should not be current")

	// Find should return an error before the challenge was added
	_, err := cs.Find(c1.ID)
	require.Equal(t, model.ErrNotFound, err)

	// Adding the challenge to the repo, should succeed
	require.NoError(t, cs.Add(&c1))

	// Now find should succeed, too, as the challenge has been added
	c, err := cs.Find(c1.ID)
	require.NoError(t, err, "existing challenge lookup should not error")
	require.Equal(t, *c, c1, "existing challenge should be returned")

	// Current challenge should return an error, as it doesn't exist
	c, err = cs.Current()
	require.Equal(t, model.ErrNotFound, err)

	// Current challenge should return the correct one, after it has
	// been added
	cs.Add(cur)
	c, err = cs.Current()
	require.NoError(t, err, "current challenge lookup should not error")
	require.Equal(t, *c, *cur, "current challenge should be returned")

	// All should return all added challenges
	cx, err := cs.All()
	require.NoError(t, err, "all challenges returned an error")
	require.Equal(t, 2, len(cx), "two challenges must be returned")

	cx0 := *cx[0]
	cx1 := *cx[1]
	require.True(t, (cx0 == *cur && cx1 == c1) || (cx1 == *cur && cx0 == c1),
		"saved challenges must be returned")
}
