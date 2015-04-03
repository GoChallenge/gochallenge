package spec

import (
	"testing"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

// MustBehaveLikeUsers tests behaviour of the given implementation
// of users repo
func MustBehaveLikeUsers(t *testing.T, us model.Users) {
	// Create and save a user
	u1 := model.User{
		Name:     "Jane Doe",
		GithubID: 66235,
		APIKey:   "deadc0ffee",
	}
	err := us.Save(&u1)
	require.NoError(t, err, "errored when saving a user")
	// ID should be auto-generated, if not specified
	require.NotEmpty(t, u1.ID, "User ID should be auto-generated")

	// added user record should find-able by its ID
	ux, err := us.Find(u1.ID)
	require.NoError(t, err, "errored when finding a user")
	require.Equal(t, u1, *ux)
	// but not if it's a wrong one
	ux, err = us.Find(u1.ID * 100)
	require.Equal(t, model.ErrNotFound, err)

	// and by its Github ID
	ux, err = us.FindByGithubID(u1.GithubID)
	require.NoError(t, err)
	require.Equal(t, u1, *ux)
	// but not if it's a wrong one
	ux, err = us.FindByGithubID(u1.GithubID * 100)
	require.Equal(t, model.ErrNotFound, err)

	// user, when added, should have the API key generated
	ak := u1.APIKey
	require.NotEmpty(t, ak)
	// which can be used to find the same user
	ux, err = us.FindByAPIKey(ak)
	require.NoError(t, err)
	require.Equal(t, u1, *ux)
	// but, again, not if it's a wrong one
	ux, err = us.FindByAPIKey("o_O")
	require.Equal(t, model.ErrNotFound, err)

	// the second user, when added, must receive a different ID
	// Create and save a user
	u2 := model.User{
		Name: "Gordon Freeman",
	}
	err = us.Save(&u2)
	require.NoError(t, err)
	// ID should be auto-generated, if not specified
	require.NotEqual(t, u2.ID, 0, "User ID should be auto-generated")
	require.NotEqual(t, u1.ID, u2.ID)
}
