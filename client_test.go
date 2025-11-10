package chromaclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.baseURL != "http://localhost:8000" {
		t.Errorf("Expected baseURL to be http://localhost:8000, got %s", client.baseURL)
	}
	if client.tenant != DefaultTenant {
		t.Errorf("Expected tenant to be %s, got %s", DefaultTenant, client.tenant)
	}
	if client.database != DefaultDatabase {
		t.Errorf("Expected database to be %s, got %s", DefaultDatabase, client.database)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client := NewClient(
		WithBaseURL("http://example.com"),
		WithTenant("custom_tenant"),
		WithDatabase("custom_database"),
	)

	if client.baseURL != "http://example.com" {
		t.Errorf("Expected baseURL to be http://example.com, got %s", client.baseURL)
	}
	if client.tenant != "custom_tenant" {
		t.Errorf("Expected tenant to be custom_tenant, got %s", client.tenant)
	}
	if client.database != "custom_database" {
		t.Errorf("Expected database to be custom_database, got %s", client.database)
	}
}

func TestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/version" {
			t.Errorf("Expected path /api/v1/version, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("0.4.24")
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	version, err := client.Version(context.Background())
	if err != nil {
		t.Fatalf("Version() error = %v", err)
	}
	if version != "0.4.24" {
		t.Errorf("Expected version 0.4.24, got %s", version)
	}
}

func TestHeartbeat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/heartbeat" {
			t.Errorf("Expected path /api/v1/heartbeat, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"nanosecond heartbeat": 1234567890})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	result, err := client.Heartbeat(context.Background())
	if err != nil {
		t.Fatalf("Heartbeat() error = %v", err)
	}
	if result["nanosecond heartbeat"] != 1234567890 {
		t.Errorf("Expected heartbeat value 1234567890, got %f", result["nanosecond heartbeat"])
	}
}

func TestReset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/reset" {
			t.Errorf("Expected path /api/v1/reset, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	result, err := client.Reset(context.Background())
	if err != nil {
		t.Fatalf("Reset() error = %v", err)
	}
	if !result {
		t.Errorf("Expected reset to return true")
	}
}

