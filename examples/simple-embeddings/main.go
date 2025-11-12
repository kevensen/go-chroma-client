package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"strings"

	chromaclient "github.com/kevensen/go-chroma-client"
)

// SimpleEmbeddingGenerator creates basic embeddings without external API calls.
// NOTE: This is a DEMO implementation only - not suitable for production use!
// For real semantic search, use proper embedding models like:
// - OpenAI's text-embedding-3-small
// - Sentence Transformers
// - Cohere embeddings
type SimpleEmbeddingGenerator struct {
	dimensions int
}

// NewSimpleEmbeddingGenerator creates a basic embedding generator
func NewSimpleEmbeddingGenerator(dimensions int) *SimpleEmbeddingGenerator {
	return &SimpleEmbeddingGenerator{
		dimensions: dimensions,
	}
}

// GenerateEmbedding creates a simple embedding from text.
// This uses a hash-based approach for demonstration purposes.
// Real embeddings should capture semantic meaning!
func (g *SimpleEmbeddingGenerator) GenerateEmbedding(text string) []float64 {
	// Normalize text
	text = strings.ToLower(strings.TrimSpace(text))

	// Create a deterministic hash
	hash := sha256.Sum256([]byte(text))

	embedding := make([]float64, g.dimensions)

	// Use hash bytes to create pseudo-random but deterministic values
	for i := 0; i < g.dimensions; i++ {
		// Get value from hash (reuse hash bytes in cycle)
		idx := i % len(hash)
		value := float64(hash[idx]) / 255.0 // Normalize to 0-1

		// Apply some variation based on position
		value = value*2.0 - 1.0 // Scale to -1 to 1

		// Add word-count based feature
		wordCount := float64(len(strings.Fields(text)))
		if i == 0 {
			embedding[i] = wordCount / 100.0
		} else {
			embedding[i] = value
		}
	}

	// Normalize the vector
	return normalizeVector(embedding)
}

// normalizeVector normalizes a vector to unit length
func normalizeVector(vec []float64) []float64 {
	var magnitude float64
	for _, v := range vec {
		magnitude += v * v
	}
	magnitude = math.Sqrt(magnitude)

	if magnitude > 0 {
		normalized := make([]float64, len(vec))
		for i, v := range vec {
			normalized[i] = v / magnitude
		}
		return normalized
	}

	return vec
}

// cosineSimilarity calculates cosine similarity between two vectors (for testing)
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, magA, magB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		magA += a[i] * a[i]
		magB += b[i] * b[i]
	}

	if magA == 0 || magB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(magA) * math.Sqrt(magB))
}

