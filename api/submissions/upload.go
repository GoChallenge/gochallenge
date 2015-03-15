package submissions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("%s", err)))
}

func boundary(r *http.Request) (string, error) {
	var bnd string

	ct := r.Header.Get("Content-Type")
	mt, args, err := mime.ParseMediaType(ct)
	bnd = args["boundary"]

	if err == nil && (!strings.HasPrefix(mt, "multipart/") || bnd == "") {
		err = fmt.Errorf("invalid content type %s", ct)
	}
	return bnd, err
}

func readSubmission(s *model.Submission, r *http.Request) error {
	var p *multipart.Part

	bnd, err := boundary(r)
	if err != nil {
		return err
	}

	mr := multipart.NewReader(r.Body, bnd)
	for p, err = mr.NextPart(); err == nil; p, err = mr.NextPart() {
		mt, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
		if err == nil {
			switch mt {
			case "application/json":
				err = parseJSON(s, p)
			case "application/zip":
				err = parseZip(s, p)
			}
		}
		if err != nil {
			break
		}
	}

	if err == io.EOF {
		return nil
	}
	return err
}

func parseJSON(s *model.Submission, p *multipart.Part) error {
	b, err := ioutil.ReadAll(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

func parseZip(s *model.Submission, p *multipart.Part) error {
	var (
		b   []byte
		err error
	)

	if p.Header.Get("Content-Transfer-Encoding") == "base64" {
		// decode base64
		dc := base64.NewDecoder(base64.StdEncoding, p)
		b, err = ioutil.ReadAll(dc)
	} else {
		// default to binary content
		b, err = ioutil.ReadAll(p)
	}
	if err != nil {
		return err
	}
	s.Data = b

	return nil
}

// Post new sumission
func Post(cs model.Challenges, ss model.Submissions) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var (
			c   model.Challenge
			s   model.Submission
			err error
			b   []byte
		)
		c, err = findChallenge(cs, ps.ByName("id"))

		if err == nil {
			err = readSubmission(&s, r)
		}

		if err == nil {
			s.Challenge = &c
			err = ss.Add(&s)
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
