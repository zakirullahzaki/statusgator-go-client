package statusgator

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getE2EClient returns a client for E2E tests.
// Skips the test if STATUSGATOR_API_TOKEN is not set.
func getE2EClient(t *testing.T) *Client {
	token := os.Getenv("STATUSGATOR_API_TOKEN")
	if token == "" {
		t.Skip("STATUSGATOR_API_TOKEN not set, skipping E2E test")
	}

	client, err := NewClient(token)
	require.NoError(t, err)
	return client
}

func TestE2E_Ping(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	err := client.Ping(ctx)
	assert.NoError(t, err, "Ping should succeed with valid token")
}

func TestE2E_Regions(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	regions, err := client.Regions.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, regions, "Should have at least one region")

	// Verify first region has all expected fields
	region := regions[0]
	assert.NotEmpty(t, region.RegionID, "region_id should not be empty")
	assert.NotEmpty(t, region.Name, "name should not be empty")
	assert.NotEmpty(t, region.Code, "code should not be empty")
	assert.NotEmpty(t, region.Provider, "provider should not be empty")
	assert.NotEmpty(t, region.DNSName, "dns_name should not be empty")
	assert.NotEmpty(t, region.IPAddress, "ip_address should not be empty")

	t.Logf("Found %d regions, first: %s (%s)", len(regions), region.Name, region.Code)
}

func TestE2E_Users(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	users, err := client.Users.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, users, "Should have at least one user")

	// Verify first user has all expected fields
	user := users[0]
	assert.NotEmpty(t, user.ID, "id should not be empty")
	assert.NotEmpty(t, user.Email, "email should not be empty")
	assert.NotEmpty(t, user.Role, "role should not be empty")
	assert.NotEmpty(t, user.CreatedAt, "created_at should not be empty")

	// Test FullName helper
	fullName := user.FullName()
	assert.NotEmpty(t, fullName, "FullName() should return non-empty string")

	t.Logf("Found %d users, first: %s <%s>", len(users), fullName, user.Email)
}

func TestE2E_Boards(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	boards, pagination, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, pagination, "pagination should not be nil")

	if len(boards) == 0 {
		t.Log("No boards found, skipping board field validation")
		return
	}

	// Verify first board has all expected fields
	board := boards[0]
	assert.NotEmpty(t, board.ID, "id should not be empty")
	assert.NotEmpty(t, board.Name, "name should not be empty")
	assert.NotEmpty(t, board.PublicToken, "public_token should not be empty")
	assert.NotEmpty(t, board.CreatedAt, "created_at should not be empty")

	t.Logf("Found %d boards, first: %s (ID: %s)", len(boards), board.Name, board.ID)

	// Test Get single board
	singleBoard, err := client.Boards.Get(ctx, board.ID)
	require.NoError(t, err)
	assert.Equal(t, board.ID, singleBoard.ID)
	assert.Equal(t, board.Name, singleBoard.Name)
}

func TestE2E_Monitors(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	// First get a board
	boards, _, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)

	if len(boards) == 0 {
		t.Skip("No boards found, skipping monitors test")
	}

	boardID := boards[0].ID

	monitors, pagination, err := client.Monitors.List(ctx, boardID, nil)
	require.NoError(t, err)
	require.NotNil(t, pagination, "pagination should not be nil")

	if len(monitors) == 0 {
		t.Log("No monitors found, skipping monitor field validation")
		return
	}

	// Verify first monitor has all expected fields from real API
	monitor := monitors[0]
	assert.NotEmpty(t, monitor.ID, "id should not be empty")
	assert.NotEmpty(t, monitor.DisplayName, "display_name should not be empty")
	assert.NotEmpty(t, monitor.MonitorType, "monitor_type should not be empty")
	assert.NotEmpty(t, monitor.FilteredStatus, "filtered_status should not be empty")
	assert.NotEmpty(t, monitor.IconURL, "icon_url should not be empty")
	assert.NotEmpty(t, monitor.CreatedAt, "created_at should not be empty")

	// Verify monitor_type is one of known types
	validTypes := []MonitorType{MonitorTypeWebsite, MonitorTypePing, MonitorTypeService, MonitorTypeCustom}
	assert.Contains(t, validTypes, monitor.MonitorType, "monitor_type should be a valid type")

	// Verify status is one of known statuses
	validStatuses := []MonitorStatus{MonitorStatusUp, MonitorStatusDown, MonitorStatusWarn, MonitorStatusMaintenance, MonitorStatusUnknown}
	assert.Contains(t, validStatuses, monitor.FilteredStatus, "filtered_status should be a valid status")
	assert.Contains(t, validStatuses, monitor.UnfilteredStatus, "unfiltered_status should be a valid status")

	// Test IsPaused helper
	_ = monitor.IsPaused() // Just verify it doesn't panic

	t.Logf("Found %d monitors, first: %s [%s] - %s",
		len(monitors), monitor.DisplayName, monitor.MonitorType, monitor.FilteredStatus)

	// If it's a ServiceMonitor, verify Service nested object
	if monitor.MonitorType == MonitorTypeService && monitor.Service != nil {
		assert.NotEmpty(t, monitor.Service.ID, "service.id should not be empty")
		assert.NotEmpty(t, monitor.Service.Name, "service.name should not be empty")
		t.Logf("  Service: %s (slug: %s)", monitor.Service.Name, monitor.Service.Slug)
	}

	// If it has a group, verify Group nested object
	if monitor.Group != nil {
		assert.NotEmpty(t, monitor.Group.ID, "group.id should not be empty")
		assert.NotEmpty(t, monitor.Group.Name, "group.name should not be empty")
		t.Logf("  Group: %s", monitor.Group.Name)
	}
}

