package submissions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// Post new sumission
func Post(cs model.Challenges, ss model.Submissions) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var (
			c model.Challenge
			s model.Submission
			b []byte
		)
		c, err := findChallenge(cs, ps.ByName("id"))

		if err == nil {
			b, err = ioutil.ReadAll(r.Body)
		}
		if err == nil {
			err = json.Unmarshal(b, &s)
			s.ChallengeID = c.ID
			s.Challenge = &c
			ss.Add(&s)
		}
		if err == nil {
			b, err = json.Marshal(s)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			w.Write(b)
		}
	}
}

// find a challenge given the value of requested ID string
func findChallenge(cs model.Challenges, id string) (model.Challenge, error) {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return model.Challenge{}, err
	}
	return cs.Find(cid)
}
