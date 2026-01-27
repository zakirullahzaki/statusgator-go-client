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

func TestSubscribersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/status_page_subscribers", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "sub123",
					"email": "user@example.com",
					"confirmed": true,
					"created_at": "2024-01-15T10:30:00Z"
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

	subscribers, pagination, err := client.Subscribers.List(context.Background(), "board123", nil)

	require.NoError(t, err)
	assert.Len(t, subscribers, 1)
	assert.Equal(t, "sub123", subscribers[0].ID)
	assert.Equal(t, "user@example.com", subscribers[0].Email)
	assert.True(t, subscribers[0].Confirmed)
	assert.Equal(t, 1, pagination.TotalCount)
}

func TestSubscribersService_Add(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/status_page_subscribers", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		var req SubscriberRequest
		json.Unmarshal(body, &req)
		assert.Equal(t, "new@example.com", req.Email)
		assert.True(t, req.SkipConfirmation)

		response := `{
			"success": true,
			"data": {
				"id": "sub456",
				"email": "new@example.com",
				"confirmed": true,
				"created_at": "2024-01-15T15:00:00Z"
			}
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	req := &SubscriberRequest{
		Email:            "new@example.com",
		SkipConfirmation: true,
	}

	subscriber, err := client.Subscribers.Add(context.Background(), "board123", req)

	require.NoError(t, err)
	assert.Equal(t, "sub456", subscriber.ID)
	assert.Equal(t, "new@example.com", subscriber.Email)
	assert.True(t, subscriber.Confirmed)
}

func TestSubscribersService_DeleteByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/status_page_subscribers", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "sub123", r.URL.Query().Get("id"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.Subscribers.DeleteByID(context.Background(), "board123", "sub123")
	assert.NoError(t, err)
}

func TestSubscribersService_DeleteByEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/status_page_subscribers", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "user@example.com", r.URL.Query().Get("email"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.Subscribers.DeleteByEmail(context.Background(), "board123", "user@example.com")
	assert.NoError(t, err)
}
