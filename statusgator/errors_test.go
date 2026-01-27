package statusgator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name: "with message only",
			err: &APIError{
				StatusCode: 404,
				Message:    "Board not found",
			},
			expected: "statusgator: 404 - Board not found",
		},
		{
			name: "with errors array",
			err: &APIError{
				StatusCode: 422,
				Message:    "Validation failed",
				Errors:     []string{"name is required", "url is invalid"},
			},
			expected: "statusgator: 422 - Validation failed: [name is required url is invalid]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "APIError with 404",
			err:      &APIError{StatusCode: 404, Message: "Not found"},
			expected: true,
		},
		{
			name:     "APIError with 500",
			err:      &APIError{StatusCode: 500, Message: "Server error"},
			expected: false,
		},
		{
			name:     "ErrBoardNotFound",
			err:      ErrBoardNotFound,
			expected: true,
		},
		{
			name:     "ErrMonitorNotFound",
			err:      ErrMonitorNotFound,
			expected: true,
		},
		{
			name:     "ErrServiceNotFound",
			err:      ErrServiceNotFound,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("random error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsNotFound(tt.err))
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "APIError with 401",
			err:      &APIError{StatusCode: 401, Message: "Unauthorized"},
			expected: true,
		},
		{
			name:     "APIError with 403",
			err:      &APIError{StatusCode: 403, Message: "Forbidden"},
			expected: false,
		},
		{
			name:     "ErrUnauthorized",
			err:      ErrUnauthorized,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("random error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsUnauthorized(tt.err))
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "APIError with 403",
			err:      &APIError{StatusCode: 403, Message: "Forbidden"},
			expected: true,
		},
		{
			name:     "ErrForbidden",
			err:      ErrForbidden,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("random error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsForbidden(tt.err))
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		idName  string
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      "abc123",
			idName:  "boardID",
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			idName:  "boardID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id, tt.idName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrEmptyID)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
