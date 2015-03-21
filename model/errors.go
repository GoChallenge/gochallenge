package model

import "fmt"

// Error type for errors returned by the model package
type Error int

// Errors defined and used by the application
const (
	ErrNotFound Error = iota
	ErrNoRemote
)

var errmsgs = map[Error]string{
	ErrNotFound: "Not found",
	ErrNoRemote: "Challenge does not have git remote",
}

func (e Error) Error() string {
	if s, ok := errmsgs[e]; ok {
		return s
	}
	return fmt.Sprintf("Error %d", e)
}
