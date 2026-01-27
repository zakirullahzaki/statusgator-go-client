package statusgator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_HasNextPage(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected bool
	}{
		{
			name:     "has next page",
			p:        Pagination{NextPage: intPtr(2)},
			expected: true,
		},
		{
			name:     "no next page",
			p:        Pagination{NextPage: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.p.HasNextPage())
		})
	}
}

func TestPagination_HasPrevPage(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected bool
	}{
		{
			name:     "has prev page",
			p:        Pagination{PrevPage: intPtr(1)},
			expected: true,
		},
		{
			name:     "no prev page",
			p:        Pagination{PrevPage: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.p.HasPrevPage())
		})
	}
}

func TestListOptions_Values(t *testing.T) {
	tests := []struct {
		name     string
		opts     *ListOptions
		expected map[string]string
	}{
		{
			name: "with page and per_page",
			opts: &ListOptions{Page: 2, PerPage: 50},
			expected: map[string]string{
				"page":     "2",
				"per_page": "50",
			},
		},
		{
			name:     "nil options",
			opts:     nil,
			expected: map[string]string{},
		},
		{
			name: "per_page exceeds max",
			opts: &ListOptions{Page: 1, PerPage: 200},
			expected: map[string]string{
				"page":     "1",
				"per_page": "100", // capped at MaxPerPage
			},
		},
		{
			name:     "zero values",
			opts:     &ListOptions{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := tt.opts.Values()
			for k, v := range tt.expected {
				assert.Equal(t, v, values.Get(k))
			}
		})
	}
}

func TestDefaultListOptions(t *testing.T) {
	opts := DefaultListOptions()
	assert.Equal(t, DefaultPage, opts.Page)
	assert.Equal(t, DefaultPerPage, opts.PerPage)
}

func TestAddQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		opts     *ListOptions
		expected string
	}{
		{
			name:     "with params",
			path:     "/boards",
			opts:     &ListOptions{Page: 1, PerPage: 25},
			expected: "/boards?page=1&per_page=25",
		},
		{
			name:     "nil params",
			path:     "/boards",
			opts:     nil,
			expected: "/boards",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addQueryParams(tt.path, tt.opts.Values())
			assert.Equal(t, tt.expected, result)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
