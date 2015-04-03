package model

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"strconv"
)

// Users collection interface
type Users interface {
	Add(*User) error
	Find(UserID) (*User, error)
	FindByGithubID(int) (*User, error)
	FindByAPIKey(string) (*User, error)
}

// NewUser generates and populates with basic details new user record
func NewUser() (*User, error) {
	u := &User{}
	err := u.ResetToken()
	return u, err
}

// UserID type
type UserID int32

// Atoid convert string value into UserID
func (uid *UserID) Atoid(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	*uid = UserID(n)
	return nil
}

// User of a challenge
type User struct {
	ID          UserID `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	AvatarURL   string `json:"avatar_url"`
	GithubID    int    `json:"-"`
	GithubURL   string `json:"github_url"`
	GithubLogin string `json:"github_login"`
	APIKey      string `json:"-"`
}

// ResetToken rewrites API key on the user record
func (u *User) ResetToken() error {
	var err error

	u.APIKey, err = generateToken()
	return err
}

func generateToken() (string, error) {
	const length = sha1.BlockSize

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", ErrCryptoFailure
	}

	s1 := sha1.Sum(b)
	return fmt.Sprintf("%x", s1), nil
}
