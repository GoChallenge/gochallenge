package mock

import (
	"errors"

	"github.com/gochallenge/gochallenge/model"
)

// Users represents users collection, mocked out as in-memory map
type Users struct {
	index       map[int]*model.User
	indexAPIKey map[string]*model.User
}

// NewUsers returns a new initialised users collection.
func NewUsers() Users {
	return Users{
		index:       make(map[int]*model.User),
		indexAPIKey: make(map[string]*model.User),
	}
}

// Add user to the mock users.
func (us *Users) Add(u *model.User) error {
	if _, ok := us.index[u.ID]; ok {
		return errors.New("Users already exists")
	}
	us.index[u.ID] = u
	us.indexAPIKey[u.APIKey] = u

	return nil
}

func (us *Users) Update(u *model.User) error {
	if _, ok := us.index[u.ID]; !ok {
		return errors.New("Unknown user ID")
	}
	us.index[u.ID] = u
	us.indexAPIKey[u.APIKey] = u

	return nil
}

// Find a user in the collection by its id.
func (us *Users) FindByID(id int) (*model.User, error) {
	var (
		u  *model.User
		ok bool
	)

	if u, ok = us.index[id]; !ok {
		return nil, model.ErrNotFound
	}
	return u, nil
}

// Find a user in the collection by its API Key.
func (us *Users) FindByAPIKey(key string) (*model.User, error) {
	var (
		u  *model.User
		ok bool
	)

	if u, ok = us.indexAPIKey[key]; !ok {
		return nil, model.ErrNotFound
	}
	return u, nil
}
