package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	maxRetries       = 3
	defaultRetryWait = 10 * time.Second
	maxRetryWait     = 60 * time.Second
)

const (
	NoProjectID  = ""
	infraAPIPath = "/v1/mgmt/infra"
)

type Response struct {
	Entity string         `json:"entity"`
	ID     string         `json:"id"`
	Data   map[string]any `json:"data"`
}

type Client struct {
	version       string
	managementKey string
	baseURL       string

	apiClients map[string]*api.Client
	lock       sync.Mutex
}

func NewClient(version, managementKey, baseURL string) *Client {
	return &Client{
		version:       version,
		managementKey: managementKey,
		baseURL:       baseURL,
		apiClients:    map[string]*api.Client{},
	}
}

func (c *Client) Create(ctx context.Context, projectID, entity string, data map[string]any) (*Response, error) {
	httpBody := map[string]any{
		"entity": entity,
		"data":   data,
	}

	tflog.Info(ctx, "Starting CREATE request", map[string]any{"body": debugRequest(httpBody)})
	httpRes, err := retryOnRateLimit(ctx, func() (*api.HTTPResponse, error) {
		return c.getAPIClient(projectID).DoPostRequest(ctx, infraAPIPath, httpBody, nil, c.managementKey)
	})
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished CREATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Read(ctx context.Context, projectID, entity, entityID string) (*Response, error) {
	httpQuery := map[string]string{
		"entity": entity,
		"id":     entityID,
	}

	tflog.Info(ctx, "Starting READ request", map[string]any{"query": debugRequest(httpQuery)})
	httpRes, err := retryOnRateLimit(ctx, func() (*api.HTTPResponse, error) {
		return c.getAPIClient(projectID).DoGetRequest(ctx, infraAPIPath, &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey)
	})
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished READ request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Update(ctx context.Context, projectID, entity, entityID string, data map[string]any) (*Response, error) {
	httpBody := map[string]any{
		"entity": entity,
		"id":     entityID,
		"data":   data,
	}

	tflog.Info(ctx, "Starting UPDATE request", map[string]any{"body": debugRequest(httpBody)})
	httpRes, err := retryOnRateLimit(ctx, func() (*api.HTTPResponse, error) {
		return c.getAPIClient(projectID).DoPutRequest(ctx, infraAPIPath, httpBody, nil, c.managementKey)
	})
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished UPDATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Delete(ctx context.Context, projectID, entity, entityID string) error {
	httpQuery := map[string]string{
		"entity": entity,
		"id":     entityID,
	}

	tflog.Info(ctx, "Starting DELETE request", map[string]any{"query": debugRequest(httpQuery)})
	httpRes, err := retryOnRateLimit(ctx, func() (*api.HTTPResponse, error) {
		return c.getAPIClient(projectID).DoDeleteRequest(ctx, infraAPIPath, &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey)
	})
	if err != nil {
		return err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return err
	}

	tflog.Info(ctx, "Finished DELETE request")
	return nil
}

func (c *Client) getAPIClient(projectID string) *api.Client {
	c.lock.Lock()
	defer c.lock.Unlock()

	apiClient, ok := c.apiClients[projectID]
	if !ok {
		apiClient = makeAPIClient(c.version, projectID, c.baseURL)
		c.apiClients[projectID] = apiClient
	}

	return apiClient
}

func makeAPIClient(version, projectID, baseURL string) *api.Client {
	headers := map[string]string{
		"user-agent": makeUserAgent(version),
	}

	params := api.ClientParams{
		ProjectID:            projectID,
		BaseURL:              baseURL,
		CustomDefaultHeaders: headers,
	}

	return api.NewClient(params)
}

func makeUserAgent(version string) string {
	if v := os.Getenv("DESCOPE_USER_AGENT"); v != "" {
		return v
	}
	return fmt.Sprintf("terraform-provider-descope/%s", version)
}

func retryOnRateLimit(ctx context.Context, fn func() (*api.HTTPResponse, error)) (*api.HTTPResponse, error) {
	for attempt := range maxRetries {
		res, err := fn()
		if err == nil {
			return res, nil
		}

		de := descope.AsError(err, descope.ErrRateLimitExceeded.Code)
		if de == nil {
			return nil, err
		}

		wait := defaultRetryWait
		if retryAfter, ok := de.Info[descope.ErrorInfoKeys.RateLimitExceededRetryAfter].(int); ok && retryAfter > 0 {
			wait = time.Duration(retryAfter) * time.Second
		}
		if wait > maxRetryWait {
			wait = maxRetryWait
		}

		tflog.Warn(ctx, "Rate limited by Descope API, retrying", map[string]any{
			"attempt": attempt + 1,
			"max":     maxRetries,
			"wait":    wait.String(),
		})

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return fn()
}
