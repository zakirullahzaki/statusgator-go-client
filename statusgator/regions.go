package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
)

// RegionsService handles monitoring region operations.
type RegionsService service

type regionsResponse struct {
	Success bool     `json:"success"`
	Data    []Region `json:"data"`
}

// List returns all available monitoring regions.
func (s *RegionsService) List(ctx context.Context) ([]Region, error) {
	body, err := s.client.get(ctx, "/monitoring_regions")
	if err != nil {
		return nil, err
	}

	var resp regionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}
