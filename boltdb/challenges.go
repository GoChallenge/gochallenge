package boltdb

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gochallenge/gochallenge/model"
)

// Open Bolt database at the given file name
func Open(file string) (*bolt.DB, error) {
	opt := &bolt.Options{
		Timeout: 1 * time.Second,
	}
	return bolt.Open(file, 0600, opt)
}

const bktChallenges = "Chals"

// Challenges repository
type Challenges struct {
	db *bolt.DB
}

// NewChallenges returns a new initialised struct of challenges
func NewChallenges(db *bolt.DB) (Challenges, error) {
	err := db.Update(initChallengeBucket)
	return Challenges{db}, err
}

// Add another challenge to the repo
func (cs *Challenges) Add(c *model.Challenge) error {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(c); err != nil {
		return err
	}
	err := cs.db.Update(putChallenge(c, b.Bytes()))
	return err
}

// Find a challenge in the repository by its id
func (cs *Challenges) Find(id int) (*model.Challenge, error) {
	var (
		b    *[]byte
		chal model.Challenge
	)

	if err := cs.db.View(getChallenge(id, &b)); err != nil {
		return nil, err
	}
	// bolt returns an empty result for unknown key lookup,
	// test it here and return ErrNotFound in this case
	err := decodeChal(b, &chal)
	return &chal, err
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
	f := func(c *model.Challenge) bool {
		return c.Current()
	}
	err := cs.db.View(findChallenge(cs, &chal, f))
	return &chal, err
}

//
// Low-level database operations
//

func decodeChal(b *[]byte, chal *model.Challenge) error {
	// bolt returns an empty result for unknown key lookup,
	// test it here and return ErrNotFound in this case
	if len(*b) == 0 {
		return model.ErrNotFound
	}
	err := gob.NewDecoder(bytes.NewReader(*b)).Decode(chal)
	return err

}

// initialises bolt bucket for challenges
func initChallengeBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte(bktChallenges))
	return err
}

// put encoded challenge data into the database
func putChallenge(c *model.Challenge, b []byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktChallenges))

		k := strconv.Itoa(c.ID)
		return bkt.Put([]byte(k), b)
	}
}

// get encoded challenge data from the database
func getChallenge(id int, b **[]byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktChallenges))

		k := strconv.Itoa(id)
		v := bkt.Get([]byte(k))
		*b = &v
		return nil
	}
}

// find the first challenge that matches given finder function criteria
func findChallenge(cs *Challenges, chal *model.Challenge,
	f func(*model.Challenge) bool) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		var err error
		bkc := tx.Bucket([]byte(bktChallenges)).Cursor()

		for k, v := bkc.First(); k != nil && err == nil; k, v = bkc.Next() {
			if err = decodeChal(&v, chal); err == nil && chal.Current() {
				// the matching challenge is found, stop here
				return nil
			}
		}
		// no matching challenge was found, if there're no errors either -
		// return ErrNotFound
		if err == nil {
			err = model.ErrNotFound
		}
		return err
	}
}

// get all challenges from the database
func getChallenges(chals *[]*model.Challenge) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktChallenges))

		err := bkt.ForEach(func(_, v []byte) error {
			chal := model.Challenge{}
			if err := decodeChal(&v, &chal); err != nil {
				return err
			}
			*chals = append(*chals, &chal)
			return nil
		})
		return err
	}
}
