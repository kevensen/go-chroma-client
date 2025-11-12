package chromaclient

import "time"

// Space represents the vector space for similarity calculation
type Space string

const (
	SpaceL2     Space = "l2"
	SpaceCosine Space = "cosine"
	SpaceIP     Space = "ip"
)

// Include represents what to include in query/get results
type Include string

const (
	IncludeDistances  Include = "distances"
	IncludeDocuments  Include = "documents"
	IncludeEmbeddings Include = "embeddings"
	IncludeMetadatas  Include = "metadatas"
	IncludeUris       Include = "uris"
)

// HnswConfiguration represents HNSW index configuration
type HnswConfiguration struct {
	EfConstruction *int     `json:"ef_construction,omitempty"`
	EfSearch       *int     `json:"ef_search,omitempty"`
	MaxNeighbors   *int     `json:"max_neighbors,omitempty"`
	ResizeFactor   *float64 `json:"resize_factor,omitempty"`
	Space          *Space   `json:"space,omitempty"`
	SyncThreshold  *int     `json:"sync_threshold,omitempty"`
}

// SpannConfiguration represents SPANN index configuration
type SpannConfiguration struct {
	EfConstruction        *int   `json:"ef_construction,omitempty"`
	EfSearch              *int   `json:"ef_search,omitempty"`
	MaxNeighbors          *int   `json:"max_neighbors,omitempty"`
	MergeThreshold        *int32 `json:"merge_threshold,omitempty"`
	ReassignNeighborCount *int32 `json:"reassign_neighbor_count,omitempty"`
	SearchNprobe          *int32 `json:"search_nprobe,omitempty"`
	Space                 *Space `json:"space,omitempty"`
	SplitThreshold        *int32 `json:"split_threshold,omitempty"`
	WriteNprobe           *int32 `json:"write_nprobe,omitempty"`
}

// EmbeddingFunctionConfiguration represents embedding function configuration
type EmbeddingFunctionConfiguration struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CollectionConfiguration represents collection-level configuration
type CollectionConfiguration struct {
	EmbeddingFunction *EmbeddingFunctionConfiguration `json:"embedding_function,omitempty"`
	Hnsw              *HnswConfiguration              `json:"hnsw,omitempty"`
	Spann             *SpannConfiguration             `json:"spann,omitempty"`
}

// ValueTypes represents index configurations for each value type
type ValueTypes struct {
	Bool         map[string]interface{} `json:"bool,omitempty"`
	Float        map[string]interface{} `json:"float,omitempty"`
	FloatList    map[string]interface{} `json:"float_list,omitempty"`
	Int          map[string]interface{} `json:"int,omitempty"`
	SparseVector map[string]interface{} `json:"sparse_vector,omitempty"`
	String       map[string]interface{} `json:"string,omitempty"`
}

// InternalSchema represents the server-side schema structure
type InternalSchema struct {
	Defaults ValueTypes            `json:"defaults"`
	Keys     map[string]ValueTypes `json:"keys,omitempty"`
}

// Tenant represents a ChromaDB tenant
type Tenant struct {
	Name string `json:"name"`
}

// GetTenantResponse represents the response from getting a tenant
type GetTenantResponse struct {
	Name         string  `json:"name"`
	ResourceName *string `json:"resource_name,omitempty"`
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
	ID                string                  `json:"id"`
	Name              string                  `json:"name"`
	ConfigurationJSON CollectionConfiguration `json:"configuration_json"`
	Tenant            string                  `json:"tenant"`
	Database          string                  `json:"database"`
	LogPosition       int64                   `json:"log_position"`
	Version           int32                   `json:"version"`
	Dimension         *int32                  `json:"dimension,omitempty"`
	Metadata          map[string]interface{}  `json:"metadata,omitempty"`
	Schema            *InternalSchema         `json:"schema,omitempty"`
}

// CreateCollection is the request body for creating a collection
type CreateCollection struct {
	Name          string                   `json:"name"`
	Metadata      map[string]interface{}   `json:"metadata,omitempty"`
	GetOrCreate   bool                     `json:"get_or_create,omitempty"`
	Configuration *CollectionConfiguration `json:"configuration,omitempty"`
	Schema        *InternalSchema          `json:"schema,omitempty"`
}

// UpdateCollection is the request body for updating a collection
type UpdateCollection struct {
	NewName          *string                  `json:"new_name,omitempty"`
	NewMetadata      map[string]interface{}   `json:"new_metadata,omitempty"`
	NewConfiguration *CollectionConfiguration `json:"new_configuration,omitempty"`
}

// AddEmbedding is the request body for adding embeddings
type AddEmbedding struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float64              `json:"embeddings,omitempty"`
	Metadatas  []map[string]interface{} `json:"metadatas,omitempty"`
	Documents  []string                 `json:"documents,omitempty"`
	Uris       []string                 `json:"uris,omitempty"`
}

// UpdateEmbedding is the request body for updating embeddings
type UpdateEmbedding struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float64              `json:"embeddings,omitempty"`
	Metadatas  []map[string]interface{} `json:"metadatas,omitempty"`
	Documents  []string                 `json:"documents,omitempty"`
	Uris       []string                 `json:"uris,omitempty"`
}

// GetEmbedding is the request body for getting embeddings
type GetEmbedding struct {
	IDs           []string               `json:"ids,omitempty"`
	Where         map[string]interface{} `json:"where,omitempty"`
	WhereDocument map[string]interface{} `json:"where_document,omitempty"`
	Sort          string                 `json:"sort,omitempty"`
	Limit         *int                   `json:"limit,omitempty"`
	Offset        *int                   `json:"offset,omitempty"`
	Include       []Include              `json:"include,omitempty"`
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
	Include         []Include              `json:"include,omitempty"`
	IDs             []string               `json:"ids,omitempty"`
}

// GetResult represents the result of a get operation
type GetResult struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float64              `json:"embeddings,omitempty"`
	Documents  []string                 `json:"documents,omitempty"`
	Metadatas  []map[string]interface{} `json:"metadatas,omitempty"`
	Uris       []string                 `json:"uris,omitempty"`
	Include    []Include                `json:"include"`
}

// QueryResult represents the result of a query operation
type QueryResult struct {
	IDs        [][]string                 `json:"ids"`
	Embeddings [][][]float64              `json:"embeddings,omitempty"`
	Documents  [][]string                 `json:"documents,omitempty"`
	Metadatas  [][]map[string]interface{} `json:"metadatas,omitempty"`
	Distances  [][]float64                `json:"distances,omitempty"`
	Uris       [][]string                 `json:"uris,omitempty"`
	Include    []Include                  `json:"include"`
}

// PreflightChecks represents the preflight check response
type PreflightChecks map[string]interface{}

// HeartbeatResponse represents the heartbeat response
type HeartbeatResponse struct {
	NanosecondHeartbeat int64 `json:"nanosecond heartbeat"`
}

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
