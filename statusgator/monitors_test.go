package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitors", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "mon123",
					"name": "API Health",
					"type": "website",
					"status": "up",
					"paused": false,
					"created_at": "2024-01-15T10:30:00Z",
					"updated_at": "2024-01-15T12:00:00Z"
				}
			],
			"pagination": {
				"current_page": 1,
				"per_page": 25,
				"total_pages": 1,
				"total_count": 1,
				"next_page": null,
				"prev_page": null
			}
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	monitors, pagination, err := client.Monitors.List(context.Background(), "board123", nil)

	require.NoError(t, err)
	assert.Len(t, monitors, 1)
	assert.Equal(t, "mon123", monitors[0].ID)
	assert.Equal(t, "API Health", monitors[0].Name)
	assert.Equal(t, MonitorTypeWebsite, monitors[0].Type)
	assert.Equal(t, MonitorStatusUp, monitors[0].Status)
	assert.False(t, monitors[0].Paused)
	assert.Equal(t, 1, pagination.TotalCount)
}

func TestMonitorsService_List_EmptyBoardID(t *testing.T) {
	client, err := NewClient("test-token")
	require.NoError(t, err)

	_, _, err = client.Monitors.List(context.Background(), "", nil)
	assert.ErrorIs(t, err, ErrEmptyID)
}

func TestMonitorsService_ListByStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "down", r.URL.Query().Get("status"))

		response := `{
			"success": true,
			"data": [
				{
					"id": "mon456",
					"name": "Database",
					"type": "ping",
					"status": "down",
					"paused": false
				}
			],
			"pagination": {}
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	monitors, err := client.Monitors.ListByStatus(context.Background(), "board123", MonitorStatusDown)

	require.NoError(t, err)
	assert.Len(t, monitors, 1)
	assert.Equal(t, MonitorStatusDown, monitors[0].Status)
}

func TestMonitorsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitors/mon123", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.Monitors.Delete(context.Background(), "board123", "mon123")
	assert.NoError(t, err)
}

func TestMonitorsService_Delete_EmptyIDs(t *testing.T) {
	client, err := NewClient("test-token")
	require.NoError(t, err)

	err = client.Monitors.Delete(context.Background(), "", "mon123")
	assert.ErrorIs(t, err, ErrEmptyID)

	err = client.Monitors.Delete(context.Background(), "board123", "")
	assert.ErrorIs(t, err, ErrEmptyID)
}
