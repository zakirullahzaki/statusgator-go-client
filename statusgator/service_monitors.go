package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ServiceMonitorsService handles service monitor operations.
type ServiceMonitorsService service

type serviceMonitorResponse struct {
	Success bool           `json:"success"`
	Data    ServiceMonitor `json:"data"`
}

// Create subscribes to an external status page.
func (s *ServiceMonitorsService) Create(ctx context.Context, boardID string, req *ServiceMonitorRequest) (*ServiceMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/service_monitors", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp serviceMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Update updates a service monitor.
func (s *ServiceMonitorsService) Update(ctx context.Context, boardID, monitorID string, req *ServiceMonitorRequest) (*ServiceMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/service_monitors/%s", url.PathEscape(boardID), url.PathEscape(monitorID))

	body, err := s.client.patch(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp serviceMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}
