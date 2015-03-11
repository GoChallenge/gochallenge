package model

// Submissions repository interface
type Submissions interface {
	All() []Submissions
	ByID(string) (Submission, error)
}

// Submission type describes details of a submitted solutions for a
// challenge
type Submission struct {
	ID string
}
