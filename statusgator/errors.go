package statusgator

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Sentinel errors for common failure cases.
var (
	ErrTokenRequired    = errors.New("api token is required")
	ErrEmptyID          = errors.New("id cannot be empty")
	ErrResponseTooLarge = errors.New("response body exceeds maximum size")
	ErrBoardNotFound    = errors.New("board not found")
	ErrMonitorNotFound  = errors.New("monitor not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrUnauthorized     = errors.New("unauthorized: invalid or missing token")
	ErrForbidden        = errors.New("forbidden: insufficient permissions")
)

// APIError represents an error response from the StatusGator API.
type APIError struct {
	StatusCode int      `json:"-"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors,omitempty"`
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("statusgator: %d - %s: %v", e.StatusCode, e.Message, e.Errors)
	}
	return fmt.Sprintf("statusgator: %d - %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error indicates a 404 response.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return errors.Is(err, ErrBoardNotFound) ||
		errors.Is(err, ErrMonitorNotFound) ||
		errors.Is(err, ErrServiceNotFound)
}

// IsUnauthorized returns true if the error indicates a 401 response.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401
	}
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden returns true if the error indicates a 403 response.
func IsForbidden(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 403
	}
	return errors.Is(err, ErrForbidden)
}

func parseAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{StatusCode: statusCode}

	if err := json.Unmarshal(body, apiErr); err != nil {
		apiErr.Message = string(body)
	}

	if apiErr.Message == "" {
		apiErr.Message = "unknown error"
	}

	switch statusCode {
	case 401:
		return fmt.Errorf("%w: %s", ErrUnauthorized, apiErr.Message)
	case 403:
		return fmt.Errorf("%w: %s", ErrForbidden, apiErr.Message)
	case 404:
		return apiErr
	}

	return apiErr
}

// validateID checks if the provided ID is non-empty.
func validateID(id, name string) error {
	if id == "" {
		return fmt.Errorf("%w: %s", ErrEmptyID, name)
	}
	return nil
}
