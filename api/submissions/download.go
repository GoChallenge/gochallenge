package submissions

import (
	"net/http"
	"strconv"

	"github.com/gochallenge/gochallenge/api/write"
	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// Download archived code of a submission
func Download(ss model.Submissions) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		s, err := ss.Find(ps.ByName("id"))

		if err != nil {
			write.Error(w, r, err)
			return
		}

		w.Header().Add("Content-Type", "application/zip")
		w.Header().Add("Content-Disposition", "attachment; filename=code.zip")
		w.Header().Add("Content-Length", strconv.Itoa(len(s.Data)))
		w.Write(s.Data)
	}
}
