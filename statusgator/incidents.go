package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// IncidentsService handles incident operations.
type IncidentsService service

type incidentsResponse struct {
	Success    bool       `json:"success"`
	Data       []Incident `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type incidentResponse struct {
	Success bool     `json:"success"`
	Data    Incident `json:"data"`
}

type incidentUpdateResponse struct {
	Success bool           `json:"success"`
	Data    IncidentUpdate `json:"data"`
}

// List returns incidents for a board with pagination.
func (s *IncidentsService) List(ctx context.Context, boardID string, opts *ListOptions) ([]Incident, *Pagination, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/boards/%s/incidents", url.PathEscape(boardID))
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp incidentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAll returns all incidents for a board across all pages.
func (s *IncidentsService) ListAll(ctx context.Context, boardID string) ([]Incident, error) {
	var allIncidents []Incident
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		incidents, pagination, err := s.List(ctx, boardID, opts)
		if err != nil {
			return nil, err
		}

		allIncidents = append(allIncidents, incidents...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allIncidents, nil
}

// Create creates a new incident or maintenance window.
func (s *IncidentsService) Create(ctx context.Context, boardID string, req *IncidentRequest) (*Incident, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/incidents", url.PathEscape(boardID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp incidentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}

// AddUpdate adds a status update to an incident.
func (s *IncidentsService) AddUpdate(ctx context.Context, boardID, incidentID string, req *IncidentUpdateRequest) (*IncidentUpdate, error) {
	if err := validateID(boardID, "boardID"); err != nil {
		return nil, err
	}
	if err := validateID(incidentID, "incidentID"); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := fmt.Sprintf("/boards/%s/incidents/%s/incident_updates", url.PathEscape(boardID), url.PathEscape(incidentID))

	body, err := s.client.post(ctx, path, req)
	if err != nil {
		return nil, err
	}

	var resp incidentUpdateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &resp.Data, nil
}
