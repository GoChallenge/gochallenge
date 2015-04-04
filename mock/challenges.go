package mock

import "github.com/gochallenge/gochallenge/model"

// CurrentID is an ID of a challenge that mock considers to be
// the current one
const CurrentID = 100001

// Challenges repository, mocked out as in-memory map
type Challenges struct {
	index  map[model.ChallengeID]*model.Challenge
	ary    []*model.Challenge
	lastID model.ChallengeID
}

// NewChallenges returns a new initialised struct of challenges
func NewChallenges() Challenges {
	return Challenges{
		index: make(map[model.ChallengeID]*model.Challenge),
		ary:   make([]*model.Challenge, 0),
	}
}

// Save a challenge into the mock repo
func (cs *Challenges) Save(c *model.Challenge) error {
	if c.ID == 0 {
		cs.lastID++
		c.ID = cs.lastID
	} else if c.ID > cs.lastID {
		cs.lastID = c.ID
	}

	cs.ary = append(cs.ary, c)
	cs.index[c.ID] = c
	return nil
}

// Find a challenge in the repository by its id
func (cs *Challenges) Find(id model.ChallengeID) (*model.Challenge, error) {
	var (
		c  *model.Challenge
		ok bool
	)

	if c, ok = cs.index[id]; !ok {
		return nil, model.ErrNotFound
	}

	return c, nil
}

// All challenges currently available
func (cs *Challenges) All() ([]*model.Challenge, error) {
	return cs.ary, nil
}

// Current challenge, mocked to return challenge with ID "0"
func (cs *Challenges) Current() (*model.Challenge, error) {
	return cs.Find(CurrentID)
}
