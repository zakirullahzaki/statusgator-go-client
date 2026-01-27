package statusgator

import "time"

// MonitorStatus represents the status of a monitor.
type MonitorStatus string

const (
	MonitorStatusUp          MonitorStatus = "up"
	MonitorStatusDown        MonitorStatus = "down"
	MonitorStatusWarn        MonitorStatus = "warn"
	MonitorStatusMaintenance MonitorStatus = "maintenance"
	MonitorStatusUnknown     MonitorStatus = "unknown"
)

// MonitorType represents the type of monitor.
type MonitorType string

const (
	MonitorTypeWebsite MonitorType = "WebsiteMonitor"
	MonitorTypePing    MonitorType = "PingMonitor"
	MonitorTypeService MonitorType = "ServiceMonitor"
	MonitorTypeCustom  MonitorType = "CustomMonitor"
)

// Board represents a StatusGator dashboard.
type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PublicToken string    `json:"public_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GroupInfo represents nested group information in a monitor response.
type GroupInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

// Monitor represents a monitor from the API v3 /boards/{id}/monitors endpoint.
// Field names match the actual API response exactly.
type Monitor struct {
	ID                 string         `json:"id"`
	DisplayName        string         `json:"display_name"`
	MonitorType        MonitorType    `json:"monitor_type"`
	FilteredStatus     MonitorStatus  `json:"filtered_status"`
	UnfilteredStatus   MonitorStatus  `json:"unfiltered_status"`
	Description        *string        `json:"description,omitempty"`
	LastMessage        *string        `json:"last_message,omitempty"`
	LastDetails        *string        `json:"last_details,omitempty"`
	OverriddenMessage  *string        `json:"overridden_message,omitempty"`
	OverriddenStatus   *MonitorStatus `json:"overridden_status,omitempty"`
	OverridesLockedAt  *time.Time     `json:"overrides_locked_at,omitempty"`
	PausedAt           *time.Time     `json:"paused_at,omitempty"`
	CheckedAt          *time.Time     `json:"checked_at,omitempty"`
	FilterCount        int            `json:"filter_count"`
	IconURL            string         `json:"icon_url"`
	Position           *int           `json:"position,omitempty"`
	EarlyWarningSignal bool           `json:"early_warning_signal"`
	Service            *ServiceInfo   `json:"service,omitempty"`
	Group              *GroupInfo     `json:"group,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// IsPaused returns true if the monitor is currently paused.
func (m *Monitor) IsPaused() bool {
	return m.PausedAt != nil
}

// WebsiteMonitor represents a website HTTP monitor.
type WebsiteMonitor struct {
	Monitor
	URL              string   `json:"url"`
	CheckInterval    int      `json:"check_interval"`
	HTTPMethod       string   `json:"http_method"`
	CheckContent     bool     `json:"check_content"`
	Content          string   `json:"content,omitempty"`
	AlertContentFound bool    `json:"alert_content_found"`
	CheckRegions     []string `json:"check_regions,omitempty"`
	AlertAnyLocation bool     `json:"alert_any_location"`
	ResponseCodes    []int    `json:"response_codes,omitempty"`
	FollowRedirects  bool     `json:"follow_redirects"`
	Timeout          int      `json:"timeout"`
	RetryCount       int      `json:"retry_count"`
	RequestBody      string   `json:"request_body,omitempty"`
	HTTPAuthUsername string   `json:"http_auth_username,omitempty"`
	HTTPAuthPassword string   `json:"http_auth_password,omitempty"`
	RequestHeaders   []string `json:"request_headers,omitempty"`
}

// WebsiteMonitorRequest represents a request to create/update a website monitor.
type WebsiteMonitorRequest struct {
	Name            string            `json:"name,omitempty"`
	URL             string            `json:"url,omitempty"`
	CheckInterval   int               `json:"check_interval,omitempty"`
	HTTPMethod      string            `json:"http_method,omitempty"`
	ExpectedStatus  int               `json:"expected_status,omitempty"`
	ContentMatch    string            `json:"content_match,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	BasicAuthUser   string            `json:"basic_auth_user,omitempty"`
	BasicAuthPass   string            `json:"basic_auth_pass,omitempty"`
	Timeout         int               `json:"timeout,omitempty"`
	FollowRedirects *bool             `json:"follow_redirects,omitempty"`
	Regions         []string          `json:"regions,omitempty"`
	GroupID         string            `json:"group_id,omitempty"`
}

// PingMonitor represents a ping/ICMP monitor.
type PingMonitor struct {
	Monitor
	Address  string   `json:"address"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
	Regions  []string `json:"regions,omitempty"`
}

