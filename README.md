# go-chroma-client

A comprehensive Go client library for ChromaDB 2.0 API with full endpoint coverage.

## Features

- ✅ Complete ChromaDB 2.0 API coverage
- ✅ Tenant operations (Create, Get)
- ✅ Database operations (Create, Get)
- ✅ Collection operations (Create, Get, List, Count, Delete, Update)
- ✅ Document operations (Add, Update, Upsert, Get, Delete, Count, Query)
- ✅ Utility operations (Version, Heartbeat, Reset, Pre-flight checks)
- ✅ Configurable client with custom HTTP client support
- ✅ Context support for all operations
- ✅ Comprehensive test coverage
- ✅ **No automatic embedding generation - you provide your own embeddings**

## Important: Embedding Generation

**This library does NOT generate embeddings for you.** Creating vectorized embeddings from your documents is the responsibility of the library user. You must provide pre-computed embeddings when adding or querying documents.

This design choice gives you full control over:
- Which embedding model to use (OpenAI, Sentence Transformers, Cohere, etc.)
- How embeddings are generated and cached
- Cost management for API-based embedding services
- Custom embedding strategies for your specific use case

## Installation

```bash
go get github.com/kevensen/go-chroma-client
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    chromaclient "github.com/kevensen/go-chroma-client"
)

func main() {
    // Create a new client
    client := chromaclient.NewClient(
        chromaclient.WithBaseURL("http://localhost:8000"),
    )

    ctx := context.Background()

    // Check server version
    version, err := client.Version(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ChromaDB version: %s\n", version)

    // Create a collection
    collection, err := client.CreateCollection(ctx, chromaclient.CreateCollection{
        Name: "my_collection",
        Metadata: map[string]interface{}{
            "description": "My first collection",
        },
    }, "", "")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created collection: %s (ID: %s)\n", collection.Name, collection.ID)

    // Add documents
    err = client.Add(ctx, collection.ID, chromaclient.AddEmbedding{
        IDs:       []string{"doc1", "doc2"},
        Documents: []string{"Hello world", "Goodbye world"},
        Metadatas: []map[string]interface{}{
            {"source": "greeting"},
            {"source": "farewell"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Count documents
    count, err := client.Count(ctx, collection.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Document count: %d\n", count)

    // Query documents
    result, err := client.Query(ctx, collection.ID, chromaclient.QueryEmbedding{
        QueryEmbeddings: [][]float64{{0.1, 0.2, 0.3}}, // Example embedding
        NResults:        2,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Query results: %+v\n", result)
}
```

