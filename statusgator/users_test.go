package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"users": [
				{
					"id": "user123",
					"email": "admin@example.com",
					"name": "Admin User",
					"role": "admin"
				},
				{
					"id": "user456",
					"email": "member@example.com",
					"name": "Team Member",
					"role": "member"
				}
			]
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	users, err := client.Users.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user123", users[0].ID)
	assert.Equal(t, "admin@example.com", users[0].Email)
	assert.Equal(t, "Admin User", users[0].Name)
	assert.Equal(t, "admin", users[0].Role)
}
