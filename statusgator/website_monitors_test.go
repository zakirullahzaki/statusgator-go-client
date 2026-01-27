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

func TestWebsiteMonitorsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/website_monitors", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		var req WebsiteMonitorRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, "My Website", req.Name)
		assert.Equal(t, "https://example.com", req.URL)

		response := `{
			"success": true,
			"data": {
				"id": "wm123",
				"display_name": "My Website",
				"monitor_type": "WebsiteMonitor",
				"filtered_status": "unknown",
				"unfiltered_status": "unknown",
				"paused_at": null,
				"url": "https://example.com",
				"check_interval": 1,
				"http_method": "GET",
				"expected_status": 200,
				"timeout": 30,
				"follow_redirects": true
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &WebsiteMonitorRequest{
		Name:           "My Website",
		URL:            "https://example.com",
		CheckInterval:  1,
		ExpectedStatus: 200,
	}

	monitor, err := client.WebsiteMonitors.Create(context.Background(), "board123", req)

	require.NoError(t, err)
	assert.Equal(t, "wm123", monitor.ID)
	assert.Equal(t, "My Website", monitor.Name)
	assert.Equal(t, "https://example.com", monitor.URL)
	assert.Equal(t, 1, monitor.CheckInterval)
}

func TestWebsiteMonitorsService_Create_NilRequest(t *testing.T) {
	client, err := NewClient("test-token")
	require.NoError(t, err)

	_, err = client.WebsiteMonitors.Create(context.Background(), "board123", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")
}

func TestWebsiteMonitorsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/website_monitors/wm123", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		response := `{
			"success": true,
			"data": {
				"id": "wm123",
				"display_name": "Updated Website",
				"monitor_type": "WebsiteMonitor",
				"filtered_status": "up",
				"unfiltered_status": "up",
				"url": "https://example.com",
				"check_interval": 5
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &WebsiteMonitorRequest{
		Name:          "Updated Website",
		CheckInterval: 5,
	}

	monitor, err := client.WebsiteMonitors.Update(context.Background(), "board123", "wm123", req)

	require.NoError(t, err)
	assert.Equal(t, "Updated Website", monitor.Name)
	assert.Equal(t, 5, monitor.CheckInterval)
}

func TestWebsiteMonitorsService_Pause(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/website_monitors/wm123/pause", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.WebsiteMonitors.Pause(context.Background(), "board123", "wm123")
	assert.NoError(t, err)
}

func TestWebsiteMonitorsService_Unpause(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/website_monitors/wm123/unpause", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.WebsiteMonitors.Unpause(context.Background(), "board123", "wm123")
	assert.NoError(t, err)
}
