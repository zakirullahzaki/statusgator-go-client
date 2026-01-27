package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
)

// UsersService handles user operations.
type UsersService service

type usersResponse struct {
	Success bool   `json:"success"`
	Data    []User `json:"data"`
}

// List returns all organization users.
func (s *UsersService) List(ctx context.Context) ([]User, error) {
	body, err := s.client.get(ctx, "/users")
	if err != nil {
		return nil, err
	}

	var resp usersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}
