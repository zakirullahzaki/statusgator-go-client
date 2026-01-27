package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// CustomMonitorsService handles custom monitor operations.
type CustomMonitorsService service

type customMonitorResponse struct {
	Success bool          `json:"success"`
	Data    CustomMonitor `json:"data"`
}

// Create creates a new custom monitor.
func (s *CustomMonitorsService) Create(ctx context.Context, boardID string, req *CustomMonitorRequest) (*CustomMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/custom_monitors", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp customMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Update updates a custom monitor.
func (s *CustomMonitorsService) Update(ctx context.Context, boardID, monitorID string, req *CustomMonitorRequest) (*CustomMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/custom_monitors/%s", url.PathEscape(boardID), url.PathEscape(monitorID))

	body, err := s.client.patch(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp customMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// SetStatus updates the status of a custom monitor.
func (s *CustomMonitorsService) SetStatus(ctx context.Context, boardID, monitorID string, status MonitorStatus) error {
	req := &CustomMonitorRequest{Status: status}
	_, err := s.Update(ctx, boardID, monitorID, req)
	return err
}
