package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

type writerFunc func(error, http.ResponseWriter, *model.Challenge) error

const gogetMeta = `<meta name="go-import" content="%s git %s">`

// Get challenge details.
// The challenge is identified via "id" route parameter, which
// can take the form of following string values:
// * current - currently running challenge, or an error
//             if there isn't one
// * 123 - numerical ID of the challenge
// * challenge-123 - prefixed numerical ID of the challenge.
func Get(cs model.Challenges) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		c, err := findChallenge(cs, ps.ByName("id"))

		if err = responder(r)(err, w, c); err != nil {
			write.Error(w, r, err)
		}
	}
}

// find a challenge given the value if requested ID string
func findChallenge(cs model.Challenges, id string) (*model.Challenge, error) {
	var cid model.ChallengeID

	idx := strings.Replace(id, "challenge-", "", 1)

	if idx == "current" {
		return cs.Current()
	} else if err := cid.Atoid(idx); err == nil {
		return cs.Find(cid)
	} else {
		return nil, err
	}
}

// determine response format by request parameters
func responder(r *http.Request) writerFunc {
	if err := r.ParseForm(); err == nil && r.Form.Get("go-get") == "1" {
		return gogeter
	}
	return jsonifier
}

// render a challenge in a way interpretable by go get tool
func gogeter(err error, w http.ResponseWriter, c *model.Challenge) error {
	if err != nil {
		return err
	}
	if c.Git == "" {
		return model.ErrNoRemote
	}

	s := fmt.Sprintf(gogetMeta, c.Import, c.Git)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte(s))
	return err
}

// render a challenge as json
func jsonifier(err error, w http.ResponseWriter, c *model.Challenge) error {
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(c)
}
