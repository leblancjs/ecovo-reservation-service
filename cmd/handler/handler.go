package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"azure.com/ecovo/reservation-service/cmd/middleware/requestid"
)

// A Handler represents a handler that can return an error.
type Handler func(http.ResponseWriter, *http.Request) error

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerErr := WrapError(handler(w, r))
	if handlerErr != nil {
		requestID, _ := requestid.FromContext(r.Context())

		log.Printf("[Request ID=%s] error: %s", requestID, handlerErr)

		type errorResponse struct {
			*Error
			RequestID string `json:"requestId"`
		}

		w.WriteHeader(handlerErr.Code)
		err := json.NewEncoder(w).Encode(errorResponse{
			handlerErr,
			requestID,
		})
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf(
					`Oh boy!
					Something went wrong while we were handling an error.
					That's embarassing.
					Please contact your system administrator.
					Request ID: %s`,
					requestID,
				),
				http.StatusInternalServerError,
			)
		}
	}
}
