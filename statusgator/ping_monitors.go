package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// PingMonitorsService handles ping monitor operations.
type PingMonitorsService service

type pingMonitorResponse struct {
	Success bool        `json:"success"`
	Data    PingMonitor `json:"data"`
}

// Create creates a new ping monitor.
func (s *PingMonitorsService) Create(ctx context.Context, boardID string, req *PingMonitorRequest) (*PingMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/ping_monitors", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp pingMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Update updates an existing ping monitor.
func (s *PingMonitorsService) Update(ctx context.Context, boardID, monitorID string, req *PingMonitorRequest) (*PingMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/ping_monitors/%s", url.PathEscape(boardID), url.PathEscape(monitorID))

	body, err := s.client.patch(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp pingMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Pause pauses the ping monitor.
func (s *PingMonitorsService) Pause(ctx context.Context, boardID, monitorID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return err
	}

	path := fmt.Sprintf("/boards/%s/ping_monitors/%s/pause", url.PathEscape(boardID), url.PathEscape(monitorID))
	_, err := s.client.post(ctx, path, nil)
	return err
}

// Unpause resumes the ping monitor.
func (s *PingMonitorsService) Unpause(ctx context.Context, boardID, monitorID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return err
	}

	path := fmt.Sprintf("/boards/%s/ping_monitors/%s/unpause", url.PathEscape(boardID), url.PathEscape(monitorID))
	_, err := s.client.post(ctx, path, nil)
	return err
}