func TestE2E_MonitorGroups(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	// First get a board
	boards, _, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)

	if len(boards) == 0 {
		t.Skip("No boards found, skipping monitor groups test")
	}

	boardID := boards[0].ID

	groups, err := client.MonitorGroups.List(ctx, boardID)
	require.NoError(t, err)

	if len(groups) == 0 {
		t.Log("No monitor groups found")
		return
	}

	// Verify first group has all expected fields
	group := groups[0]
	assert.NotEmpty(t, group.ID, "id should not be empty")
	assert.NotEmpty(t, group.Name, "name should not be empty")
	assert.NotEmpty(t, group.CreatedAt, "created_at should not be empty")

	t.Logf("Found %d monitor groups, first: %s (position: %d)", len(groups), group.Name, group.Position)
}

func TestE2E_BoardHistory(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	// First get a board
	boards, _, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)

	if len(boards) == 0 {
		t.Skip("No boards found, skipping board history test")
	}

	boardID := boards[0].ID

	history, err := client.Boards.GetHistory(ctx, boardID, nil)
	require.NoError(t, err)

	if len(history) == 0 {
		t.Log("No history events found")
		return
	}

	// Verify first history event has all expected fields from real API
	event := history[0]
	assert.NotEmpty(t, event.MonitorID, "monitor_id should not be empty")
	assert.NotEmpty(t, event.Name, "name should not be empty")
	assert.NotEmpty(t, event.Status, "status should not be empty")
	assert.NotEmpty(t, event.StartedAt, "started_at should not be empty")

	t.Logf("Found %d history events, first: %s [%s] at %s",
		len(history), event.Name, event.Status, event.StartedAt)
}

func TestE2E_Incidents(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	// First get a board
	boards, _, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)

	if len(boards) == 0 {
		t.Skip("No boards found, skipping incidents test")
	}

	boardID := boards[0].ID

	incidents, pagination, err := client.Incidents.List(ctx, boardID, nil)
	require.NoError(t, err)
	require.NotNil(t, pagination, "pagination should not be nil")

	if len(incidents) == 0 {
		t.Log("No incidents found (this is normal for new accounts)")
		return
	}

	// Verify first incident has all expected fields from real API
	incident := incidents[0]
	assert.NotEmpty(t, incident.ID, "id should not be empty")
	assert.NotEmpty(t, incident.Name, "name should not be empty")
	assert.NotEmpty(t, incident.Severity, "severity should not be empty")
	assert.NotEmpty(t, incident.Phase, "phase should not be empty")
	assert.NotEmpty(t, incident.BoardID, "board_id should not be empty")
	assert.NotEmpty(t, incident.CreatedAt, "created_at should not be empty")

	t.Logf("Found %d incidents, first: %s [%s - %s]",
		len(incidents), incident.Name, incident.Severity, incident.Phase)
}

func TestE2E_Subscribers(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	// First get a board
	boards, _, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)

	if len(boards) == 0 {
		t.Skip("No boards found, skipping subscribers test")
	}

	boardID := boards[0].ID

	subscribers, pagination, err := client.Subscribers.List(ctx, boardID, nil)
	require.NoError(t, err)
	require.NotNil(t, pagination, "pagination should not be nil")

	if len(subscribers) == 0 {
		t.Log("No subscribers found (this is normal for new accounts)")
		return
	}

	// Verify first subscriber has all expected fields from real API
	subscriber := subscribers[0]
	assert.NotEmpty(t, subscriber.ID, "id should not be empty")
	assert.NotEmpty(t, subscriber.Email, "email should not be empty")
	assert.NotEmpty(t, subscriber.CreatedAt, "created_at should not be empty")

	t.Logf("Found %d subscribers, first: %s (confirmed: %v)",
		len(subscribers), subscriber.Email, subscriber.Confirmed)
}

func TestE2E_Components(t *testing.T) {
	client := getE2EClient(t)
	ctx := context.Background()

	// First get a board
	boards, _, err := client.Boards.List(ctx, nil)
	require.NoError(t, err)

	if len(boards) == 0 {
		t.Skip("No boards found, skipping components test")
	}

	boardID := boards[0].ID

	// Get monitors for the board
	monitors, _, err := client.Monitors.List(ctx, boardID, nil)
	require.NoError(t, err)

	// Find a ServiceMonitor (it has components)
	var serviceMonitorID string
	for _, m := range monitors {
		if m.MonitorType == MonitorTypeService {
			serviceMonitorID = m.ID
			break
		}
	}

	if serviceMonitorID == "" {
		t.Skip("No ServiceMonitor found, skipping components test")
	}

	components, pagination, err := client.Components.ListForMonitor(ctx, boardID, serviceMonitorID, nil)
	require.NoError(t, err)
	require.NotNil(t, pagination, "pagination should not be nil")

	if len(components) == 0 {
		t.Log("No components found")
		return
	}

	// Verify first component has all expected fields from real API
	component := components[0]
	assert.NotEmpty(t, component.ID, "id should not be empty")
	assert.NotEmpty(t, component.Name, "name should not be empty")
	assert.NotEmpty(t, component.ServiceID, "service_id should not be empty")
	assert.NotEmpty(t, component.Status, "status should not be empty")
	assert.NotEmpty(t, component.CreatedAt, "created_at should not be empty")

	t.Logf("Found %d components, first: %s [%s]",
		len(components), component.Name, component.Status)
}
