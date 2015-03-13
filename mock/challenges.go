package mock

import (
	"errors"

	"github.com/gochallenge/gochallenge/model"
)

// CurrentID is an ID of a challenge that mock considers to be
// the current one
const CurrentID = 0

// Challenges repository, mocked out as in-memory map
type Challenges map[int]model.Challenge

// Add another challenge to the mock repo
func (cs Challenges) Add(c model.Challenge) {
	cs[c.ID] = c
}

// Find a challenge in the repository by its id
func (cs Challenges) Find(id int) (model.Challenge, error) {
	var (
		c   model.Challenge
		ok  bool
		err error
	)

	if c, ok = cs[id]; !ok {
		err = errors.New("Unknown challenge ID")
	}

	return c, err
}

// Current challenge, mocked to return challenge with ID "0"
func (cs Challenges) Current() (model.Challenge, error) {
	return cs.Find(CurrentID)
}
