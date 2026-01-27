package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentsService_ListForMonitor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitors/mon123/components", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "comp123",
					"name": "API",
					"group_name": "Core Services",
					"service_id": "svc123",
					"status": "up",
					"created_at": "2024-01-15T10:30:00Z",
					"updated_at": "2024-01-15T12:00:00Z"
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

	components, pagination, err := client.Components.ListForMonitor(context.Background(), "board123", "mon123", nil)

	require.NoError(t, err)
	assert.Len(t, components, 1)
	assert.Equal(t, "comp123", components[0].ID)
	assert.Equal(t, "API", components[0].Name)
	assert.Equal(t, "Core Services", components[0].GroupName)
	assert.Equal(t, MonitorStatusUp, components[0].Status)
	assert.Equal(t, 1, pagination.TotalCount)
}

func TestComponentsService_ListByStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "affected", r.URL.Query().Get("status"))

		response := `{
			"success": true,
			"data": [
				{
					"id": "comp456",
					"name": "Database",
					"group_name": "Infrastructure",
					"service_id": "svc123",
					"status": "down"
				}
			],
			"pagination": {}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	components, err := client.Components.ListByStatus(context.Background(), "board123", "mon123", "affected")

	require.NoError(t, err)
	assert.Len(t, components, 1)
	assert.Equal(t, MonitorStatusDown, components[0].Status)
}
