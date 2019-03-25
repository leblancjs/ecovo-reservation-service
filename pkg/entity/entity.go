package entity

import "strings"

// An ID is an entity's unique identifier.
type ID string

// NilID is the zero value for an ID.
var NilID ID

// NewIDFromHex creates a new unique identifier from a hex string.
func NewIDFromHex(hex string) ID {
	return ID(hex)
}

// Hex returns the hex encoding of the unique identifier as a string.
func (id ID) Hex() string {
	return string(id)
}

// IsZero compares the ID with the zero value and returns whether or not it is
// nil.
func (id ID) IsZero() bool {
	return strings.Compare(string(id), string(NilID)) == 0
}
