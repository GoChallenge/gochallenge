package boltdb

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gochallenge/gochallenge/model"
)

var bktUsers = []byte("Users")

// Users repository
type Users struct {
	db *bolt.DB
}

// NewUsers returns bolt-backed Users repo
func NewUsers(db *bolt.DB) (Users, error) {
	err := db.Update(initBucket(bktUsers))
	return Users{db}, err
}

// Add another user to the repo
func (us *Users) Add(u *model.User) error {
	return chain(us.db.Update,
		validateNewUser(u),
		store(bktUsers, strconv.Itoa(u.ID), u),
	)
}

// Find returns a user record for the given ID
func (us *Users) Find(id int) (*model.User, error) {
	var u model.User
	return &u, us.db.View(load(bktUsers, strconv.Itoa(id), &u))
}

// FindByGithubID returns a user record with the given Github ID
func (us *Users) FindByGithubID(ghid int) (*model.User, error) {
	return us.findBy(func(u *model.User) bool {
		return u.GithubID == ghid
	})
}

// FindByAPIKey returns a user record with the given API key
func (us *Users) FindByAPIKey(k string) (*model.User, error) {
	return us.findBy(func(u *model.User) bool {
		return u.APIKey == k
	})
}

//
// Low-level database operations
//

func (us *Users) findBy(f func(*model.User) bool) (*model.User, error) {
	var u model.User
	return &u, us.db.View(findUser(us, &u, f))
}

// validates given user record as a new record - e.g. making sure is does
// not conflict with another existing record, etc
func validateNewUser(u *model.User) boltf {
	return func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers)
		k := strconv.Itoa(u.ID)

		if bkt.Get([]byte(k)) != nil {
			return model.ErrDuplicateRecord
		}
		return nil
	}
}

// find the first user that matches given finder function criteria
func findUser(us *Users, u *model.User, f func(*model.User) bool) boltf {
	return func(tx *bolt.Tx) error {
		var err error
		bkt := tx.Bucket(bktUsers).Cursor()

		for k, v := bkt.First(); k != nil && err == nil; k, v = bkt.Next() {
			if err = decode(&v, u); err == nil && f(u) {
				// the matching record is found, stop here
				return nil
			}
		}
		// no matching record was found, if there're no errors either -
		// return ErrNotFound
		if err == nil {
			err = model.ErrNotFound
		}
		return err
	}
}
