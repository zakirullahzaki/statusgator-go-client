package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestMonitorsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitors", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "mon123",
					"display_name": "API Health",
					"monitor_type": "WebsiteMonitor",
					"filtered_status": "up",
					"unfiltered_status": "up",
					"paused_at": null,
					"checked_at": "2024-01-15T12:00:00Z",
					"filter_count": 0,
					"icon_url": "https://example.com/icon.png",
					"position": 1,
					"early_warning_signal": false,
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
	assert.False(t, monitors[0].IsPaused())
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
					"display_name": "Database",
					"monitor_type": "PingMonitor",
					"filtered_status": "down",
					"unfiltered_status": "down",
					"paused_at": null
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

func TestMonitor_IsPaused(t *testing.T) {
	t.Run("paused monitor", func(t *testing.T) {
		pausedTime := parseTime("2024-01-15T10:00:00Z")
		m := Monitor{PausedAt: &pausedTime}
		assert.True(t, m.IsPaused())
	})

	t.Run("not paused monitor", func(t *testing.T) {
		m := Monitor{PausedAt: nil}
		assert.False(t, m.IsPaused())
	})
}

func TestMonitor_WithServiceAndGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"success": true,
			"data": [
				{
					"id": "mon789",
					"display_name": "Cloudflare",
					"monitor_type": "ServiceMonitor",
					"filtered_status": "warn",
					"unfiltered_status": "warn",
					"last_message": "CGK (Jakarta) on 2025-09-05",
					"last_details": "In progress - Scheduled maintenance",
					"paused_at": null,
					"checked_at": "2025-09-05T23:05:23Z",
					"filter_count": 0,
					"icon_url": "https://example.com/cloudflare.png",
					"position": 3,
					"early_warning_signal": false,
					"service": {
						"id": "svc123",
						"name": "Cloudflare",
						"slug": "cloudflare",
						"home_page_url": "http://www.cloudflare.com",
						"status_page_url": "https://www.cloudflarestatus.com",
						"icon_url": "https://example.com/cloudflare.png",
						"landing_page_url": "https://statusgator.com/services/cloudflare",
						"official": true
					},
					"group": {
						"id": "grp456",
						"name": "Infrastructure",
						"position": 1
					},
					"created_at": "2023-07-06T14:18:35Z",
					"updated_at": "2023-07-06T14:18:35Z"
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

	monitors, _, err := client.Monitors.List(context.Background(), "board123", nil)

	require.NoError(t, err)
	require.Len(t, monitors, 1)

	m := monitors[0]
	assert.Equal(t, "mon789", m.ID)
	assert.Equal(t, "Cloudflare", m.Name)
	assert.Equal(t, MonitorTypeService, m.Type)
	assert.Equal(t, MonitorStatusWarn, m.Status)
	assert.Equal(t, MonitorStatusWarn, m.UnfilteredStatus)

	require.NotNil(t, m.LastMessage)
	assert.Equal(t, "CGK (Jakarta) on 2025-09-05", *m.LastMessage)

	require.NotNil(t, m.LastDetails)
	assert.Equal(t, "In progress - Scheduled maintenance", *m.LastDetails)

	require.NotNil(t, m.Service)
	assert.Equal(t, "svc123", m.Service.ID)
	assert.Equal(t, "Cloudflare", m.Service.Name)
	assert.Equal(t, "cloudflare", m.Service.Slug)
	assert.Equal(t, "http://www.cloudflare.com", m.Service.HomePageURL)
	assert.Equal(t, "https://www.cloudflarestatus.com", m.Service.StatusPageURL)
	assert.Equal(t, "https://statusgator.com/services/cloudflare", m.Service.LandingPageURL)
	assert.True(t, m.Service.Official)

	require.NotNil(t, m.Group)
	assert.Equal(t, "grp456", m.Group.ID)
	assert.Equal(t, "Infrastructure", m.Group.Name)
	assert.Equal(t, 1, m.Group.Position)
}