func TestCreateTenant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tenants" {
			t.Errorf("Expected path /api/v1/tenants, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req CreateTenant
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Name != "test_tenant" {
			t.Errorf("Expected tenant name test_tenant, got %s", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Tenant{Name: req.Name})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	tenant, err := client.CreateTenant(context.Background(), CreateTenant{Name: "test_tenant"})
	if err != nil {
		t.Fatalf("CreateTenant() error = %v", err)
	}
	if tenant.Name != "test_tenant" {
		t.Errorf("Expected tenant name test_tenant, got %s", tenant.Name)
	}
}

func TestGetTenant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tenants/test_tenant" {
			t.Errorf("Expected path /api/v1/tenants/test_tenant, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Tenant{Name: "test_tenant"})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	tenant, err := client.GetTenant(context.Background(), "test_tenant")
	if err != nil {
		t.Fatalf("GetTenant() error = %v", err)
	}
	if tenant.Name != "test_tenant" {
		t.Errorf("Expected tenant name test_tenant, got %s", tenant.Name)
	}
}

func TestCreateDatabase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/databases" {
			t.Errorf("Expected path /api/v1/databases, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		tenant := r.URL.Query().Get("tenant")
		if tenant != "default_tenant" {
			t.Errorf("Expected tenant default_tenant, got %s", tenant)
		}

		var req CreateDatabase
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Name != "test_database" {
			t.Errorf("Expected database name test_database, got %s", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Database{
			ID:     "db-123",
			Name:   req.Name,
			Tenant: tenant,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	database, err := client.CreateDatabase(context.Background(), CreateDatabase{Name: "test_database"})
	if err != nil {
		t.Fatalf("CreateDatabase() error = %v", err)
	}
	if database.Name != "test_database" {
		t.Errorf("Expected database name test_database, got %s", database.Name)
	}
}

func TestGetDatabase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/databases/test_database" {
			t.Errorf("Expected path /api/v1/databases/test_database, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Database{
			ID:     "db-123",
			Name:   "test_database",
			Tenant: "default_tenant",
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	database, err := client.GetDatabase(context.Background(), "test_database")
	if err != nil {
		t.Fatalf("GetDatabase() error = %v", err)
	}
	if database.Name != "test_database" {
		t.Errorf("Expected database name test_database, got %s", database.Name)
	}
}

func TestListCollections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections" {
			t.Errorf("Expected path /api/v1/collections, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Collection{
			{ID: "col-1", Name: "collection1"},
			{ID: "col-2", Name: "collection2"},
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	collections, err := client.ListCollections(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListCollections() error = %v", err)
	}
	if len(collections) != 2 {
		t.Errorf("Expected 2 collections, got %d", len(collections))
	}
}

func TestCountCollections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/count_collections" {
			t.Errorf("Expected path /api/v1/count_collections, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(5)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	count, err := client.CountCollections(context.Background(), "", "")
	if err != nil {
		t.Fatalf("CountCollections() error = %v", err)
	}
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestCreateCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections" {
			t.Errorf("Expected path /api/v1/collections, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req CreateCollection
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Name != "test_collection" {
			t.Errorf("Expected collection name test_collection, got %s", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Collection{
			ID:       "col-123",
			Name:     req.Name,
			Metadata: req.Metadata,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	collection, err := client.CreateCollection(context.Background(), CreateCollection{
		Name:     "test_collection",
		Metadata: map[string]interface{}{"key": "value"},
	}, "", "")
	if err != nil {
		t.Fatalf("CreateCollection() error = %v", err)
	}
	if collection.Name != "test_collection" {
		t.Errorf("Expected collection name test_collection, got %s", collection.Name)
	}
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/test_collection" {
			t.Errorf("Expected path /api/v1/collections/test_collection, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Collection{
			ID:   "col-123",
			Name: "test_collection",
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	collection, err := client.GetCollection(context.Background(), "test_collection", "", "")
	if err != nil {
		t.Fatalf("GetCollection() error = %v", err)
	}
	if collection.Name != "test_collection" {
		t.Errorf("Expected collection name test_collection, got %s", collection.Name)
	}
}

func TestDeleteCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/test_collection" {
			t.Errorf("Expected path /api/v1/collections/test_collection, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Collection{
			ID:   "col-123",
			Name: "test_collection",
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	err := client.DeleteCollection(context.Background(), "test_collection", "", "")
	if err != nil {
		t.Fatalf("DeleteCollection() error = %v", err)
	}
}

func TestUpdateCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/col-123" {
			t.Errorf("Expected path /api/v1/collections/col-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		var req UpdateCollection
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.NewName == nil || *req.NewName != "updated_collection" {
			t.Errorf("Expected new name updated_collection")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Collection{
			ID:   "col-123",
			Name: *req.NewName,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	newName := "updated_collection"
	collection, err := client.UpdateCollection(context.Background(), "col-123", UpdateCollection{
		NewName: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateCollection() error = %v", err)
	}
	if collection.Name != "updated_collection" {
		t.Errorf("Expected collection name updated_collection, got %s", collection.Name)
	}
}

func TestAdd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/col-123/add" {
			t.Errorf("Expected path /api/v1/collections/col-123/add, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req AddEmbedding
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if len(req.IDs) != 2 {
			t.Errorf("Expected 2 IDs, got %d", len(req.IDs))
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	err := client.Add(context.Background(), "col-123", AddEmbedding{
		IDs:       []string{"id1", "id2"},
		Documents: []string{"doc1", "doc2"},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
}

func TestCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/col-123/count" {
			t.Errorf("Expected path /api/v1/collections/col-123/count, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(42)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	count, err := client.Count(context.Background(), "col-123")
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 42 {
		t.Errorf("Expected count 42, got %d", count)
	}
}

func TestQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/collections/col-123/query" {
			t.Errorf("Expected path /api/v1/collections/col-123/query, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req QueryEmbedding
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(QueryResult{
			IDs:       [][]string{{"id1", "id2"}},
			Documents: [][]string{{"doc1", "doc2"}},
			Distances: [][]float64{{0.1, 0.2}},
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	result, err := client.Query(context.Background(), "col-123", QueryEmbedding{
		QueryEmbeddings: [][]float64{{0.1, 0.2, 0.3}},
		NResults:        10,
	})
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(result.IDs) != 1 {
		t.Errorf("Expected 1 result group, got %d", len(result.IDs))
	}
}

func TestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Collection not found"))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.GetCollection(context.Background(), "nonexistent", "", "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("Expected HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, httpErr.StatusCode)
	}
}
