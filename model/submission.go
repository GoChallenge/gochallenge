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
	ID          string        `json:"id"`
	User        *User         `json:"user"`
	ChallengeID int           `json:"challenge_id"`
	Type        Participation `json:"type"`
	Challenge   *Challenge    `json:"-"`
	Data        *[]byte       `json:"-"`
	Created     time.Time     `json:"created"`
}

// type aliases to aid in custom marshalling of Submission structs
type submissionExport Submission
type submissionImport Submission

// MarshalJSON exports submission data, populating dynamic
// fields
func (s Submission) MarshalJSON() ([]byte, error) {
	if s.Challenge != nil {
		s.ChallengeID = s.Challenge.ID
	}
	return json.Marshal(submissionExport(s))
}

// Hydrate submission model, mapping numeric IDs (e.g. ChallengeID)
// to their full objects
func (s *Submission) Hydrate(cs Challenges) error {
	c, err := cs.Find(s.ChallengeID)
	if err != nil {
		return err
	}
	s.Challenge = c
	return nil
}
