package model

import "time"

// Challenges repository interface
type Challenges interface {
	Add(*Challenge) error
	Find(int) (*Challenge, error)
	Current() (*Challenge, error)
	All() ([]*Challenge, error)
}

// Author of a challenge
type Author struct {
	Name string `json:"name"`
}

// Challenge type describes details of a Go challenge
type Challenge struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Author Author    `json:"author"`
	URL    string    `json:"url"`
	Import string    `json:"import"`
	Git    string    `json:"-"`
	Status Lifecycle `json:"status"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
