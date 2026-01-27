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

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		opts    []ClientOption
		wantErr error
	}{
		{
			name:    "valid token",
			token:   "test-token",
			wantErr: nil,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: ErrTokenRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token, tt.opts...)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.token, client.Token)
				assert.Equal(t, DefaultBaseURL, client.BaseURL)
				assert.Equal(t, DefaultUserAgent, client.UserAgent)
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	customHTTPClient := &http.Client{Timeout: 60 * time.Second}
	customBaseURL := "https://custom.api.com/v3"
	customUserAgent := "custom-agent/1.0"

	client, err := NewClient("test-token",
		WithHTTPClient(customHTTPClient),
		WithBaseURL(customBaseURL+"/"),
		WithUserAgent(customUserAgent),
	)

	require.NoError(t, err)
	assert.Equal(t, customHTTPClient, client.client)
	assert.Equal(t, customBaseURL, client.BaseURL) // trailing slash removed
	assert.Equal(t, customUserAgent, client.UserAgent)
}

func TestClient_Ping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ping", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.Ping(context.Background())
	assert.NoError(t, err)
}

func TestClient_PingUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message": "Invalid token"}`))
	}))
	defer server.Close()

	client, err := NewClient("invalid-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.Ping(context.Background())
	assert.Error(t, err)
	assert.True(t, IsUnauthorized(err))
}

func TestClient_ResponseTooLarge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write more than DefaultMaxResponseSize
		largeResponse := make([]byte, DefaultMaxResponseSize+100)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(largeResponse)
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	err = client.Ping(context.Background())
	assert.ErrorIs(t, err, ErrResponseTooLarge)
}

func TestClient_ServicesInitialized(t *testing.T) {
	client, err := NewClient("test-token")
	require.NoError(t, err)

	assert.NotNil(t, client.Boards)
	assert.NotNil(t, client.Monitors)
	assert.NotNil(t, client.WebsiteMonitors)
	assert.NotNil(t, client.PingMonitors)
	assert.NotNil(t, client.ServiceMonitors)
	assert.NotNil(t, client.CustomMonitors)
	assert.NotNil(t, client.MonitorGroups)
	assert.NotNil(t, client.Components)
	assert.NotNil(t, client.Incidents)
	assert.NotNil(t, client.Services)
	assert.NotNil(t, client.Subscribers)
	assert.NotNil(t, client.Users)
	assert.NotNil(t, client.Regions)
}

