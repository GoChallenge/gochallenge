package auth

import (
	"net/http"

	"github.com/gochallenge/gochallenge/model"
)

// User that is authenticated in the given request
func User(r *http.Request, us model.Users) (*model.User, error) {
	k := r.Header.Get("Auth-ApiKey")
	if k == "" {
		return nil, model.ErrAuthFailure
	}

	u, err := us.FindByAPIKey(k)
	if err == model.ErrNotFound {
		err = model.ErrAuthFailure
	}
	return u, err
}
