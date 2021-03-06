package auth

import (
	"net/http"

	"github.com/gochallenge/gochallenge/model"
)

// HTTPHeader containing API key for the request
const HTTPHeader = "Auth-ApiKey"

// User that is authenticated in the given request
func User(r *http.Request, us model.Users) (*model.User, error) {
	k := r.Header.Get(HTTPHeader)
	if k == "" {
		return nil, model.ErrAuthFailure
	}

	u, err := us.FindByAPIKey(k)
	if err == model.ErrNotFound {
		// If user is not found here - it means the API key is wrong,
		// so rewriting it back as auth failure error
		err = model.ErrAuthFailure
	}
	return u, err
}
