package model

import "fmt"

// Lifecycle of a challenge
type Lifecycle int

// Steps of Go challenge lifecycle
const (
	Unreleased = iota // not released to public yet
	Open              // open and running
	Closed            // done and closed
)

var lifecycleEncoder = map[Lifecycle]([]byte){
	Unreleased: []byte(`"unreleased"`),
	Open:       []byte(`"open"`),
	Closed:     []byte(`"closed"`),
}

var lifecycleDecoder = map[string]Lifecycle{
	`"unreleased"`: Unreleased,
	`"open"`:       Open,
	`"closed"`:     Closed,
}

// MarshalJSON marshals lifecycle into its string-based JSON form
func (l Lifecycle) MarshalJSON() ([]byte, error) {
	if b, ok := lifecycleEncoder[l]; ok {
		return b, nil
	}

	return []byte{}, fmt.Errorf("Unexpected lifecycle value %d", l)
}

// UnmarshalJSON loads lifecycle from JSON form
func (l *Lifecycle) UnmarshalJSON(b []byte) error {
	if nl, ok := lifecycleDecoder[string(b)]; ok {
		*l = nl
		return nil
	}
	return fmt.Errorf("Unknown lifecycle encoding %s", b)
}
