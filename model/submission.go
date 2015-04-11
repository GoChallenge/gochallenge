package model

import (
	"encoding/json"
	"time"
)

// Submissions repository interface
type Submissions interface {
	All() ([]*Submission, error)
	AllForChallenge(*Challenge) ([]*Submission, error)
	Find(string) (*Submission, error)
	Add(*Submission) error
}

// Submission type describes details of a submitted solutions for a
// challenge
type Submission struct {
	ID        string        `json:"id"`
	User      *User         `json:"user"`
	Type      Participation `json:"type"`
	Challenge *Challenge    `json:"-"`
	Data      *[]byte       `json:"-"`
	Created   time.Time     `json:"created"`
}

// unexported type to use as a basis for JSON representation
// of a submission, mostly to replace associated objects
// with their IDs
type submissionEx struct {
	ID          string        `json:"id"`
	UserID      UserID        `json:"user_id"`
	ChallengeID ChallengeID   `json:"challenge_id"`
	Type        Participation `json:"type"`
	Created     time.Time     `json:"created"`
}

// MarshalJSON exports submission data, substituting associations
// with their IDs
func (s Submission) MarshalJSON() ([]byte, error) {
	se := &submissionEx{
		ID:      s.ID,
		Type:    s.Type,
		Created: s.Created,
	}
	if s.Challenge != nil {
		se.ChallengeID = s.Challenge.ID
	}
	if s.User != nil {
		se.UserID = s.User.ID
	}

	return json.Marshal(se)
}

// Unmarshal imports submission data, hydrating associated objects
// based on their ID values received
func (s *Submission) Unmarshal(b []byte, cs Challenges, us Users) error {
	var err error

	var se submissionEx
	if err = json.Unmarshal(b, &se); err != nil {
		return err
	}

	s.ID = se.ID
	s.Type = se.Type
	s.Created = se.Created

	if se.ChallengeID != 0 {
		s.Challenge, err = cs.Find(se.ChallengeID)
	}
	if err == nil && se.UserID != 0 {
		s.User, err = us.Find(se.UserID)
	}

	return err
}
