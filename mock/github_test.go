package mock_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestGithubMock(t *testing.T) {
	gh := mock.NewGithub()
	require.Contains(t, gh.AuthURL("hello"), "hello",
		"auth url should include state string")

	c, err := gh.NewClientWithToken("faketoken")
	require.NoError(t, err)

	// before a user is set - call to User API should error
	ux, err := gh.User(c)
	require.Equal(t, err, model.ErrGithubAPIError)

	// after a user is set - the call should succeed
	u0 := &model.GithubUser{
		ID: 12345,
	}
	gh.SetUser(u0)
	ux, err = gh.User(c)
	require.NoError(t, err)
	require.Equal(t, u0, ux, "pre-set user should be returned")
}
