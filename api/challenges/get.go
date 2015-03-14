package challenges

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

type marshalFunc func(model.Challenge) ([]byte, error)

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
		var b []byte

		c, err := findChallenge(cs, ps.ByName("id"))

		if err == nil {
			f := responseFormat(r)
			b, err = f(c)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			w.Write(b)
		}
	}
}

// find a challenge given the value if requested ID string
func findChallenge(cs model.Challenges, id string) (model.Challenge, error) {
	idx := strings.Replace(id, "challenge-", "", 1)

	if idx == "current" {
		return cs.Current()
	} else if cid, err := strconv.Atoi(idx); err == nil {
		return cs.Find(cid)
	} else {
		return model.Challenge{}, err
	}
}

// determine response format by request parameters
func responseFormat(r *http.Request) marshalFunc {
	if err := r.ParseForm(); err == nil && r.Form.Get("go-get") == "1" {
		return gogeter
	}
	return jsonifier
}

// render a challenge in a way interpretable by go get tool
func gogeter(c model.Challenge) ([]byte, error) {
	if c.Git != "" {
		return []byte(fmt.Sprintf(gogetMeta, c.Import, c.Git)), nil
	}
	return []byte{}, errors.New("challenge does not have git remote")
}

// render a challenge as json
func jsonifier(c model.Challenge) ([]byte, error) {
	return json.Marshal(c)
}
