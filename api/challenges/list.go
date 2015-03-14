package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

// List all available challenges
func List(cs model.Challenges) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		var b []byte

		cx, err := cs.All()

		if err == nil {
			b, err = json.Marshal(cx)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			w.Write(b)
		}
	}
}
