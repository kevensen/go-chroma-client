package chromaclient

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

// Client is a ChromaDB client
type Client struct {
	baseURL    string
	httpClient *http.Client
	tenant     string
	database   string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTenant sets the default tenant for the client
func WithTenant(tenant string) ClientOption {
	return func(c *Client) {
		c.tenant = tenant
	}
}

// WithDatabase sets the default database for the client
func WithDatabase(database string) ClientOption {
	return func(c *Client) {
		c.database = database
	}
}

// NewClient creates a new ChromaDB client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: "http://localhost:8000",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		tenant:   DefaultTenant,
		database: DefaultDatabase,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
			Timestamp:  time.Now(),
		}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Version returns the ChromaDB version
func (c *Client) Version(ctx context.Context) (string, error) {
	var version string
	err := c.doRequest(ctx, http.MethodGet, "/api/v2/version", nil, &version)
	return version, err
}

// Heartbeat checks if the ChromaDB server is alive
func (c *Client) Heartbeat(ctx context.Context) (map[string]float64, error) {
	var result map[string]float64
	err := c.doRequest(ctx, http.MethodGet, "/api/v2/heartbeat", nil, &result)
	return result, err
}

// Reset resets the ChromaDB database (WARNING: This deletes all data)
func (c *Client) Reset(ctx context.Context) (bool, error) {
	var result bool
	err := c.doRequest(ctx, http.MethodPost, "/api/v2/reset", nil, &result)
	return result, err
}

// PreFlightChecks returns preflight check results
func (c *Client) PreFlightChecks(ctx context.Context) (PreflightChecks, error) {
	var result PreflightChecks
	err := c.doRequest(ctx, http.MethodGet, "/api/v2/pre-flight-checks", nil, &result)
	return result, err
}

// Root returns root endpoint information
func (c *Client) Root(ctx context.Context) (map[string]float64, error) {
	var result map[string]float64
	err := c.doRequest(ctx, http.MethodGet, "/api/v2", nil, &result)
	return result, err
}

// CreateTenant creates a new tenant
func (c *Client) CreateTenant(ctx context.Context, req CreateTenant) (*Tenant, error) {
	var result Tenant
	err := c.doRequest(ctx, http.MethodPost, "/api/v2/tenants", req, &result)
	return &result, err
}

// GetTenant gets a tenant by name
func (c *Client) GetTenant(ctx context.Context, name string) (*Tenant, error) {
	var result Tenant
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v2/tenants/%s", name), nil, &result)
	return &result, err
}

// CreateDatabase creates a new database
func (c *Client) CreateDatabase(ctx context.Context, req CreateDatabase, tenant ...string) (*Database, error) {
	tenantName := c.tenant
	if len(tenant) > 0 {
		tenantName = tenant[0]
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases", url.QueryEscape(tenantName))
	var result Database
	err := c.doRequest(ctx, http.MethodPost, path, req, &result)
	return &result, err
}

// GetDatabase gets a database by name
func (c *Client) GetDatabase(ctx context.Context, name string, tenant ...string) (*Database, error) {
	tenantName := c.tenant
	if len(tenant) > 0 {
		tenantName = tenant[0]
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases/%s", url.QueryEscape(tenantName), name)
	var result Database
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// ListCollections lists all collections
func (c *Client) ListCollections(ctx context.Context, tenant, database string) ([]Collection, error) {
	if tenant == "" {
		tenant = c.tenant
	}
	if database == "" {
		database = c.database
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections",
		url.QueryEscape(tenant), url.QueryEscape(database))

	var result []Collection
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

// CountCollections returns the number of collections
func (c *Client) CountCollections(ctx context.Context, tenant, database string) (int, error) {
	if tenant == "" {
		tenant = c.tenant
	}
	if database == "" {
		database = c.database
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections_count",
		url.QueryEscape(tenant), url.QueryEscape(database))

	var result int
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

// CreateCollection creates a new collection
func (c *Client) CreateCollection(ctx context.Context, req CreateCollection, tenant, database string) (*Collection, error) {
	if tenant == "" {
		tenant = c.tenant
	}
	if database == "" {
		database = c.database
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections",
		url.QueryEscape(tenant), url.QueryEscape(database))

	var result Collection
	err := c.doRequest(ctx, http.MethodPost, path, req, &result)
	return &result, err
}

// GetCollection gets a collection by name
func (c *Client) GetCollection(ctx context.Context, name string, tenant, database string) (*Collection, error) {
	if tenant == "" {
		tenant = c.tenant
	}
	if database == "" {
		database = c.database
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections/%s",
		url.QueryEscape(tenant), url.QueryEscape(database), name)

	var result Collection
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// DeleteCollection deletes a collection by name
func (c *Client) DeleteCollection(ctx context.Context, name string, tenant, database string) error {
	if tenant == "" {
		tenant = c.tenant
	}
	if database == "" {
		database = c.database
	}

	path := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections/%s",
		url.QueryEscape(tenant), url.QueryEscape(database), name)

	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// UpdateCollection updates a collection
func (c *Client) UpdateCollection(ctx context.Context, collectionID string, req UpdateCollection) (*Collection, error) {
	path := fmt.Sprintf("/api/v2/collections/%s", collectionID)
	var result Collection
	err := c.doRequest(ctx, http.MethodPut, path, req, &result)
	return &result, err
}

// Add adds embeddings to a collection
func (c *Client) Add(ctx context.Context, collectionID string, req AddEmbedding) error {
	path := fmt.Sprintf("/api/v2/collections/%s/add", collectionID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// Update updates embeddings in a collection
func (c *Client) Update(ctx context.Context, collectionID string, req UpdateEmbedding) error {
	path := fmt.Sprintf("/api/v2/collections/%s/update", collectionID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// Upsert upserts embeddings in a collection
func (c *Client) Upsert(ctx context.Context, collectionID string, req AddEmbedding) error {
	path := fmt.Sprintf("/api/v2/collections/%s/upsert", collectionID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// Get gets embeddings from a collection
func (c *Client) Get(ctx context.Context, collectionID string, req GetEmbedding) (*GetResult, error) {
	path := fmt.Sprintf("/api/v2/collections/%s/get", collectionID)
	var result GetResult
	err := c.doRequest(ctx, http.MethodPost, path, req, &result)
	return &result, err
}

// Delete deletes embeddings from a collection
func (c *Client) Delete(ctx context.Context, collectionID string, req DeleteEmbedding) ([]string, error) {
	path := fmt.Sprintf("/api/v2/collections/%s/delete", collectionID)
	var result []string
	err := c.doRequest(ctx, http.MethodPost, path, req, &result)
	return result, err
}

// Count returns the number of embeddings in a collection
func (c *Client) Count(ctx context.Context, collectionID string) (int, error) {
	path := fmt.Sprintf("/api/v2/collections/%s/count", collectionID)
	var result int
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

// Query queries a collection for nearest neighbors
func (c *Client) Query(ctx context.Context, collectionID string, req QueryEmbedding) (*QueryResult, error) {
	path := fmt.Sprintf("/api/v2/collections/%s/query", collectionID)
	var result QueryResult
	err := c.doRequest(ctx, http.MethodPost, path, req, &result)
	return &result, err
}