For example, you might use:
- [OpenAI's embedding API](https://platform.openai.com/docs/guides/embeddings)
- [Sentence Transformers](https://www.sbert.net/) (via Python or other bindings)
- [Cohere's Embed API](https://docs.cohere.com/docs/embeddings)
- Any other embedding service or model of your choice

## API Reference

### Client Configuration

```go
// Create a client with default settings (localhost:8000)
client := chromaclient.NewClient()

// Create a client with custom options
client := chromaclient.NewClient(
    chromaclient.WithBaseURL("http://chroma.example.com:8000"),
    chromaclient.WithTenant("custom_tenant"),
    chromaclient.WithDatabase("custom_database"),
    chromaclient.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

### Utility Operations

```go
// Get server version
version, err := client.Version(ctx)

// Check server health
heartbeat, err := client.Heartbeat(ctx)

// Reset database (WARNING: Deletes all data)
success, err := client.Reset(ctx)

// Pre-flight checks
checks, err := client.PreFlightChecks(ctx)
```

### Tenant Operations

```go
// Create a tenant
tenant, err := client.CreateTenant(ctx, chromaclient.CreateTenant{
    Name: "my_tenant",
})

// Get a tenant
tenant, err := client.GetTenant(ctx, "my_tenant")
```

### Database Operations

```go
// Create a database
db, err := client.CreateDatabase(ctx, chromaclient.CreateDatabase{
    Name: "my_database",
}, "my_tenant")

// Get a database
db, err := client.GetDatabase(ctx, "my_database", "my_tenant")
```

### Collection Operations

```go
// List all collections
collections, err := client.ListCollections(ctx, "", "")

// Count collections
count, err := client.CountCollections(ctx, "", "")

// Create a collection
collection, err := client.CreateCollection(ctx, chromaclient.CreateCollection{
    Name: "my_collection",
    Metadata: map[string]interface{}{
        "key": "value",
    },
}, "", "")

// Get a collection
collection, err := client.GetCollection(ctx, "my_collection", "", "")

// Update a collection
newName := "updated_collection"
collection, err := client.UpdateCollection(ctx, collectionID, chromaclient.UpdateCollection{
    NewName: &newName,
    NewMetadata: map[string]interface{}{
        "updated": true,
    },
})

// Delete a collection
err := client.DeleteCollection(ctx, "my_collection", "", "")
```

### Document Operations

**Note:** All document operations that work with embeddings require you to provide pre-computed embedding vectors. The `Embeddings` field is optional in the API, but you must provide embeddings if you want to perform similarity searches using the `Query` operation.

```go
// Add documents with pre-computed embeddings
// You are responsible for generating these embeddings using your chosen model
err := client.Add(ctx, collectionID, chromaclient.AddEmbedding{
    IDs:        []string{"id1", "id2"},
    Documents:  []string{"doc1", "doc2"},
    Embeddings: [][]float64{{0.1, 0.2}, {0.3, 0.4}}, // Your pre-computed embeddings
    Metadatas: []map[string]interface{}{
        {"key": "value1"},
        {"key": "value2"},
    },
})

// Update documents
err := client.Update(ctx, collectionID, chromaclient.UpdateEmbedding{
    IDs:       []string{"id1"},
    Documents: []string{"updated doc"},
})

// Upsert documents (insert or update)
err := client.Upsert(ctx, collectionID, chromaclient.AddEmbedding{
    IDs:       []string{"id1", "id2"},
    Documents: []string{"doc1", "doc2"},
})

// Get documents
result, err := client.Get(ctx, collectionID, chromaclient.GetEmbedding{
    IDs:     []string{"id1", "id2"},
    Include: []string{chromaclient.IncludeDocuments, chromaclient.IncludeMetadatas},
})

// Delete documents
deletedIDs, err := client.Delete(ctx, collectionID, chromaclient.DeleteEmbedding{
    IDs: []string{"id1", "id2"},
})

// Count documents in a collection
count, err := client.Count(ctx, collectionID)

// Query for nearest neighbors using pre-computed query embeddings
// You must provide the query embedding vector(s)
result, err := client.Query(ctx, collectionID, chromaclient.QueryEmbedding{
    QueryEmbeddings: [][]float64{{0.1, 0.2, 0.3}}, // Your pre-computed query embedding
    NResults:        10,
    Where: map[string]interface{}{
        "key": "value",
    },
    Include: []string{
        chromaclient.IncludeDocuments,
        chromaclient.IncludeMetadatas,
        chromaclient.IncludeDistances,
    },
})
```

## Example: Using Your Own Embeddings

Here's a complete example showing how to use this library with your own embedding generation:

```go
package main

import (
    "context"
    "log"
    chromaclient "github.com/kevensen/go-chroma-client"
)

// generateEmbedding is a placeholder for your embedding generation logic
// In practice, you would call an embedding service or model here
func generateEmbedding(text string) []float64 {
    // TODO: Replace with actual embedding generation
    // Examples:
    // - Call OpenAI's embedding API
    // - Use a local embedding model
    // - Call Cohere, Anthropic, or other embedding services
    
    // This is just a dummy example - not real embeddings!
    return []float64{0.1, 0.2, 0.3, 0.4, 0.5}
}

func main() {
    client := chromaclient.NewClient(
        chromaclient.WithBaseURL("http://localhost:8000"),
    )
    ctx := context.Background()

    // Create a collection
    collection, err := client.CreateCollection(ctx, chromaclient.CreateCollection{
        Name: "my_documents",
    }, "", "")
    if err != nil {
        log.Fatal(err)
    }

    // Documents to add
    documents := []string{
        "The quick brown fox jumps over the lazy dog",
        "Machine learning is transforming technology",
    }

    // Generate embeddings for each document (YOUR RESPONSIBILITY)
    embeddings := make([][]float64, len(documents))
    for i, doc := range documents {
        embeddings[i] = generateEmbedding(doc)
    }

    // Add documents with their embeddings
    err = client.Add(ctx, collection.ID, chromaclient.AddEmbedding{
        IDs:        []string{"doc1", "doc2"},
        Documents:  documents,
        Embeddings: embeddings, // Your pre-computed embeddings
    })
    if err != nil {
        log.Fatal(err)
    }

    // Query with a search query
    queryText := "artificial intelligence"
    queryEmbedding := generateEmbedding(queryText) // Generate embedding for query

    results, err := client.Query(ctx, collection.ID, chromaclient.QueryEmbedding{
        QueryEmbeddings: [][]float64{queryEmbedding}, // Your pre-computed query embedding
        NResults:        5,
        Include: []string{
            chromaclient.IncludeDocuments,
            chromaclient.IncludeDistances,
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Query results: %+v\n", results)
}
```

## Error Handling

The client returns detailed HTTP errors:

```go
collection, err := client.GetCollection(ctx, "nonexistent", "", "")
if err != nil {
    if httpErr, ok := err.(*chromaclient.HTTPError); ok {
        fmt.Printf("HTTP Error %d: %s\n", httpErr.StatusCode, httpErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Examples

This repository includes several examples demonstrating different approaches to embedding generation:

### Basic Example
Location: `examples/basic/`

A basic example showing all ChromaDB operations without embeddings for queries. Good for understanding the client API.

### Simple Embeddings Example  
Location: `examples/simple-embeddings/`

A self-contained example using hash-based embeddings. **For demonstration only** - shows the integration pattern without requiring API keys. Not suitable for production (doesn't capture semantic meaning).

### Real Embeddings with OpenAI
Location: `examples/with-embeddings/`

A production-ready example using OpenAI's embedding API. Shows:
- Generating real semantic embeddings
- Performing similarity searches
- Best practices for embedding-based search

**Requires**: OpenAI API key (set `OPENAI_API_KEY` environment variable)

Run any example:
```bash
cd examples/with-embeddings
go run main.go
```

## Testing

Run the tests:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
