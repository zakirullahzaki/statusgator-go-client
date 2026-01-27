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
	MonitorTypeWebsite MonitorType = "website"
	MonitorTypePing    MonitorType = "ping"
	MonitorTypeService MonitorType = "service"
	MonitorTypeCustom  MonitorType = "custom"
)

// Board represents a StatusGator dashboard.
type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PublicToken string    `json:"public_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Monitor represents a base monitor structure.
type Monitor struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Type      MonitorType   `json:"type"`
	Status    MonitorStatus `json:"status"`
	Paused    bool          `json:"paused"`
	GroupID   *string       `json:"group_id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// WebsiteMonitor represents a website HTTP monitor.
type WebsiteMonitor struct {
	Monitor
	URL             string            `json:"url"`
	CheckInterval   int               `json:"check_interval"`
	HTTPMethod      string            `json:"http_method"`
	ExpectedStatus  int               `json:"expected_status"`
	ContentMatch    string            `json:"content_match,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	BasicAuthUser   string            `json:"basic_auth_user,omitempty"`
	BasicAuthPass   string            `json:"basic_auth_pass,omitempty"`
	Timeout         int               `json:"timeout"`
	FollowRedirects bool              `json:"follow_redirects"`
	Regions         []string          `json:"regions,omitempty"`
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
	Host          string   `json:"host"`
	CheckInterval int      `json:"check_interval"`
	Regions       []string `json:"regions,omitempty"`
}

// PingMonitorRequest represents a request to create/update a ping monitor.
type PingMonitorRequest struct {
	Name          string   `json:"name,omitempty"`
	Host          string   `json:"host,omitempty"`
	CheckInterval int      `json:"check_interval,omitempty"`
	Regions       []string `json:"regions,omitempty"`
	GroupID       string   `json:"group_id,omitempty"`
}

// ServiceMonitor represents a subscription to an external status page.
type ServiceMonitor struct {
	Monitor
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
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
	Description string `json:"description,omitempty"`
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
type Incident struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Severity     IncidentSeverity `json:"severity"`
	Phase        IncidentPhase    `json:"phase"`
	MonitorIDs   []string         `json:"monitor_ids"`
	Updates      []IncidentUpdate `json:"updates,omitempty"`
	ScheduledFor *time.Time       `json:"scheduled_for,omitempty"`
	ScheduledEnd *time.Time       `json:"scheduled_end,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// IncidentUpdate represents an update to an incident.
type IncidentUpdate struct {
	ID        string        `json:"id"`
	Message   string        `json:"message"`
	Phase     IncidentPhase `json:"phase"`
	CreatedAt time.Time     `json:"created_at"`
}

// IncidentRequest represents a request to create an incident.
type IncidentRequest struct {
	Title        string           `json:"title"`
	Message      string           `json:"message"`
	Severity     IncidentSeverity `json:"severity"`
	Phase        IncidentPhase    `json:"phase,omitempty"`
	MonitorIDs   []string         `json:"monitor_ids"`
	ScheduledFor *time.Time       `json:"scheduled_for,omitempty"`
	ScheduledEnd *time.Time       `json:"scheduled_end,omitempty"`
}

// IncidentUpdateRequest represents a request to update an incident.
type IncidentUpdateRequest struct {
	Message string        `json:"message"`
	Phase   IncidentPhase `json:"phase"`
}

// Service represents an external service that can be monitored.
type Service struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	StatusURL string    `json:"status_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Subscriber represents a status page email subscriber.
type Subscriber struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Confirmed bool      `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
}

// SubscriberRequest represents a request to add a subscriber.
type SubscriberRequest struct {
	Email            string `json:"email"`
	SkipConfirmation bool   `json:"skip_confirmation,omitempty"`
}

// User represents an organization user.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Region represents a geographic monitoring region.
type Region struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	IPAddrs  []string `json:"ip_addresses"`
	DNSNames []string `json:"dns_names"`
}

// HistoryEvent represents a historical status event.
type HistoryEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	MonitorID string    `json:"monitor_id"`
	Details   string    `json:"details"`
}

// HistoryOptions specifies filters for history queries.
type HistoryOptions struct {
	StartDate string
	EndDate   string
	MonitorID string
}
