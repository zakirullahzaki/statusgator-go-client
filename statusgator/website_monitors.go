package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// WebsiteMonitorsService handles website monitor operations.
type WebsiteMonitorsService service

type websiteMonitorResponse struct {
	Success bool           `json:"success"`
	Data    WebsiteMonitor `json:"data"`
}

// Create creates a new website monitor.
func (s *WebsiteMonitorsService) Create(ctx context.Context, boardID string, req *WebsiteMonitorRequest) (*WebsiteMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/website_monitors", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp websiteMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Update updates an existing website monitor.
func (s *WebsiteMonitorsService) Update(ctx context.Context, boardID, monitorID string, req *WebsiteMonitorRequest) (*WebsiteMonitor, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/website_monitors/%s", url.PathEscape(boardID), url.PathEscape(monitorID))

	body, err := s.client.patch(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp websiteMonitorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// Pause pauses the website monitor.
func (s *WebsiteMonitorsService) Pause(ctx context.Context, boardID, monitorID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return err
	}

	path := fmt.Sprintf("/boards/%s/website_monitors/%s/pause", url.PathEscape(boardID), url.PathEscape(monitorID))
	_, err := s.client.post(ctx, path, nil)
	return err
}

// Unpause resumes the website monitor.
func (s *WebsiteMonitorsService) Unpause(ctx context.Context, boardID, monitorID string) error {
	if err := validateID(boardID, "boardID"); err != nil {
		return err
	}
	if err := validateID(monitorID, "monitorID"); err != nil {
		return err
	}

	path := fmt.Sprintf("/boards/%s/website_monitors/%s/unpause", url.PathEscape(boardID), url.PathEscape(monitorID))
	_, err := s.client.post(ctx, path, nil)
	return err
}
