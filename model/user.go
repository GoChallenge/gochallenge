package model

import (
	"crypto/rand"
	"fmt"
	"io"
)

// Users collection interface
type Users interface {
	Add(*User) error
	Update(*User) error
	FindByID(int) (*User, error)
	FindByAPIKey(string) (*User, error)
}

// User of a challenge
type User struct {
	ID        int    `json:"-"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url"`
	APIKey    string `json:"-"`
}

// GitHubUser represents user of GitHub.
type GitHubUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (gu *GitHubUser) ToUser() *User {
	return &User{
		ID:        gu.ID,
		Name:      gu.Name,
		Email:     gu.Email,
		AvatarURL: gu.AvatarURL,
		APIKey:    fmt.Sprintf("%d-%s", gu.ID, generateToken()),
	}
}

var (
	chars  = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	length = 40
)

func generateToken() string {
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		io.ReadFull(rand.Reader, r)
		for _, c := range r {
			if c >= maxrb {
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
