package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

type marshalFunc func(model.Challenge) ([]byte, error)

const gogetMeta = `<meta name="go-import" content="%s git %s">`

func (a *app) getChallenge(w http.ResponseWriter,
	r *http.Request, ps httprouter.Params) {
	var (
		js []byte
	)

	c, err := findChallenge(a.cfg.Challenges, ps.ByName("id"))

	if err == nil {
		f := responseFormat(r)
		js, err = f(c)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write(js)
	}
}

func findChallenge(cs model.Challenges, id string) (model.Challenge, error) {
	if id == "current" {
		return cs.Current()
	} else if cid, err := strconv.Atoi(id); err == nil {
		return cs.Find(cid)
	} else {
		return model.Challenge{}, err
	}
}

func responseFormat(r *http.Request) marshalFunc {
	if err := r.ParseForm(); err == nil && r.Form.Get("go-get") == "1" {
		return gogeter
	}
	return jsonifier
}

func gogeter(c model.Challenge) ([]byte, error) {
	if c.Git != "" {
		return []byte(fmt.Sprintf(gogetMeta, c.Import, c.Git)), nil
	}
	return []byte{}, errors.New("challenge does not have git remote")
}

func jsonifier(c model.Challenge) ([]byte, error) {
	return json.Marshal(c)
}
