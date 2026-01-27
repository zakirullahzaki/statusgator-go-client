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
					"region_id": "us-east",
					"dns_name": "monitor-us-east.statusgator.com",
					"ip_address": "34.230.98.138",
					"provider": "AWS",
					"code": "us-east-1",
					"name": "US East",
					"desc": "US East (N. Virginia)",
					"icon_url": "https://example.com/aws.png",
					"color": "#1a5db2"
				},
				{
					"region_id": "eu-west",
					"dns_name": "monitor-eu-west.statusgator.com",
					"ip_address": "34.250.135.124",
					"provider": "AWS",
					"code": "eu-west-1",
					"name": "EU West",
					"desc": "Europe (Ireland)",
					"icon_url": "https://example.com/aws.png",
					"color": "#a66f00"
				}
			]
		}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient("test-token", WithBaseURL(server.URL))
	require.NoError(t, err)

	regions, err := client.Regions.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, regions, 2)

	assert.Equal(t, "us-east", regions[0].RegionID)
	assert.Equal(t, "US East", regions[0].Name)
	assert.Equal(t, "us-east-1", regions[0].Code)
	assert.Equal(t, "US East (N. Virginia)", regions[0].Desc)
	assert.Equal(t, "AWS", regions[0].Provider)
	assert.Equal(t, "monitor-us-east.statusgator.com", regions[0].DNSName)
	assert.Equal(t, "34.230.98.138", regions[0].IPAddress)
	assert.Equal(t, "#1a5db2", regions[0].Color)

	assert.Equal(t, "EU West", regions[1].Name)
}
