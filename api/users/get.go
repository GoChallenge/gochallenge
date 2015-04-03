package users

import (
	"encoding/json"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

func Get(us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		u, err := findUser(us, ps.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// find a user of a specified ID.
func findUser(us model.Users, idstr string) (*model.User, error) {
	var (
		id  model.UserID
		u   *model.User
		err error
	)

	err = id.Atoid(idstr)
	if err != nil {
		return nil, err
	}
	u, err = us.Find(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}
