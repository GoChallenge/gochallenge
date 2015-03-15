package model

// Submissions repository interface
type Submissions interface {
	All() ([]*Submission, error)
	AllForChallenge(*Challenge) ([]*Submission, error)
	Find(int) (Submission, error)
	Add(*Submission) error
}

// User of a challenge
type User struct {
	Name string `json:"name"`
}

// Submission type describes details of a submitted solutions for a
// challenge
type Submission struct {
	ID          string        `json:"id"`
	User        User          `json:"user"`
	ChallengeID int           `json:"challenge_id"`
	Type        Participation `json:"type"`
	Challenge   *Challenge
	Data        []byte
}
