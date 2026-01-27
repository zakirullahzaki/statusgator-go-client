package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
)

// UsersService handles user operations.
type UsersService service

type usersResponse struct {
	Users []User `json:"users"`
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

	return resp.Users, nil
}
