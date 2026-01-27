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
			"success": true,
			"data": [
				{
					"id": "user123",
					"email": "admin@example.com",
					"first_name": "Admin",
					"last_name": "User",
					"company": "Example Corp",
					"job_title": null,
					"role": "org_admin",
					"confirmed": true,
					"two_factor_enabled": true,
					"created_at": "2024-01-01T00:00:00Z",
					"last_sign_in_at": "2024-01-15T10:00:00Z"
				},
				{
					"id": "user456",
					"email": "member@example.com",
					"first_name": "Team",
					"last_name": "Member",
					"company": "Example Corp",
					"job_title": "Developer",
					"role": "member",
					"confirmed": true,
					"two_factor_enabled": false,
					"created_at": "2024-01-05T00:00:00Z",
					"last_sign_in_at": null
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
	assert.Equal(t, "Admin", users[0].FirstName)
	assert.Equal(t, "User", users[0].LastName)
	assert.Equal(t, "Admin User", users[0].FullName())
	assert.Equal(t, "Example Corp", users[0].Company)
	assert.Equal(t, "org_admin", users[0].Role)
	assert.True(t, users[0].Confirmed)
	assert.True(t, users[0].TwoFactorEnabled)
	assert.NotNil(t, users[0].LastSignInAt)

	assert.Equal(t, "Team Member", users[1].FullName())
	assert.False(t, users[1].TwoFactorEnabled)
	assert.Nil(t, users[1].LastSignInAt)
}
