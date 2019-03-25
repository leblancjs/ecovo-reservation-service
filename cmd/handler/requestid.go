package handler

import (
	"context"
	"net/http"

	"azure.com/ecovo/reservation-service/cmd/middleware/requestid"
	"github.com/google/uuid"
)

// RequestID extracts the request ID from a request's headers, if it is
// present, and stores it in the request's context.
//
// If no request ID is present in the request's headers, it will be generated.
func RequestID(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), requestid.RequestIDContextKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))

		return nil
	}
}
