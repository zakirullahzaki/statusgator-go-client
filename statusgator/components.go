package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ComponentsService handles component operations.
type ComponentsService service

type componentsResponse struct {
	Success    bool        `json:"success"`
	Data       []Component `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// ListForMonitor returns components for a monitor with pagination.
func (s *ComponentsService) ListForMonitor(ctx context.Context, boardID, monitorID string, opts *ListOptions) ([]Component, *Pagination, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, nil, err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/boards/%s/monitors/%s/components", url.PathEscape(boardID), url.PathEscape(monitorID))
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp componentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAllForMonitor returns all components for a monitor across all pages.
func (s *ComponentsService) ListAllForMonitor(ctx context.Context, boardID, monitorID string) ([]Component, error) {
	var allComponents []Component
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		components, pagination, err := s.ListForMonitor(ctx, boardID, monitorID, opts)
		if err != nil {
			return nil, err
		}

		allComponents = append(allComponents, components...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allComponents, nil
}

// ListByStatus returns components filtered by status.
// Use "affected" to get all non-up components, or specify status values
// (comma-separated for multiple: "down,warn").
func (s *ComponentsService) ListByStatus(ctx context.Context, boardID, monitorID, status string) ([]Component, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("status", status)

	path := fmt.Sprintf("/boards/%s/monitors/%s/components", url.PathEscape(boardID), url.PathEscape(monitorID))
	path = addQueryParams(path, params)

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp componentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}
