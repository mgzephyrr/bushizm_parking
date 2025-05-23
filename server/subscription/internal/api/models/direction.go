package models

import (
	"encoding/json"
	"fmt"
)

type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)

func (d *Direction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch Direction(s) {
	case DirectionIn, DirectionOut:
		*d = Direction(s)
		return nil
	default:
		return fmt.Errorf("invalid direction value: %s", s)
	}
}

func (d Direction) MarshalJSON() ([]byte, error) {
	if !d.IsValid() {
		return nil, fmt.Errorf("invalid direction value: %s", d)
	}
	return json.Marshal(string(d))
}

func (d Direction) IsValid() bool {
	return d == DirectionIn || d == DirectionOut
}

func (d Direction) String() string {
	return string(d)
}
