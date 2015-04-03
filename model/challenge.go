package model

import "time"

// Challenges repository interface
type Challenges interface {
	Add(*Challenge) error
	Find(ChallengeID) (*Challenge, error)
	Current() (*Challenge, error)
	All() ([]*Challenge, error)
}

// Author of a challenge
type Author struct {
	Name string `json:"name"`
}

// ChallengeID type
type ChallengeID int32

// Atoid convert string value into ChallengeID
func (uid *ChallengeID) Atoid(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	*uid = ChallengeID(n)
	return nil
}

// Challenge type describes details of a Go challenge
type Challenge struct {
	ID     ChallengeID `json:"id"`
	Name   string      `json:"name"`
	Author Author      `json:"author"`
	URL    string      `json:"url"`
	Import string      `json:"import"`
	Git    string      `json:"-"`
	Status Lifecycle   `json:"status"`
	Start  time.Time   `json:"start"`
	End    time.Time   `json:"end"`
}

// Current status of the challenge
func (ch Challenge) Current() bool {
	now := time.Now()
	return ch.Start.Before(now) && ch.End.After(now)
}
