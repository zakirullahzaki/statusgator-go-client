package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// StatusPageSubscribersService handles subscriber operations.
type StatusPageSubscribersService service

type subscribersResponse struct {
	Success    bool         `json:"success"`
	Data       []Subscriber `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

type subscriberResponse struct {
	Success bool       `json:"success"`
	Data    Subscriber `json:"data"`
}

// List returns all subscribers for a board with pagination.
func (s *StatusPageSubscribersService) List(ctx context.Context, boardID string, opts *ListOptions) ([]Subscriber, *Pagination, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/boards/%s/status_page_subscribers", url.PathEscape(boardID))
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp subscribersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAll returns all subscribers for a board across all pages.
func (s *StatusPageSubscribersService) ListAll(ctx context.Context, boardID string) ([]Subscriber, error) {
	var allSubscribers []Subscriber
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		subscribers, pagination, err := s.List(ctx, boardID, opts)
		if err != nil {
			return nil, err
		}

		allSubscribers = append(allSubscribers, subscribers...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allSubscribers, nil
}

// Add adds a new subscriber to a board.
func (s *StatusPageSubscribersService) Add(ctx context.Context, boardID string, req *SubscriberRequest) (*Subscriber, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/status_page_subscribers", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp subscriberResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// DeleteByID removes a subscriber by ID.
func (s *StatusPageSubscribersService) DeleteByID(ctx context.Context, boardID, subscriberID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(subscriberID, "subscriberID"); err != nil {
		return err
	}

	params := url.Values{}
	params.Set("id", subscriberID)

	path := fmt.Sprintf("/boards/%s/status_page_subscribers", url.PathEscape(boardID))
	return s.client.deleteWithQuery(ctx, path, params)
}

// DeleteByEmail removes a subscriber by email.
func (s *StatusPageSubscribersService) DeleteByEmail(ctx context.Context, boardID, email string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	params := url.Values{}
	params.Set("email", email)

	path := fmt.Sprintf("/boards/%s/status_page_subscribers", url.PathEscape(boardID))
	return s.client.deleteWithQuery(ctx, path, params)
}
