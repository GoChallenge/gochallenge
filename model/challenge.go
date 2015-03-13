package model

// Challenges repository interface
type Challenges interface {
	Find(int) (Challenge, error)
	Current() (Challenge, error)
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
	Status Lifecycle `json:"status"`
}
