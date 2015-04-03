package boltdb

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gochallenge/gochallenge/model"
)

var bktChallenges = []byte("Chals")

// Challenges repository
type Challenges struct {
	db *bolt.DB
}

// NewChallenges returns a new initialised struct of challenges
func NewChallenges(db *bolt.DB) (Challenges, error) {
	err := db.Update(initBucket(bktChallenges))
	return Challenges{db}, err
}

// Add another challenge to the repo
func (cs *Challenges) Add(c *model.Challenge) error {
	return cs.db.Update(store(bktChallenges, strconv.Itoa(c.ID), c))
}

// Find a challenge in the repository by its id
func (cs *Challenges) Find(id int) (*model.Challenge, error) {
	var chal model.Challenge
	return &chal, cs.db.View(load(bktChallenges, strconv.Itoa(id), &chal))
}

// All challenges currently available
func (cs *Challenges) All() ([]*model.Challenge, error) {
	var chals []*model.Challenge

	err := cs.db.View(getChallenges(&chals))
	return chals, err
}

// Current challenge, according to the rules defined on challenge
// model itself
func (cs *Challenges) Current() (*model.Challenge, error) {
	var chal model.Challenge

	// TODO: this can be optimised by keeping a pointer to
	// the current challenge in the database. When retrieved,
	// the challenge can be re-verified as current, and only
	// if it is out of date - full re-scan be triggered
	f := func(c interface{}) bool {
		return c.(*model.Challenge).Current()
	}
	return &chal, cs.db.View(first(bktChallenges, f, &chal))
}

//
// Low-level database operations
//

// get all challenges from the database
func getChallenges(chals *[]*model.Challenge) boltf {
	return func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktChallenges)

		err := bkt.ForEach(func(_, v []byte) error {
			chal := model.Challenge{}
			if err := decode(&v, &chal); err != nil {
				return err
			}
			*chals = append(*chals, &chal)
			return nil
		})
		return err
	}
}
