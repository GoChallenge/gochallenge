package submissions

import (
	"encoding/json"
	"net/http"

	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// List all available submissions for a challenge
func List(cs model.Challenges, ss model.Submissions) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var sx []*model.Submission

		c, err := findChallenge(cs, ps.ByName("id"))
		sx, err = listSubmissions(err, ss, &c)
		err = writeSubmissions(err, w, sx)

		if err != nil {
			write.Error(w, r, err)
		}
	}
}

func listSubmissions(err error, ss model.Submissions,
	c *model.Challenge) ([]*model.Submission, error) {

	if err != nil {
		return nil, err
	}
	return ss.AllForChallenge(c)
}

func writeSubmissions(err error, w http.ResponseWriter,
	sx []*model.Submission) error {

	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(sx)
}
