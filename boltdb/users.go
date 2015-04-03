package boltdb

import (
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
	return Users{db: db}, err
}

// Save a user into the repo. If this is a new user record, and it doesn't
// have its ID specified yet - it will be set to the next available value
func (us *Users) Save(u *model.User) error {
	return chain(us.db.Update,
		prefillID(u),
		store(bktUsers, &u.ID, u),
	)
}

// Find returns a user record for the given ID
func (us *Users) Find(id model.UserID) (*model.User, error) {
	var u model.User
	return &u, us.db.View(load(bktUsers, id, &u))
}

// FindByGithubID returns a user record with the given Github ID
func (us *Users) FindByGithubID(ghid int) (*model.User, error) {
	return us.findBy(func(u interface{}) bool {
		return u.(*model.User).GithubID == ghid
	})
}

// FindByAPIKey returns a user record with the given API key
func (us *Users) FindByAPIKey(k string) (*model.User, error) {
	return us.findBy(func(u interface{}) bool {
		return u.(*model.User).APIKey == k
	})
}

//
// Low-level database operations
//

func (us *Users) findBy(f func(interface{}) bool) (*model.User, error) {
	var u model.User
	return &u, us.db.View(first(bktUsers, f, &u))
}

// prefills user's ID with the next available unique value. If user already
// has its ID set - does nothing.
func prefillID(u *model.User) boltf {
	return func(tx *bolt.Tx) error {
		if u.ID != 0 {
			return nil
		}
		var id model.UserID
		if err := maxKey(tx, bktUsers, &id); err != nil && err.Error() != "EOF" {
			return err
		}
		u.ID = id + 1
		return nil
	}
}
