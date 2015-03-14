package mock

import (
	"errors"

	"github.com/gochallenge/gochallenge/model"
)

// CurrentID is an ID of a challenge that mock considers to be
// the current one
const CurrentID = 1

// Challenges repository, mocked out as in-memory map
type Challenges struct {
	index map[int]*model.Challenge
	ary   []model.Challenge
}

// NewChallenges returns a new initialised struct of challenges
func NewChallenges() Challenges {
	return Challenges{
		index: make(map[int]*model.Challenge),
		ary:   make([]model.Challenge, 0),
	}
}

// Add another challenge to the mock repo
func (cs *Challenges) Add(c model.Challenge) {
	cs.ary = append(cs.ary, c)
	cs.index[c.ID] = &c
}

// Find a challenge in the repository by its id
func (cs *Challenges) Find(id int) (model.Challenge, error) {
	var (
		c  *model.Challenge
		ok bool
	)

	if c, ok = cs.index[id]; !ok {
		return model.Challenge{},
			errors.New("Unknown challenge ID")
	}

	return *c, nil
}

// All challenges currently available
func (cs *Challenges) All() ([]model.Challenge, error) {
	return cs.ary, nil
}

// Current challenge, mocked to return challenge with ID "0"
func (cs *Challenges) Current() (model.Challenge, error) {
	return cs.Find(CurrentID)
}
