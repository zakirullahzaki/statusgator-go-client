package statusgator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL         = "https://statusgator.com/api/v3"
	DefaultTimeout         = 30 * time.Second
	DefaultMaxResponseSize = 10 * 1024 * 1024 // 10MB
	DefaultUserAgent       = "statusgator-go-client/1.0.0"
)

// Client manages communication with the StatusGator API.
type Client struct {
	client    *http.Client
	BaseURL   string
	Token     string
	UserAgent string

	common service // Reuse a single struct for all services

	// Services
	Boards          *BoardsService
	Monitors        *MonitorsService
	WebsiteMonitors *WebsiteMonitorsService
	PingMonitors    *PingMonitorsService
	ServiceMonitors *ServiceMonitorsService
	CustomMonitors  *CustomMonitorsService
	MonitorGroups   *MonitorGroupsService
	Components      *ComponentsService
	Incidents       *IncidentsService
	Services        *ServicesService
	Subscribers     *StatusPageSubscribersService
	Users           *UsersService
	Regions         *RegionsService
}

type service struct {
	client *Client
}

// ClientOption configures the client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.client = httpClient
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithUserAgent sets a custom user agent.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// WithTimeout sets a custom timeout for the HTTP client.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// NewClient creates a new StatusGator API client.
func NewClient(token string, opts ...ClientOption) (*Client, error) {
	if token == "" {
		return nil, ErrTokenRequired
	}

	c := &Client{
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
		BaseURL:   DefaultBaseURL,
		Token:     token,
		UserAgent: DefaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.common.client = c
	c.Boards = (*BoardsService)(&c.common)
	c.Monitors = (*MonitorsService)(&c.common)
	c.WebsiteMonitors = (*WebsiteMonitorsService)(&c.common)
	c.PingMonitors = (*PingMonitorsService)(&c.common)
	c.ServiceMonitors = (*ServiceMonitorsService)(&c.common)
	c.CustomMonitors = (*CustomMonitorsService)(&c.common)
	c.MonitorGroups = (*MonitorGroupsService)(&c.common)
	c.Components = (*ComponentsService)(&c.common)
	c.Incidents = (*IncidentsService)(&c.common)
	c.Services = (*ServicesService)(&c.common)
	c.Subscribers = (*StatusPageSubscribersService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.Regions = (*RegionsService)(&c.common)

	return c, nil
}

// Ping verifies API connectivity and authentication.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.get(ctx, "/ping")
	return err
}

// doRequest executes an HTTP request and returns the response body.
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	reqURL := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	c.setHeaders(req, body != nil)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	limitedReader := io.LimitReader(resp.Body, DefaultMaxResponseSize+1)
	respBody, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if len(respBody) > DefaultMaxResponseSize {
		return nil, ErrResponseTooLarge
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

func (c *Client) setHeaders(req *http.Request, hasBody bool) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

func (c *Client) post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithJSON(ctx, http.MethodPost, path, body)
}

func (c *Client) patch(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithJSON(ctx, http.MethodPatch, path, body)
}

func (c *Client) delete(ctx context.Context, path string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

func (c *Client) deleteWithQuery(ctx context.Context, path string, params url.Values) error {
	fullPath := addQueryParams(path, params)
	_, err := c.doRequest(ctx, http.MethodDelete, fullPath, nil)
	return err
}

func (c *Client) doRequestWithJSON(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
	}
	return c.doRequest(ctx, method, path, &buf)
}
