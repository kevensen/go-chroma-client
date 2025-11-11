package chromaclient

import "time"

// Tenant represents a ChromaDB tenant
type Tenant struct {
	Name string `json:"name"`
}

// CreateTenant is the request body for creating a tenant
type CreateTenant struct {
	Name string `json:"name"`
}

// Database represents a ChromaDB database
type Database struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Tenant string `json:"tenant"`
}

// CreateDatabase is the request body for creating a database
type CreateDatabase struct {
	Name string `json:"name"`
}

// Collection represents a ChromaDB collection
type Collection struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreateCollection is the request body for creating a collection
type CreateCollection struct {
	Name        string                 `json:"name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	GetOrCreate bool                   `json:"get_or_create,omitempty"`
}

// UpdateCollection is the request body for updating a collection
type UpdateCollection struct {
	NewName     *string                `json:"new_name,omitempty"`
	NewMetadata map[string]interface{} `json:"new_metadata,omitempty"`
}

// AddEmbedding is the request body for adding embeddings
type AddEmbedding struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float64              `json:"embeddings,omitempty"`
	Metadatas  []map[string]interface{} `json:"metadatas,omitempty"`
	Documents  []string                 `json:"documents,omitempty"`
}

// UpdateEmbedding is the request body for updating embeddings
type UpdateEmbedding struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float64              `json:"embeddings,omitempty"`
	Metadatas  []map[string]interface{} `json:"metadatas,omitempty"`
	Documents  []string                 `json:"documents,omitempty"`
}

// GetEmbedding is the request body for getting embeddings
type GetEmbedding struct {
	IDs           []string               `json:"ids,omitempty"`
	Where         map[string]interface{} `json:"where,omitempty"`
	WhereDocument map[string]interface{} `json:"where_document,omitempty"`
	Sort          string                 `json:"sort,omitempty"`
	Limit         *int                   `json:"limit,omitempty"`
	Offset        *int                   `json:"offset,omitempty"`
	Include       []string               `json:"include,omitempty"`
}

// DeleteEmbedding is the request body for deleting embeddings
type DeleteEmbedding struct {
	IDs           []string               `json:"ids,omitempty"`
	Where         map[string]interface{} `json:"where,omitempty"`
	WhereDocument map[string]interface{} `json:"where_document,omitempty"`
}

// QueryEmbedding is the request body for querying embeddings
type QueryEmbedding struct {
	QueryEmbeddings [][]float64            `json:"query_embeddings"`
	NResults        int                    `json:"n_results,omitempty"`
	Where           map[string]interface{} `json:"where,omitempty"`
	WhereDocument   map[string]interface{} `json:"where_document,omitempty"`
	Include         []string               `json:"include,omitempty"`
}

// GetResult represents the result of a get operation
type GetResult struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float64              `json:"embeddings,omitempty"`
	Documents  []string                 `json:"documents,omitempty"`
	Metadatas  []map[string]interface{} `json:"metadatas,omitempty"`
}

// QueryResult represents the result of a query operation
type QueryResult struct {
	IDs        [][]string                 `json:"ids"`
	Embeddings [][][]float64              `json:"embeddings,omitempty"`
	Documents  [][]string                 `json:"documents,omitempty"`
	Metadatas  [][]map[string]interface{} `json:"metadatas,omitempty"`
	Distances  [][]float64                `json:"distances,omitempty"`
}

// PreflightChecks represents the preflight check response
type PreflightChecks map[string]interface{}

// IncludeOption represents what to include in results
const (
	IncludeDocuments  = "documents"
	IncludeEmbeddings = "embeddings"
	IncludeMetadatas  = "metadatas"
	IncludeDistances  = "distances"
)

// DefaultTenant is the default tenant name
const DefaultTenant = "default_tenant"

// DefaultDatabase is the default database name
const DefaultDatabase = "default_database"

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int
	Message    string
	Timestamp  time.Time
}

func (e *HTTPError) Error() string {
	return e.Message
}
