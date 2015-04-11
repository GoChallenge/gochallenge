package submissions

import (
	"encoding/json"
	"net/http"

	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// List all available submissions for a challenge
func List(cs model.Challenges, ss model.Submissions,
	us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var sx []*model.Submission

		u, err := auth.User(r, us)
		c, err := findChallenge(err, cs, ps.ByName("id"))
		sx, err = listSubmissions(err, ss, c, u)
		err = writeSubmissions(err, w, sx)

		if err != nil {
			write.Error(w, r, err)
		}
	}
}

func listSubmissions(err error, ss model.Submissions,
	c *model.Challenge, u *model.User) ([]*model.Submission, error) {
	var subs []*model.Submission
	// non-initialised array will marshal to "null", not to "[]",
	// so this is a little hack to work around that
	subs = make([]*model.Submission, 0)

	if err != nil {
		return subs, err
	}
	sx, err := ss.AllForChallenge(c)
	if err != nil {
		return subs, err
	}

	for _, s := range sx {
		if readable(u, s) {
			subs = append(subs, s)
		}
	}
	return subs, nil
}

func writeSubmissions(err error, w http.ResponseWriter,
	sx []*model.Submission) error {

	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(sx)
}

func readable(u *model.User, s *model.Submission) bool {
	return s.User != nil && u.ID == s.User.ID
}
