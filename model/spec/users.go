package spec

import (
	"testing"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

// MustBehaveLikeUsers tests behaviour of the given implementation
// of users repo
func MustBehaveLikeUsers(t *testing.T, us model.Users) {
	u1 := model.User{
		ID:       1,
		Name:     "Jane Doe",
		GithubID: 66235,
		APIKey:   "deadc0ffee",
	}
	// first adding should succeed
	err := us.Add(&u1)
	require.NoError(t, err)
	// but the second should fail with a duplicate record error
	err = us.Add(&u1)
	require.Equal(t, model.ErrDuplicateRecord, err,
		"adding a duplicate record did not error")

	// added user record should find-able by its ID
	ux, err := us.FindByID(u1.ID)
	require.NoError(t, err)
	require.Equal(t, u1, *ux)
	// but not if it's a wrong one
	ux, err = us.FindByID(u1.ID * 100)
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
}