// PingMonitorRequest represents a request to create/update a ping monitor.
type PingMonitorRequest struct {
	Name     string   `json:"name,omitempty"`
	Address  string   `json:"address,omitempty"`
	Interval int      `json:"interval,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
	Regions  []string `json:"regions,omitempty"`
	GroupID  string   `json:"group_id,omitempty"`
}

// ServiceInfo represents nested service information in a monitor response.
type ServiceInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	HomePageURL    string `json:"home_page_url"`
	StatusPageURL  string `json:"status_page_url"`
	IconURL        string `json:"icon_url"`
	LandingPageURL string `json:"landing_page_url"`
	Official       bool   `json:"official"`
}

// ServiceMonitor represents a subscription to an external status page.
type ServiceMonitor struct {
	Monitor
	Service *ServiceInfo `json:"service,omitempty"`
}

// GetServiceID returns the service ID from the nested service object.
func (sm *ServiceMonitor) GetServiceID() string {
	if sm.Service != nil {
		return sm.Service.ID
	}
	return ""
}

// GetServiceName returns the service name from the nested service object.
func (sm *ServiceMonitor) GetServiceName() string {
	if sm.Service != nil {
		return sm.Service.Name
	}
	return ""
}

// ServiceMonitorRequest represents a request to create/update a service monitor.
type ServiceMonitorRequest struct {
	Name      string `json:"name,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
}

// CustomMonitor represents a manually-managed monitor.
type CustomMonitor struct {
	Monitor
}

// CustomMonitorRequest represents a request to create/update a custom monitor.
type CustomMonitorRequest struct {
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Status      MonitorStatus `json:"status,omitempty"`
	GroupID     string        `json:"group_id,omitempty"`
}

// MonitorGroup represents a group of monitors.
type MonitorGroup struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	Collapsed bool      `json:"collapsed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MonitorGroupRequest represents a request to create/update a monitor group.
type MonitorGroupRequest struct {
	Name      string `json:"name,omitempty"`
	Position  int    `json:"position,omitempty"`
	Collapsed *bool  `json:"collapsed,omitempty"`
}

// Component represents a service component.
type Component struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	GroupName string        `json:"group_name"`
	ServiceID string        `json:"service_id"`
	Status    MonitorStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// IncidentSeverity represents incident severity level.
type IncidentSeverity string

const (
	IncidentSeverityMinor       IncidentSeverity = "minor"
	IncidentSeverityMajor       IncidentSeverity = "major"
	IncidentSeverityMaintenance IncidentSeverity = "maintenance"
)

// IncidentPhase represents the phase of an incident.
type IncidentPhase string

const (
	IncidentPhaseInvestigating IncidentPhase = "investigating"
	IncidentPhaseIdentified    IncidentPhase = "identified"
	IncidentPhaseMonitoring    IncidentPhase = "monitoring"
	IncidentPhaseResolved      IncidentPhase = "resolved"
	IncidentPhaseScheduled     IncidentPhase = "scheduled"
	IncidentPhaseInProgress    IncidentPhase = "in_progress"
	IncidentPhaseVerifying     IncidentPhase = "verifying"
	IncidentPhaseCompleted     IncidentPhase = "completed"
)

// Incident represents an incident or maintenance window.
// Matches the actual API response from /boards/{id}/incidents endpoint.
type Incident struct {
	ID                      string           `json:"id"`
	Name                    string           `json:"name"`
	Details                 string           `json:"details"`
	Severity                IncidentSeverity `json:"severity"`
	Phase                   IncidentPhase    `json:"phase"`
	StartedAt               *time.Time       `json:"started_at,omitempty"`
	ResolvedAt              *time.Time       `json:"resolved_at,omitempty"`
	WillStartAt             *time.Time       `json:"will_start_at,omitempty"`
	WillEndAt               *time.Time       `json:"will_end_at,omitempty"`
	AutoCompleteMaintenance bool             `json:"auto_complete_maintenance"`
	BoardID                 string           `json:"board_id"`
	Duration                *string          `json:"duration,omitempty"`
	MaintenanceDuration     *string          `json:"maintenance_duration,omitempty"`
	ScheduledMaintenance    bool             `json:"scheduled_maintenance"`
	ResolvedOrCompleted     bool             `json:"resolved_or_completed"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
}

