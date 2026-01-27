package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// MonitorsService handles common monitor operations.
type MonitorsService service

type monitorsResponse struct {
	Success    bool       `json:"success"`
	Data       []Monitor  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// List returns all monitors for a board with pagination.
func (s *MonitorsService) List(ctx context.Context, boardID string, opts *ListOptions) ([]Monitor, *Pagination, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/boards/%s/monitors", url.PathEscape(boardID))
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp monitorsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAll returns all monitors for a board across all pages.
func (s *MonitorsService) ListAll(ctx context.Context, boardID string) ([]Monitor, error) {
	var allMonitors []Monitor
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		monitors, pagination, err := s.List(ctx, boardID, opts)
		if err != nil {
			return nil, err
		}

		allMonitors = append(allMonitors, monitors...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allMonitors, nil
}

// ListByStatus returns monitors filtered by status.
func (s *MonitorsService) ListByStatus(ctx context.Context, boardID string, status MonitorStatus) ([]Monitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("status", string(status))

	path := fmt.Sprintf("/boards/%s/monitors", url.PathEscape(boardID))
	path = addQueryParams(path, params)

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp monitorsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}

// Delete removes a monitor by ID.
func (s *MonitorsService) Delete(ctx context.Context, boardID, monitorID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return err
	}

	path := fmt.Sprintf("/boards/%s/monitors/%s", url.PathEscape(boardID), url.PathEscape(monitorID))
	return s.client.delete(ctx, path)
}
