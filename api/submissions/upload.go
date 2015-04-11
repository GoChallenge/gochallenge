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
	"strings"
	"time"

	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// Post new sumission
func Post(cs model.Challenges, ss model.Submissions, us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var s model.Submission

		u, err := auth.User(r, us)
		s.Challenge, err = findChallenge(err, cs, ps.ByName("id"))
		err = readSubmission(err, &s, r)
		err = storeSubmission(err, ss, &s, u)
		err = writeSubmission(err, w, s)

		write.Error(w, r, err)
	}
}

// find a challenge given the value of requested ID string
func findChallenge(err error, cs model.Challenges,
	id string) (*model.Challenge, error) {
	if err != nil {
		return nil, err
	}

	var cid model.ChallengeID
	err = cid.Atoid(id)
	if err != nil {
		return nil, err
	}
	return cs.Find(cid)
}

func readSubmission(err error, s *model.Submission, r *http.Request) error {
	var bnd string

	if err != nil {
		return err
	}

	if bnd, err = boundary(r); err != nil {
		return err
	}

	mr := multipart.NewReader(r.Body, bnd)
	for ; err == nil; err = parsePart(mr, s) {
	}

	// io.EOF means we completed the parsing, this is not a reportable
	// error
	if err == io.EOF {
		err = nil
	}

	return err
}

func boundary(r *http.Request) (string, error) {
	ct := r.Header.Get("Content-Type")
	mt, args, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", err
	}

	bnd := args["boundary"]
	if !strings.HasPrefix(mt, "multipart/") || bnd == "" {
		return "", fmt.Errorf("invalid content type %s", ct)
	}
	return bnd, err
}

// Parse the next part of multipart message, and handle its content
// depending on this part's content type
func parsePart(mr *multipart.Reader, s *model.Submission) error {
	var (
		p   *multipart.Part
		mt  string
		err error
	)
	if p, err = mr.NextPart(); err != nil {
		return err
	}

	mt, _, err = mime.ParseMediaType(p.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	switch mt {
	case "application/json":
		err = parseJSON(s, p)
	case "application/zip":
		err = parseZip(s, p)
	}

	return err
}

// JSON part is submission's metadata
func parseJSON(s *model.Submission, p *multipart.Part) error {
	b, err := ioutil.ReadAll(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// ZIP part is submission's binary archive
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
	s.Data = &b

	return nil
}

func writeSubmission(err error, w http.ResponseWriter,
	s model.Submission) error {

	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(s)
}

func storeSubmission(err error, ss model.Submissions, s *model.Submission, u *model.User) error {
	if err != nil {
		return err
	}
	s.User = u
	s.Created = time.Now().UTC()
	return ss.Add(s)
}
