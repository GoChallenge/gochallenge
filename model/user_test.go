package model_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestUserResetToken(t *testing.T) {
	u, err := model.NewUser()
	require.NoError(t, err)

	k := u.APIKey
	u.ResetToken()
	require.NotEmpty(t, u.APIKey)
	require.NotEqual(t, u.APIKey, k, "ResetToken must reset API key")
}
