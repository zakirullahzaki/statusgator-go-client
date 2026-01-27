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

func TestMonitorGroupsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitor_groups", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "grp123",
					"name": "Production",
					"position": 1,
					"collapsed": false,
					"created_at": "2024-01-15T10:30:00Z",
					"updated_at": "2024-01-15T12:00:00Z"
				}
			]
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	groups, err := client.MonitorGroups.List(context.Background(), "board123")

	require.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, "grp123", groups[0].ID)
	assert.Equal(t, "Production", groups[0].Name)
	assert.Equal(t, 1, groups[0].Position)
	assert.False(t, groups[0].Collapsed)
}

func TestMonitorGroupsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitor_groups/grp123", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": {
				"id": "grp123",
				"name": "Production",
				"position": 1,
				"collapsed": false
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	group, err := client.MonitorGroups.Get(context.Background(), "board123", "grp123")

	require.NoError(t, err)
	assert.Equal(t, "grp123", group.ID)
	assert.Equal(t, "Production", group.Name)
}

func TestMonitorGroupsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitor_groups", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		var req MonitorGroupRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, "New Group", req.Name)

		response := `{
			"success": true,
			"data": {
				"id": "grp456",
				"name": "New Group",
				"position": 2,
				"collapsed": false
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &MonitorGroupRequest{Name: "New Group", Position: 2}
	group, err := client.MonitorGroups.Create(context.Background(), "board123", req)

	require.NoError(t, err)
	assert.Equal(t, "grp456", group.ID)
	assert.Equal(t, "New Group", group.Name)
}

func TestMonitorGroupsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitor_groups/grp123", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		response := `{
			"success": true,
			"data": {
				"id": "grp123",
				"name": "Updated Group",
				"position": 1,
				"collapsed": true
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	collapsed := true
	req := &MonitorGroupRequest{Name: "Updated Group", Collapsed: &collapsed}
	group, err := client.MonitorGroups.Update(context.Background(), "board123", "grp123", req)

	require.NoError(t, err)
	assert.Equal(t, "Updated Group", group.Name)
	assert.True(t, group.Collapsed)
}

func TestMonitorGroupsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/monitor_groups/grp123", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.MonitorGroups.Delete(context.Background(), "board123", "grp123")
	assert.NoError(t, err)
}