// IncidentUpdate represents an update to an incident.
// Matches the actual API response from incident_updates endpoint.
type IncidentUpdate struct {
	ID                string           `json:"id"`
	IncidentID        string           `json:"incident_id"`
	Details           string           `json:"details"`
	Phase             IncidentPhase    `json:"phase"`
	Severity          IncidentSeverity `json:"severity"`
	PostedAt          *time.Time       `json:"posted_at,omitempty"`
	NotifySubscribers bool             `json:"notify_subscribers"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// IncidentRequest represents a request to create an incident.
type IncidentRequest struct {
	Name                    string           `json:"name"`
	Details                 string           `json:"details"`
	Severity                IncidentSeverity `json:"severity"`
	Phase                   IncidentPhase    `json:"phase,omitempty"`
	WillStartAt             *time.Time       `json:"will_start_at,omitempty"`
	WillEndAt               *time.Time       `json:"will_end_at,omitempty"`
	AutoCompleteMaintenance *bool            `json:"auto_complete_maintenance,omitempty"`
}

// IncidentUpdateRequest represents a request to add an update to an incident.
type IncidentUpdateRequest struct {
	Details           string           `json:"details"`
	Phase             IncidentPhase    `json:"phase,omitempty"`
	Severity          IncidentSeverity `json:"severity,omitempty"`
	NotifySubscribers *bool            `json:"notify_subscribers,omitempty"`
}

// Service represents an external service that can be monitored.
type Service struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	HomePageURL    string    `json:"home_page_url"`
	StatusPageURL  string    `json:"status_page_url"`
	IconURL        string    `json:"icon_url"`
	LandingPageURL string    `json:"landing_page_url"`
	Official       bool      `json:"official"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Subscriber represents a status page email subscriber.
// Matches the actual API response from status_page_subscribers endpoint.
type Subscriber struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Confirmed   bool       `json:"confirmed"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	MonitorIDs  []string   `json:"monitor_ids,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SubscriberRequest represents a request to add a subscriber.
type SubscriberRequest struct {
	Email            string `json:"email"`
	SkipConfirmation bool   `json:"skip_confirmation,omitempty"`
}

// User represents an organization user.
type User struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Company          string     `json:"company"`
	JobTitle         *string    `json:"job_title,omitempty"`
	Role             string     `json:"role"`
	Confirmed        bool       `json:"confirmed"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	CreatedAt        time.Time  `json:"created_at"`
	LastSignInAt     *time.Time `json:"last_sign_in_at,omitempty"`
}

// FullName returns user's full name.
func (u *User) FullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// Region represents a geographic monitoring region.
type Region struct {
	RegionID  string `json:"region_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Desc      string `json:"desc"`
	Provider  string `json:"provider"`
	DNSName   string `json:"dns_name"`
	IPAddress string `json:"ip_address"`
	IconURL   string `json:"icon_url"`
	Color     string `json:"color"`
}

// HistoryEvent represents a historical status event from board history.
// Matches the actual API response from /boards/{id}/history endpoint.
type HistoryEvent struct {
	MonitorID          string        `json:"monitor_id"`
	Name               string        `json:"name"`
	IconURL            string        `json:"icon_url"`
	Status             MonitorStatus `json:"status"`
	StartedAt          time.Time     `json:"started_at"`
	EndedAt            *time.Time    `json:"ended_at,omitempty"`
	Duration           string        `json:"duration"`
	Message            string        `json:"message"`
	Details            string        `json:"details"`
	EarlyWarningSignal bool          `json:"early_warning_signal"`
}

// HistoryOptions specifies filters for history queries.
type HistoryOptions struct {
	StartDate string
	EndDate   string
	MonitorID string
}
