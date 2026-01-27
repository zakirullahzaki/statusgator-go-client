package statusgator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/monitoring_regions", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{
			"success": true,
			"data": [
				{
					"id": "region123",
					"name": "US East",
					"code": "us-east-1",
					"ip_addresses": ["1.2.3.4", "5.6.7.8"],
					"dns_names": ["us-east.monitor.example.com"]
				},
				{
					"id": "region456",
					"name": "EU West",
					"code": "eu-west-1",
					"ip_addresses": ["10.20.30.40"],
					"dns_names": ["eu-west.monitor.example.com"]
				}
			]
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	regions, err := client.Regions.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, regions, 2)

	assert.Equal(t, "region123", regions[0].ID)
	assert.Equal(t, "US East", regions[0].Name)
	assert.Equal(t, "us-east-1", regions[0].Code)
	assert.Equal(t, []string{"1.2.3.4", "5.6.7.8"}, regions[0].IPAddrs)
	assert.Equal(t, []string{"us-east.monitor.example.com"}, regions[0].DNSNames)

	assert.Equal(t, "EU West", regions[1].Name)
}
