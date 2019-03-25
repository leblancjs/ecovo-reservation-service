package requestid

import (
	"context"
	"fmt"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	// RequestIDContextKey represents the key used to store and retrieve the
	// request ID from the request's context.
	RequestIDContextKey = contextKey("X-Request-ID")
)

// FromContext extracts the request ID from a request's context.
func FromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("requestid: context is nil")
	}

	requestID, ok := ctx.Value(RequestIDContextKey).(string)
	if !ok {
		return "", fmt.Errorf("requestid: %s not found in context", RequestIDContextKey)
	}

	return requestID, nil
}
