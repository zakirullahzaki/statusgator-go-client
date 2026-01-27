package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := `{
			"success": true,
			"data": [
				{
					"id": "board123",
					"name": "Production Services",
					"public_token": "pub_token_123",
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
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	boards, pagination, err := client.Boards.List(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, boards, 1)
	assert.Equal(t, "board123", boards[0].ID)
	assert.Equal(t, "Production Services", boards[0].Name)
	assert.Equal(t, "pub_token_123", boards[0].PublicToken)
	assert.Equal(t, 1, pagination.CurrentPage)
	assert.Equal(t, 1, pagination.TotalCount)
	assert.False(t, pagination.HasNextPage())
}

func TestBoardsService_ListWithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("per_page"))

		response := `{
			"success": true,
			"data": [],
			"pagination": {
				"current_page": 2,
				"per_page": 10,
				"total_pages": 2,
				"total_count": 15,
				"next_page": null,
				"prev_page": 1
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	opts := &ListOptions{Page: 2, PerPage: 10}
	_, pagination, err := client.Boards.List(context.Background(), opts)

	require.NoError(t, err)
	assert.Equal(t, 2, pagination.CurrentPage)
	assert.True(t, pagination.HasPrevPage())
}

func TestBoardsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": {
				"id": "board123",
				"name": "Production Services",
				"public_token": "pub_token_123",
				"created_at": "2024-01-15T10:30:00Z",
				"updated_at": "2024-01-15T12:00:00Z"
			}
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	board, err := client.Boards.Get(context.Background(), "board123")

	require.NoError(t, err)
	assert.Equal(t, "board123", board.ID)
	assert.Equal(t, "Production Services", board.Name)
}

func TestBoardsService_Get_EmptyID(t *testing.T) {
	client, err := NewClient("test-token")
	require.NoError(t, err)

	_, err = client.Boards.Get(context.Background(), "")
	assert.ErrorIs(t, err, ErrEmptyID)
}

func TestBoardsService_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "Board not found"}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.Boards.Get(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.True(t, IsNotFound(err))
}

func TestBoardsService_GetHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/boards/board123/history", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "2024-01-01", r.URL.Query().Get("start_date"))
		assert.Equal(t, "2024-01-31", r.URL.Query().Get("end_date"))

		response := `{
			"success": true,
			"data": [
				{
					"monitor_id": "mon123",
					"name": "API Service",
					"icon_url": "https://example.com/icon.png",
					"status": "down",
					"started_at": "2024-01-15T14:30:00Z",
					"ended_at": "2024-01-15T15:00:00Z",
					"duration": "30.0",
					"message": "Service unavailable",
					"details": "Connection timeout",
					"early_warning_signal": false
				}
			]
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	opts := &HistoryOptions{
		StartDate: "2024-01-01",
		EndDate:   "2024-01-31",
	}
	history, err := client.Boards.GetHistory(context.Background(), "board123", opts)

	require.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, "mon123", history[0].MonitorID)
	assert.Equal(t, "API Service", history[0].Name)
	assert.Equal(t, MonitorStatusDown, history[0].Status)
	assert.Equal(t, "Connection timeout", history[0].Details)
}

func TestBoardsService_ListAll(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		var response string
		if callCount == 1 {
			response = `{
				"success": true,
				"data": [{"id": "board1", "name": "Board 1"}],
				"pagination": {
					"current_page": 1,
					"per_page": 100,
					"total_pages": 2,
					"total_count": 2,
					"next_page": 2,
					"prev_page": null
				}
			}`
		} else {
			response = `{
				"success": true,
				"data": [{"id": "board2", "name": "Board 2"}],
				"pagination": {
					"current_page": 2,
					"per_page": 100,
					"total_pages": 2,
					"total_count": 2,
					"next_page": null,
					"prev_page": 1
				}
			}`
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	boards, err := client.Boards.ListAll(context.Background())

	require.NoError(t, err)
	assert.Len(t, boards, 2)
	assert.Equal(t, "board1", boards[0].ID)
	assert.Equal(t, "board2", boards[1].ID)
	assert.Equal(t, 2, callCount)
}
