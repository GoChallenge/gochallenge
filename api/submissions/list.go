package submissions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// List all available submissions for a challenge
func List(cs model.Challenges, ss model.Submissions) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var (
			b  []byte
			sx []*model.Submission
		)

		c, err := findChallenge(cs, ps.ByName("id"))

		if err == nil {
			sx, err = ss.AllForChallenge(&c)
		}

		if err == nil {
			b, err = json.Marshal(sx)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			w.Write(b)
		}
	}
}
