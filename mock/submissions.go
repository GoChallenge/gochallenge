package mock

import (
	"strconv"

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
	s.ID = strconv.Itoa(len(ss.ary) + 1)
	ss.ary = append(ss.ary, s)
	return nil
}

// Find a challenge in the repository by its id
func (ss *Submissions) Find(id string) (*model.Submission, error) {
	for _, s := range ss.ary {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, model.ErrNotFound
}

// All submissions received
func (ss *Submissions) All() ([]*model.Submission, error) {
	return ss.ary, nil
}

// AllForChallenge return submissions received for the given challenge
func (ss *Submissions) AllForChallenge(c *model.Challenge) ([]*model.Submission, error) {
	var sx []*model.Submission
	sx = make([]*model.Submission, 0)

	for _, s := range ss.ary {
		if s.Challenge.ID == c.ID {
			sx = append(sx, s)
		}
	}
	return sx, nil
}
