package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"azure.com/ecovo/reservation-service/cmd/middleware/auth"
)

// Auth validates a request's authorization header using the given validator
// to ensure that the user is authorized to access an endpoint and extracts the
// authenticated user's information.
//
// The authenticated user's information placed in the request's context and can
// be accessed by using the auth.FromContext utility function.
func Auth(validators map[string]auth.Validator, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		header := r.Header.Get("Authorization")

		authType, authCredentials, err := parseHeader(header)
		if err != nil {
			return auth.UnauthorizedError{Msg: fmt.Sprintf("auth: %s", err)}
		}

		authType = strings.ToLower(authType)

		v, ok := validators[authType]
		if !ok {
			return auth.UnauthorizedError{Msg: fmt.Sprintf("auth: no validator found for %s", authType)}
		}

		userInfo, err := v.Validate(authCredentials)
		if err != nil {
			return err
		}

		ctx := context.WithValue(r.Context(), auth.UserInfoContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))

		return nil
	}
}

func parseHeader(header string) (string, string, error) {
	headerParts := strings.Split(header, " ")
	if len(headerParts) < 2 {
		return "", "", fmt.Errorf("failed to parse authorization header")
	}

	return headerParts[0], headerParts[1], nil
}
