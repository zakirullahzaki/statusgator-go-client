package statusgator

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 25
	MaxPerPage     = 100
)

// Pagination contains pagination metadata from API responses.
type Pagination struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
}

// HasNextPage returns true if there are more pages available.
func (p *Pagination) HasNextPage() bool {
	return p.NextPage != nil
}

// HasPrevPage returns true if there are previous pages.
func (p *Pagination) HasPrevPage() bool {
	return p.PrevPage != nil
}

// ListOptions specifies pagination parameters for list operations.
type ListOptions struct {
	Page    int
	PerPage int
}

// DefaultListOptions returns default pagination options.
func DefaultListOptions() *ListOptions {
	return &ListOptions{
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
	}
}

// Values returns URL query parameters for the list options.
func (o *ListOptions) Values() url.Values {
	v := url.Values{}
	if o == nil {
		return v
	}
	if o.Page > 0 {
		v.Set("page", strconv.Itoa(o.Page))
	}
	if o.PerPage > 0 {
		perPage := o.PerPage
		if perPage > MaxPerPage {
			perPage = MaxPerPage
		}
		v.Set("per_page", strconv.Itoa(perPage))
	}
	return v
}

// addToPath appends query parameters to a path.
func addQueryParams(path string, params url.Values) string {
	if len(params) == 0 {
		return path
	}
	return fmt.Sprintf("%s?%s", path, params.Encode())
}