func main() {
	// Create ChromaDB client
	chromaClient := chromaclient.NewClient(
		chromaclient.WithBaseURL("http://localhost:8000"),
	)

	// Create a simple embedding generator (128 dimensions)
	embeddingGen := NewSimpleEmbeddingGenerator(128)

	ctx := context.Background()

	// Check server version
	fmt.Println("=== Connecting to ChromaDB ===")
	version, err := chromaClient.Version(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to ChromaDB: %v", err)
	}
	fmt.Printf("ChromaDB Version: %s\n", version)

	// Create a collection
	fmt.Println("\n=== Creating Collection ===")
	collection, err := chromaClient.CreateCollection(ctx, chromaclient.CreateCollection{
		Name:        "simple_embeddings_demo",
		GetOrCreate: true,
		Metadata: map[string]interface{}{
			"description": "Demo with simple hash-based embeddings",
			"note":        "Not for production - use real embedding models!",
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	fmt.Printf("Collection: %s (ID: %s)\n", collection.Name, collection.ID)

	// Sample documents
	documents := []string{
		"The cat sat on the mat",
		"Dogs are loyal companions",
		"Birds can fly in the sky",
		"Fish swim in the ocean",
		"Programming is fun and creative",
		"Go is a statically typed language",
		"Python is popular for data science",
		"Machine learning requires lots of data",
	}

	// Generate embeddings for all documents
	fmt.Println("\n=== Generating Simple Embeddings ===")
	fmt.Printf("Generating embeddings for %d documents...\n", len(documents))
	embeddings := make([][]float64, len(documents))
	for i, doc := range documents {
		embeddings[i] = embeddingGen.GenerateEmbedding(doc)
	}
	fmt.Printf("Generated %d embeddings (dimension: %d)\n", len(embeddings), len(embeddings[0]))

	// Demonstrate that similar texts get different embeddings
	// (because this is just a hash-based demo, not true semantic embeddings)
	fmt.Println("\n=== Note About This Demo ===")
	fmt.Println("⚠️  These are hash-based embeddings for demonstration only!")
	fmt.Println("⚠️  They do NOT capture semantic meaning like real embedding models.")
	fmt.Println("⚠️  For production, use OpenAI, Cohere, or Sentence Transformers.")

	test1 := embeddingGen.GenerateEmbedding("The cat sat on the mat")
	test2 := embeddingGen.GenerateEmbedding("The cat sat on the mat")
	test3 := embeddingGen.GenerateEmbedding("A feline rested on the rug") // Semantically similar

	fmt.Printf("\nSame text similarity: %.4f (should be 1.0)\n", cosineSimilarity(test1, test2))
	fmt.Printf("Similar meaning similarity: %.4f (would be high with real embeddings)\n", cosineSimilarity(test1, test3))

	// Create IDs and metadata
	ids := make([]string, len(documents))
	metadatas := make([]map[string]interface{}, len(documents))
	for i := range documents {
		ids[i] = fmt.Sprintf("doc_%d", i+1)
		category := "animals"
		if i >= 4 {
			category = "programming"
		}
		metadatas[i] = map[string]interface{}{
			"category": category,
			"index":    i,
		}
	}

	// Add documents to ChromaDB
	fmt.Println("\n=== Adding Documents to ChromaDB ===")
	err = chromaClient.Add(ctx, collection.ID, chromaclient.AddEmbedding{
		IDs:        ids,
		Documents:  documents,
		Embeddings: embeddings,
		Metadatas:  metadatas,
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}
	fmt.Printf("Successfully added %d documents\n", len(documents))

	// Count documents
	count, err := chromaClient.Count(ctx, collection.ID, "", "")
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	fmt.Printf("Total documents in collection: %d\n", count)

	// Perform a search
	fmt.Println("\n=== Performing Search ===")
	queryText := "The cat sat on the mat"
	fmt.Printf("Query: %q\n", queryText)

	queryEmbedding := embeddingGen.GenerateEmbedding(queryText)

	results, err := chromaClient.Query(ctx, collection.ID, chromaclient.QueryEmbedding{
		QueryEmbeddings: [][]float64{queryEmbedding},
		NResults:        3,
		Include: []chromaclient.Include{
			chromaclient.IncludeDocuments,
			chromaclient.IncludeDistances,
			chromaclient.IncludeMetadatas,
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}

	fmt.Println("\nTop 3 Results:")
	for i := 0; i < len(results.IDs[0]); i++ {
		fmt.Printf("\n%d. ID: %s\n", i+1, results.IDs[0][i])
		fmt.Printf("   Document: %s\n", results.Documents[0][i])
		fmt.Printf("   Distance: %.4f\n", results.Distances[0][i])
		if i < len(results.Metadatas[0]) {
			fmt.Printf("   Category: %v\n", results.Metadatas[0][i]["category"])
		}
	}

	// Filter by metadata
	fmt.Println("\n=== Searching with Metadata Filter ===")
	queryText2 := "Go is a statically typed language"
	queryEmbedding2 := embeddingGen.GenerateEmbedding(queryText2)

	results2, err := chromaClient.Query(ctx, collection.ID, chromaclient.QueryEmbedding{
		QueryEmbeddings: [][]float64{queryEmbedding2},
		NResults:        5,
		Where: map[string]interface{}{
			"category": "programming", // Only search programming documents
		},
		Include: []chromaclient.Include{
			chromaclient.IncludeDocuments,
			chromaclient.IncludeDistances,
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}

	fmt.Printf("Query: %q (category: programming)\n", queryText2)
	fmt.Println("\nResults:")
	for i := 0; i < len(results2.IDs[0]); i++ {
		fmt.Printf("\n%d. Document: %s\n", i+1, results2.Documents[0][i])
		fmt.Printf("   Distance: %.4f\n", results2.Distances[0][i])
	}

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("✅ You can integrate any embedding generation approach")
	fmt.Println("✅ ChromaDB stores and searches the embeddings you provide")
	fmt.Println("✅ For production: use proper embedding models (OpenAI, etc.)")
	fmt.Println("✅ This example shows the mechanics without requiring API keys")

	fmt.Println("\n=== Example completed successfully! ===")
}
