package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServicesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/services", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "svc123",
					"name": "GitHub",
					"slug": "github",
					"status_url": "https://www.githubstatus.com",
					"created_at": "2024-01-01T00:00:00Z",
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
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	services, pagination, err := client.Services.List(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Equal(t, "svc123", services[0].ID)
	assert.Equal(t, "GitHub", services[0].Name)
	assert.Equal(t, "github", services[0].Slug)
	assert.Equal(t, "https://www.githubstatus.com", services[0].StatusURL)
	assert.Equal(t, 1, pagination.TotalCount)
}

func TestServicesService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/services/search", r.URL.Path)
		assert.Equal(t, "aws", r.URL.Query().Get("query"))

		response := `{
			"success": true,
			"data": [
				{
					"id": "svc456",
					"name": "AWS",
					"slug": "aws",
					"status_url": "https://health.aws.amazon.com"
				}
			]
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	services, err := client.Services.Search(context.Background(), "aws")

	require.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Equal(t, "AWS", services[0].Name)
}

func TestServicesService_ListComponents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/services/svc123/components", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "comp123",
					"name": "API Requests",
					"group_name": "GitHub Actions",
					"service_id": "svc123",
					"status": "up"
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

	components, pagination, err := client.Services.ListComponents(context.Background(), "svc123", nil)

	require.NoError(t, err)
	assert.Len(t, components, 1)
	assert.Equal(t, "comp123", components[0].ID)
	assert.Equal(t, 1, pagination.TotalCount)
}
