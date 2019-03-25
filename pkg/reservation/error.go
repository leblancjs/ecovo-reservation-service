package reservation

// A NotFoundError is an error that represents that no reservation was found.
type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	return e.msg
}

// A AlreadyExistsError is an error that represents that a reservation already exists
// with a given unique identifier.
type AlreadyExistsError struct {
	msg string
}

func (e AlreadyExistsError) Error() string {
	return e.msg
}
