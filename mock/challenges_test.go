package mock_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestChallengeMockRepo(t *testing.T) {
	cs := mock.Challenges{}
	c1 := model.Challenge{
		ID: 123,
	}
	cs.Add(c1)

	// Search for First should return the challenge data
	c, err := cs.Find(c1.ID)
	require.NoError(t, err, "existing challenge lookup should not error")
	require.Equal(t, c, c1, "existing challenge should be returned")

	// Current challenge should return an error, as it doesn't exist
	c, err = cs.Find(mock.CurrentID)
	require.Error(t, err, "unknown challenge lookup should error")

	// Current challenge should return the correct one, after it has
	// been added
	c0 := model.Challenge{
		ID: mock.CurrentID,
	}
	cs.Add(c0)
	c, err = cs.Find(c0.ID)
	require.NoError(t, err, "current challenge lookup should not error")
	require.Equal(t, c, c0, "current challenge should be returned")
}
