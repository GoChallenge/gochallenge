package write

import (
	"fmt"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
)

var errcodes = map[model.Error]int{
	model.ErrNotFound:    404,
	model.ErrAuthFailure: 401,
}

// Error is being reported back to an API client
func Error(w http.ResponseWriter, _ *http.Request, err error) {
	if err == nil {
		return
	}

	switch err.(type) {
	case model.Error:
		code, ok := errcodes[err.(model.Error)]
		if !ok {
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	}
}
