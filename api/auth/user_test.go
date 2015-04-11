package auth_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestUserAuth(t *testing.T) {
	u := &model.User{
		ID:     5,
		Name:   "Jane Doe",
		APIKey: "c001c0ffee",
	}
	us := mock.NewUsers()
	us.Save(u)

	req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))
	require.NoError(t, err)

	// when no auth header is set - auth error
	_, err = auth.User(req, &us)
	require.Equal(t, model.ErrAuthFailure, err)

	// when a correct auth header is set - valid user returned
	// req.Header["Auth-ApiKey"] = []string{u.APIKey}
	req.Header.Set("Auth-ApiKey", u.APIKey)
	u0, err := auth.User(req, &us)
	require.NoError(t, err)
	require.Equal(t, u, u0)

	// when invalid auth header is set - auth error
	req.Header.Set("Auth-ApiKey", "badc0ffee")
	u0, err = auth.User(req, &us)
	require.Equal(t, model.ErrAuthFailure, err)
}
