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

```go
// Add documents
err := client.Add(ctx, collectionID, chromaclient.AddEmbedding{
    IDs:        []string{"id1", "id2"},
    Documents:  []string{"doc1", "doc2"},
    Embeddings: [][]float64{{0.1, 0.2}, {0.3, 0.4}},
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

// Query for nearest neighbors
result, err := client.Query(ctx, collectionID, chromaclient.QueryEmbedding{
    QueryEmbeddings: [][]float64{{0.1, 0.2, 0.3}},
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
