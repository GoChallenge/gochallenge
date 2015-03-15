package model

import "fmt"

// Participation level for a submission
type Participation int

// Possible levels of challenge participation
const (
	lvlNormal = iota
	lvlBonus
	lvlFun
	lvlAnonymous
)

var participationEncoder = map[Participation]([]byte){
	lvlNormal:    []byte(`"normal"`),
	lvlBonus:     []byte(`"bonus"`),
	lvlFun:       []byte(`"fun"`),
	lvlAnonymous: []byte(`"anonymous"`),
}

var participationDecoder = map[string]Participation{
	`"normal"`:    lvlNormal,
	`"bonus"`:     lvlBonus,
	`"fun"`:       lvlFun,
	`"anonymous"`: lvlAnonymous,
}

// MarshalJSON marshals participation into its string-based JSON form
func (l Participation) MarshalJSON() ([]byte, error) {
	if b, ok := participationEncoder[l]; ok {
		return b, nil
	}

	return []byte{}, fmt.Errorf("Unexpected participation value %d", l)
}

// UnmarshalJSON loads participation from JSON form
func (l *Participation) UnmarshalJSON(b []byte) error {
	if nl, ok := participationDecoder[string(b)]; ok {
		*l = nl
		return nil
	}
	return fmt.Errorf("Unknown participation encoding %s", b)
}
