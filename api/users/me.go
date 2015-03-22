package users

import (
	"encoding/json"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

func Me(us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		apiKey := r.Header.Get("Auth-ApiKey")
		u, err := us.FindByAPIKey(apiKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err = json.NewEncoder(w).Encode(u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
