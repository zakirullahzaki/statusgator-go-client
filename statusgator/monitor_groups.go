package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// MonitorGroupsService handles monitor group operations.
type MonitorGroupsService service

type monitorGroupsResponse struct {
	Success bool           `json:"success"`
	Data    []MonitorGroup `json:"data"`
}

type monitorGroupResponse struct {
	Success bool         `json:"success"`
	Data    MonitorGroup `json:"data"`
}

// List returns all monitor groups for a board.
func (s *MonitorGroupsService) List(ctx context.Context, boardID string) ([]MonitorGroup, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/boards/%s/monitor_groups", url.PathEscape(boardID))

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp monitorGroupsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}

// Get returns a specific monitor group.
func (s *MonitorGroupsService) Get(ctx context.Context, boardID, groupID string) (*MonitorGroup, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(groupID, "groupID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/boards/%s/monitor_groups/%s", url.PathEscape(boardID), url.PathEscape(groupID))

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp monitorGroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Create creates a new monitor group.
func (s *MonitorGroupsService) Create(ctx context.Context, boardID string, req *MonitorGroupRequest) (*MonitorGroup, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/monitor_groups", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp monitorGroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Update updates a monitor group.
func (s *MonitorGroupsService) Update(ctx context.Context, boardID, groupID string, req *MonitorGroupRequest) (*MonitorGroup, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(groupID, "groupID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/monitor_groups/%s", url.PathEscape(boardID), url.PathEscape(groupID))

	body, err := s.client.patch(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp monitorGroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Delete removes a monitor group.
func (s *MonitorGroupsService) Delete(ctx context.Context, boardID, groupID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(groupID, "groupID"); err != nil {
		return err
	}

	path := fmt.Sprintf("/boards/%s/monitor_groups/%s", url.PathEscape(boardID), url.PathEscape(groupID))
	return s.client.delete(ctx, path)
}
