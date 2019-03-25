package trip

// A RequestError is an error caused by a non successful API call to trip-service
type RequestError struct {
	msg string
}

func (e RequestError) Error() string {
	return e.msg
}
