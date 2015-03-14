package mock

import (
	"errors"

	"github.com/gochallenge/gochallenge/model"
)

// Submissions repository, mocked out as in-memory array
type Submissions struct {
	ary []*model.Submission
}

// NewSubmissions returns a new initialised struct of challenges
func NewSubmissions() Submissions {
	var ss []*model.Submission
	return Submissions{
		ary: ss,
	}
}

// Add another challenge to the mock repo
func (ss *Submissions) Add(s *model.Submission) error {
	s.ID = len(ss.ary) + 1
	ss.ary = append(ss.ary, s)
	return nil
}

// Find a challenge in the repository by its id
func (ss *Submissions) Find(id int) (model.Submission, error) {
	return model.Submission{}, errors.New("Not implemented")
}

// All submissions received
func (ss *Submissions) All() ([]*model.Submission, error) {
	return ss.ary, nil
}
