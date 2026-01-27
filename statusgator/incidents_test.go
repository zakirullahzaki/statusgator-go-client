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
					"title": "API Degradation",
					"severity": "minor",
					"phase": "investigating",
					"monitor_ids": ["mon123"],
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
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	incidents, pagination, err := client.Incidents.List(context.Background(), "board123", nil)

	require.NoError(t, err)
	assert.Len(t, incidents, 1)
	assert.Equal(t, "inc123", incidents[0].ID)
	assert.Equal(t, "API Degradation", incidents[0].Title)
	assert.Equal(t, IncidentSeverityMinor, incidents[0].Severity)
	assert.Equal(t, IncidentPhaseInvestigating, incidents[0].Phase)
	assert.Equal(t, []string{"mon123"}, incidents[0].MonitorIDs)
	assert.Equal(t, 1, pagination.TotalCount)
}

func TestIncidentsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/incidents", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		var req IncidentRequest
		json.Unmarshal(body, &req)
		assert.Equal(t, "Service Outage", req.Title)
		assert.Equal(t, IncidentSeverityMajor, req.Severity)

		response := `{
			"success": true,
			"data": {
				"id": "inc456",
				"title": "Service Outage",
				"severity": "major",
				"phase": "investigating",
				"monitor_ids": ["mon123", "mon456"],
				"created_at": "2024-01-15T15:00:00Z",
				"updated_at": "2024-01-15T15:00:00Z"
			}
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &IncidentRequest{
		Title:      "Service Outage",
		Message:    "We are investigating reports of service issues",
		Severity:   IncidentSeverityMajor,
		Phase:      IncidentPhaseInvestigating,
		MonitorIDs: []string{"mon123", "mon456"},
	}

	incident, err := client.Incidents.Create(context.Background(), "board123", req)

	require.NoError(t, err)
	assert.Equal(t, "inc456", incident.ID)
	assert.Equal(t, "Service Outage", incident.Title)
	assert.Equal(t, IncidentSeverityMajor, incident.Severity)
}

func TestIncidentsService_AddUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/incidents/inc123/incident_updates", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		var req IncidentUpdateRequest
		json.Unmarshal(body, &req)
		assert.Equal(t, "Issue has been identified", req.Message)
		assert.Equal(t, IncidentPhaseIdentified, req.Phase)

		response := `{
			"success": true,
			"data": {
				"id": "upd123",
				"message": "Issue has been identified",
				"phase": "identified",
				"created_at": "2024-01-15T15:30:00Z"
			}
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &IncidentUpdateRequest{
		Message: "Issue has been identified",
		Phase:   IncidentPhaseIdentified,
	}

	update, err := client.Incidents.AddUpdate(context.Background(), "board123", "inc123", req)

	require.NoError(t, err)
	assert.Equal(t, "upd123", update.ID)
	assert.Equal(t, "Issue has been identified", update.Message)
	assert.Equal(t, IncidentPhaseIdentified, update.Phase)
}
