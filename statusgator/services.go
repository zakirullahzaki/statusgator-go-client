package statusgator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ServicesService handles external service catalog operations.
// Note: Some operations require Firehose access.
type ServicesService service

type servicesResponse struct {
	Success    bool       `json:"success"`
	Data       []Service  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type serviceComponentsResponse struct {
	Success    bool        `json:"success"`
	Data       []Component `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// List returns all available services with pagination.
// Requires Firehose access.
func (s *ServicesService) List(ctx context.Context, opts *ListOptions) ([]Service, *Pagination, error) {
	path := "/services"
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp servicesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAll returns all available services across all pages.
// Requires Firehose access.
func (s *ServicesService) ListAll(ctx context.Context) ([]Service, error) {
	var allServices []Service
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		services, pagination, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allServices = append(allServices, services...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allServices, nil
}

// Search searches for services by query.
func (s *ServicesService) Search(ctx context.Context, query string) ([]Service, error) {
	params := url.Values{}
	params.Set("query", query)

	path := addQueryParams("/services/search", params)

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp servicesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, nil
}

// ListComponents returns components for a service with pagination.
func (s *ServicesService) ListComponents(ctx context.Context, serviceID string, opts *ListOptions) ([]Component, *Pagination, error) {
	if err := validateID(serviceID, "serviceID"); err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/services/%s/components", url.PathEscape(serviceID))
	if opts != nil {
		path = addQueryParams(path, opts.Values())
	}

	body, err := s.client.get(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	var resp serviceComponentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %w", err)
	}

	return resp.Data, &resp.Pagination, nil
}

// ListAllComponents returns all components for a service across all pages.
func (s *ServicesService) ListAllComponents(ctx context.Context, serviceID string) ([]Component, error) {
	var allComponents []Component
	opts := &ListOptions{Page: 1, PerPage: MaxPerPage}

	for {
		components, pagination, err := s.ListComponents(ctx, serviceID, opts)
		if err != nil {
			return nil, err
		}

		allComponents = append(allComponents, components...)

		if !pagination.HasNextPage() {
			break
		}
		opts.Page++
	}

	return allComponents, nil
}
