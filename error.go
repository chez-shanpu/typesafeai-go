// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package typesafeai

import (
	"fmt"
)

// ValidationError represents a locally detected invalid configuration or input.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("typesafeai: validation error: %s", e.Message)
	}
	return fmt.Sprintf("typesafeai: validation error on %s: %s", e.Field, e.Message)
}

// APIError represents an HTTP-level failure returned by the API.
type APIError struct {
	StatusCode int
	RequestID  string
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("typesafeai: api error: status %d, request id %s, body %s", e.StatusCode, e.RequestID, string(e.Body))
}

// ConnectionError represents a failure to communicate with the API, such as a
// network error, a canceled context, or an expired deadline.
type ConnectionError struct {
	Err error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("typesafeai: connection error: %s", e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}
