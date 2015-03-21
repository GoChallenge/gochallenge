package challenges

import (
	"encoding/json"
	"net/http"

	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// List all available challenges
func List(cs model.Challenges) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		cx, err := cs.All()
		err = writeChallenges(err, w, &cx)

		if err != nil {
			write.Error(w, r, err)
		}
	}
}

func writeChallenges(err error, w http.ResponseWriter, cx *[]model.Challenge) error {
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(cx)
}
