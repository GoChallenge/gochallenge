package model

// Challenges repository interface
type Challenges interface {
	Find(string) (Challenge, error)
	Current() (Challenge, error)
}

// Challenge type describes details of a Go challenge
type Challenge struct {
	ID     string
	Name   string
	Author string
}
