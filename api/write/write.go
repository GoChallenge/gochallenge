package write

import (
	"fmt"
	"net/http"
)

// Error is being reported back to an API client
func Error(w http.ResponseWriter, _ *http.Request, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("%s", err)))
}
