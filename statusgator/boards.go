package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// BoardsService handles board-related API operations.
type BoardsService service

type boardsResponse struct {
	Success    bool       `json:"success"`
	Data       []Board    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type boardResponse struct {
	Success bool  `json:"success"`
	Data    Board `json:"data"`
}

type historyResponse struct {
	Success bool           `json:"success"`
	Data    []HistoryEvent `json:"data"`
}

// List returns all boards with pagination.
func (s *BoardsService) List(ctx context.Context, opts *ListOptions) ([]Board, *Pagination, error) {
	path := "/boards"
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp boardsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAll returns all boards across all pages.
func (s *BoardsService) ListAll(ctx context.Context) ([]Board, error) {
	var allBoards []Board
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		boards, pagination, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allBoards = append(allBoards, boards...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allBoards, nil
}

// Get returns a specific board by ID.
func (s *BoardsService) Get(ctx context.Context, boardID string) (*Board, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/boards/%s", url.PathEscape(boardID))

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp boardResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// GetHistory returns historical events for a board.
func (s *BoardsService) GetHistory(ctx context.Context, boardID string, opts *HistoryOptions) ([]HistoryEvent, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}

	params := url.Values{}
	if opts != nil {
		if opts.StartDate != "" {
			params.Set("start_date", opts.StartDate)
		}
		if opts.EndDate != "" {
			params.Set("end_date", opts.EndDate)
		}
		if opts.MonitorID != "" {
			params.Set("monitor_id", opts.MonitorID)
		}
	}

	path := fmt.Sprintf("/boards/%s/history", url.PathEscape(boardID))
	path = addQueryParams(path, params)

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp historyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}
