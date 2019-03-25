package entity

// A ValidationError is an error that occurs when validating an entity fails.
// This can happen when, for example, an entity is missing a required field a
// value is incorrect or out of bounds.
type ValidationError struct {
	msg string
}

func (e ValidationError) Error() string {
	return e.msg
}
