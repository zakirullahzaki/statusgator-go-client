package statusgator

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/incidents", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "inc123",
					"name": "API Degradation",
					"details": "We are investigating reports of API issues",
					"severity": "minor",
					"phase": "investigating",
					"board_id": "board123",
					"scheduled_maintenance": false,
					"resolved_or_completed": false,
					"auto_complete_maintenance": false,
					"created_at": "2024-01-15T14:30:00Z",
					"updated_at": "2024-01-15T14:35:00Z"
				}
			],
			"pagination": {
				"current_page": 1,
				"per_page": 25,
				"total_pages": 1,
				"total_count": 1
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	incidents, pagination, err := client.Incidents.List(context.Background(), "board123", nil)

	require.NoError(t, err)
	assert.Len(t, incidents, 1)
	assert.Equal(t, "inc123", incidents[0].ID)
	assert.Equal(t, "API Degradation", incidents[0].Name)
	assert.Equal(t, "We are investigating reports of API issues", incidents[0].Details)
	assert.Equal(t, IncidentSeverityMinor, incidents[0].Severity)
	assert.Equal(t, IncidentPhaseInvestigating, incidents[0].Phase)
	assert.Equal(t, "board123", incidents[0].BoardID)
	assert.Equal(t, 1, pagination.TotalCount)
}

func TestIncidentsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/incidents", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		var req IncidentRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, "Service Outage", req.Name)
		assert.Equal(t, IncidentSeverityMajor, req.Severity)

		response := `{
			"success": true,
			"data": {
				"id": "inc456",
				"name": "Service Outage",
				"details": "We are investigating reports of service issues",
				"severity": "major",
				"phase": "investigating",
				"board_id": "board123",
				"scheduled_maintenance": false,
				"resolved_or_completed": false,
				"auto_complete_maintenance": false,
				"created_at": "2024-01-15T15:00:00Z",
				"updated_at": "2024-01-15T15:00:00Z"
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &IncidentRequest{
		Name:     "Service Outage",
		Details:  "We are investigating reports of service issues",
		Severity: IncidentSeverityMajor,
		Phase:    IncidentPhaseInvestigating,
	}

	incident, err := client.Incidents.Create(context.Background(), "board123", req)

	require.NoError(t, err)
	assert.Equal(t, "inc456", incident.ID)
	assert.Equal(t, "Service Outage", incident.Name)
	assert.Equal(t, IncidentSeverityMajor, incident.Severity)
}

func TestIncidentsService_AddUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/incidents/inc123/incident_updates", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		var req IncidentUpdateRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, "Issue has been identified", req.Details)
		assert.Equal(t, IncidentPhaseIdentified, req.Phase)

		response := `{
			"success": true,
			"data": {
				"id": "upd123",
				"incident_id": "inc123",
				"details": "Issue has been identified",
				"phase": "identified",
				"severity": "major",
				"notify_subscribers": true,
				"created_at": "2024-01-15T15:30:00Z",
				"updated_at": "2024-01-15T15:30:00Z"
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &IncidentUpdateRequest{
		Details: "Issue has been identified",
		Phase:   IncidentPhaseIdentified,
	}

	update, err := client.Incidents.AddUpdate(context.Background(), "board123", "inc123", req)

	require.NoError(t, err)
	assert.Equal(t, "upd123", update.ID)
	assert.Equal(t, "inc123", update.IncidentID)
	assert.Equal(t, "Issue has been identified", update.Details)
	assert.Equal(t, IncidentPhaseIdentified, update.Phase)
}
