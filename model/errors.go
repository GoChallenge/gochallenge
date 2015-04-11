package model

import "fmt"

// Error type for errors returned by the model package
type Error int

// Errors defined and used by the application
const (
	ErrNotFound Error = iota
	ErrNoRemote
	ErrCryptoFailure
	ErrGithubAPIError
	ErrAuthFailure
	ErrNotImplemented
	ErrDuplicateRecord
	ErrAccessDenied
)

var errmsgs = map[Error]string{
	ErrNotFound:        "Not found",
	ErrNoRemote:        "Challenge does not have git remote",
	ErrCryptoFailure:   "Error in cryptographical operation",
	ErrGithubAPIError:  "Error communicating with Github API",
	ErrAuthFailure:     "Invalid authentication",
	ErrNotImplemented:  "Not implemented",
	ErrDuplicateRecord: "Record already exists",
	ErrAccessDenied:    "Access denied",
}

func (e Error) Error() string {
	if s, ok := errmsgs[e]; ok {
		return s
	}
	return fmt.Sprintf("Error %d", e)
}
