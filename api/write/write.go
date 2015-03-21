package write

import (
	"fmt"
	"net/http"

	"github.com/gochallenge/gochallenge/model"
)

// Error is being reported back to an API client
func Error(w http.ResponseWriter, _ *http.Request, err error) {
	if err == nil {
		return
	}

	switch err.(type) {
	case model.Error:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
	}
}
