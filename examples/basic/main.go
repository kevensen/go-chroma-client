package main

import (
	"context"
	"fmt"
	"log"

	chromaclient "github.com/kevensen/go-chroma-client"
)

func main() {
	// Create a new ChromaDB client
	client := chromaclient.NewClient(
		chromaclient.WithBaseURL("http://localhost:8000"),
	)

	ctx := context.Background()

	// 1. Check server version and health
	fmt.Println("=== Server Information ===")
	version, err := client.Version(ctx)
	if err != nil {
		log.Printf("Warning: Could not get version: %v\n", err)
	} else {
		fmt.Printf("ChromaDB Version: %s\n", version)
	}

	heartbeat, err := client.Heartbeat(ctx)
	if err != nil {
		log.Printf("Warning: Could not get heartbeat: %v\n", err)
	} else {
		fmt.Printf("Server Heartbeat: %v\n", heartbeat)
	}

	// 2. Tenant operations (optional, using default tenant)
	fmt.Println("\n=== Tenant Operations ===")
	tenant, err := client.GetTenant(ctx, chromaclient.DefaultTenant)
	if err != nil {
		log.Printf("Note: Default tenant check failed (this is normal if multi-tenancy is not enabled): %v\n", err)
	} else {
		fmt.Printf("Default Tenant: %s\n", tenant.Name)
	}

	// 3. Collection operations
	fmt.Println("\n=== Collection Operations ===")
	collectionName := "example_collection"

	// Create or get a collection
	collection, err := client.CreateCollection(ctx, chromaclient.CreateCollection{
		Name:        collectionName,
		GetOrCreate: true,
		Metadata: map[string]interface{}{
			"description": "Example collection for demo",
			"created_by":  "go-chroma-client",
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	fmt.Printf("Collection: %s (ID: %s)\n", collection.Name, collection.ID)

	// List all collections
	collections, err := client.ListCollections(ctx, "", "")
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}
	fmt.Printf("Total collections: %d\n", len(collections))
	for _, col := range collections {
		fmt.Printf("  - %s (ID: %s)\n", col.Name, col.ID)
	}

	// 4. Document operations
	fmt.Println("\n=== Document Operations ===")

	// Add documents
	err = client.Add(ctx, collection.ID, chromaclient.AddEmbedding{
		IDs: []string{"doc1", "doc2", "doc3"},
		Documents: []string{
			"The quick brown fox jumps over the lazy dog",
			"Machine learning is a subset of artificial intelligence",
			"ChromaDB is a vector database for embeddings",
		},
		Metadatas: []map[string]interface{}{
			{"category": "proverb", "length": "short"},
			{"category": "tech", "topic": "AI"},
			{"category": "tech", "topic": "database"},
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}
	fmt.Println("Added 3 documents")

	// Count documents
	count, err := client.Count(ctx, collection.ID, "", "")
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	fmt.Printf("Document count: %d\n", count)

	// Get documents
	fmt.Println("\n=== Retrieving Documents ===")
	getResult, err := client.Get(ctx, collection.ID, chromaclient.GetEmbedding{
		IDs: []string{"doc1", "doc2"},
		Include: []chromaclient.Include{
			chromaclient.IncludeDocuments,
			chromaclient.IncludeMetadatas,
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to get documents: %v", err)
	}
	for i, id := range getResult.IDs {
		fmt.Printf("ID: %s\n", id)
		if i < len(getResult.Documents) {
			fmt.Printf("  Document: %s\n", getResult.Documents[i])
		}
		if i < len(getResult.Metadatas) {
			fmt.Printf("  Metadata: %v\n", getResult.Metadatas[i])
		}
	}

	// Update a document
	fmt.Println("\n=== Updating Document ===")
	err = client.Update(ctx, collection.ID, chromaclient.UpdateEmbedding{
		IDs:       []string{"doc1"},
		Documents: []string{"The quick brown fox jumps over the lazy dog - Updated!"},
		Metadatas: []map[string]interface{}{
			{"category": "proverb", "length": "short", "updated": true},
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}
	fmt.Println("Updated doc1")

	// Upsert (insert or update)
	fmt.Println("\n=== Upserting Document ===")
	err = client.Upsert(ctx, collection.ID, chromaclient.AddEmbedding{
		IDs:       []string{"doc4"},
		Documents: []string{"This is a new document added via upsert"},
		Metadatas: []map[string]interface{}{
			{"category": "demo", "method": "upsert"},
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to upsert document: %v", err)
	}
	fmt.Println("Upserted doc4")

	// Count again
	count, err = client.Count(ctx, collection.ID, "", "")
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	fmt.Printf("Document count after upsert: %d\n", count)

	// Delete a document
	fmt.Println("\n=== Deleting Document ===")
	err = client.Delete(ctx, collection.ID, chromaclient.DeleteEmbedding{
		IDs: []string{"doc4"},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to delete document: %v", err)
	}
	fmt.Println("Deleted doc4")

	// Final count
	count, err = client.Count(ctx, collection.ID, "", "")
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	fmt.Printf("Final document count: %d\n", count)

	// 5. Collection update
	fmt.Println("\n=== Updating Collection ===")
	err = client.UpdateCollection(ctx, collection.ID, chromaclient.UpdateCollection{
		NewMetadata: map[string]interface{}{
			"description": "Example collection for demo - Updated",
			"created_by":  "go-chroma-client",
			"updated":     true,
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to update collection: %v", err)
	}
	fmt.Println("Updated collection metadata")

	// 6. Cleanup (optional - uncomment to delete the collection)
	// fmt.Println("\n=== Cleanup ===")
	// err = client.DeleteCollection(ctx, collectionName, "", "")
	// if err != nil {
	// 	log.Fatalf("Failed to delete collection: %v", err)
	// }
	// fmt.Printf("Deleted collection: %s\n", collectionName)

	fmt.Println("\n=== Example completed successfully! ===")
}
