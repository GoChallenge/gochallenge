package submissions

import (
	"net/http"
	"strconv"

	"github.com/gochallenge/gochallenge/api/auth"
	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// Download archived code of a submission
func Download(ss model.Submissions, us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		u, err := auth.User(r, us)
		s, err := findSubmission(err, ss, ps.ByName("id"))
		err = ensureAccess(err, u, s)

		if err != nil {
			write.Error(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=code.zip")
		w.Header().Set("Content-Length", strconv.Itoa(len(*s.Data)))
		w.Write(*s.Data)
	}
}

func findSubmission(err error, ss model.Submissions,
	sid string) (*model.Submission, error) {
	if err != nil {
		return nil, err
	}
	return ss.Find(sid)
}

func ensureAccess(err error, u *model.User, s *model.Submission) error {
	if err != nil {
		return err
	}
	if !readable(u, s) {
		return model.ErrAccessDenied
	}
	return nil
}
