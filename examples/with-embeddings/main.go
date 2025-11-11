package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	chromaclient "github.com/kevensen/go-chroma-client"
)

// OpenAI API structures
type OpenAIEmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format,omitempty"`
}

type OpenAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// EmbeddingGenerator handles embedding creation
type EmbeddingGenerator struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewEmbeddingGenerator creates a new embedding generator
func NewEmbeddingGenerator(apiKey string) *EmbeddingGenerator {
	return &EmbeddingGenerator{
		apiKey:     apiKey,
		model:      "text-embedding-3-small", // OpenAI's latest small embedding model
		httpClient: &http.Client{},
	}
}

// GenerateEmbeddings calls OpenAI's API to generate embeddings for the given texts
func (g *EmbeddingGenerator) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
	reqBody := OpenAIEmbeddingRequest{
		Input: texts,
		Model: g.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/embeddings", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var embeddingResp OpenAIEmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract embeddings in the correct order
	embeddings := make([][]float64, len(texts))
	for _, data := range embeddingResp.Data {
		if data.Index < len(embeddings) {
			embeddings[data.Index] = data.Embedding
		}
	}

	return embeddings, nil
}

func main() {
	// Get OpenAI API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable must be set")
	}

	// Create ChromaDB client
	chromaClient := chromaclient.NewClient(
		chromaclient.WithBaseURL("http://localhost:8000"),
	)

	// Create embedding generator
	embeddingGen := NewEmbeddingGenerator(apiKey)

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
		Name:        "articles_with_embeddings",
		GetOrCreate: true,
		Metadata: map[string]interface{}{
			"description": "Article collection with OpenAI embeddings",
		},
	}, "", "")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	fmt.Printf("Collection: %s (ID: %s)\n", collection.Name, collection.ID)

	// Sample documents to add
	documents := []string{
		"Artificial intelligence is transforming the way we interact with technology.",
		"Machine learning algorithms can identify patterns in large datasets.",
		"Natural language processing enables computers to understand human language.",
		"Deep learning neural networks are inspired by the human brain.",
		"Computer vision allows machines to interpret and understand visual information.",
	}

	// Generate embeddings for documents
	fmt.Println("\n=== Generating Embeddings ===")
	fmt.Printf("Generating embeddings for %d documents using OpenAI...\n", len(documents))
	embeddings, err := embeddingGen.GenerateEmbeddings(ctx, documents)
	if err != nil {
		log.Fatalf("Failed to generate embeddings: %v", err)
	}
	fmt.Printf("Generated %d embeddings (dimension: %d)\n", len(embeddings), len(embeddings[0]))

	// Create IDs for documents
	ids := make([]string, len(documents))
	metadatas := make([]map[string]interface{}, len(documents))
	for i := range documents {
		ids[i] = fmt.Sprintf("doc_%d", i+1)
		metadatas[i] = map[string]interface{}{
			"topic":  "AI/ML",
			"index":  i,
			"source": "example",
		}
	}

	// Add documents with embeddings to ChromaDB
	fmt.Println("\n=== Adding Documents to ChromaDB ===")
	err = chromaClient.Add(ctx, collection.ID, chromaclient.AddEmbedding{
		IDs:        ids,
		Documents:  documents,
		Embeddings: embeddings, // Our generated embeddings
		Metadatas:  metadatas,
	})
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}
	fmt.Printf("Successfully added %d documents with embeddings\n", len(documents))

	// Count documents
	count, err := chromaClient.Count(ctx, collection.ID)
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}
	fmt.Printf("Total documents in collection: %d\n", count)

	// Perform a semantic search
	fmt.Println("\n=== Performing Semantic Search ===")
	queryText := "understanding images and visual data"
	fmt.Printf("Query: %q\n", queryText)

	// Generate embedding for the query
	queryEmbeddings, err := embeddingGen.GenerateEmbeddings(ctx, []string{queryText})
	if err != nil {
		log.Fatalf("Failed to generate query embedding: %v", err)
	}

	// Query ChromaDB
	results, err := chromaClient.Query(ctx, collection.ID, chromaclient.QueryEmbedding{
		QueryEmbeddings: queryEmbeddings,
		NResults:        3,
		Include: []string{
			chromaclient.IncludeDocuments,
			chromaclient.IncludeDistances,
			chromaclient.IncludeMetadatas,
		},
	})
	if err != nil {
		log.Fatalf("Failed to query documents: %v", err)
	}

	// Display results
	fmt.Println("\nTop 3 Results:")
	for i := 0; i < len(results.IDs[0]); i++ {
		fmt.Printf("\n%d. ID: %s\n", i+1, results.IDs[0][i])
		fmt.Printf("   Document: %s\n", results.Documents[0][i])
		fmt.Printf("   Distance: %.4f\n", results.Distances[0][i])
		if i < len(results.Metadatas[0]) {
			fmt.Printf("   Metadata: %v\n", results.Metadatas[0][i])
		}
	}

	// Try another query
	fmt.Println("\n=== Second Semantic Search ===")
	queryText2 := "neural networks and brain-inspired computing"
	fmt.Printf("Query: %q\n", queryText2)

	queryEmbeddings2, err := embeddingGen.GenerateEmbeddings(ctx, []string{queryText2})
	if err != nil {
		log.Fatalf("Failed to generate query embedding: %v", err)
	}

	results2, err := chromaClient.Query(ctx, collection.ID, chromaclient.QueryEmbedding{
		QueryEmbeddings: queryEmbeddings2,
		NResults:        3,
		Include: []string{
			chromaclient.IncludeDocuments,
			chromaclient.IncludeDistances,
		},
	})
	if err != nil {
		log.Fatalf("Failed to query documents: %v", err)
	}

	fmt.Println("\nTop 3 Results:")
	for i := 0; i < len(results2.IDs[0]); i++ {
		fmt.Printf("\n%d. ID: %s\n", i+1, results2.IDs[0][i])
		fmt.Printf("   Document: %s\n", results2.Documents[0][i])
		fmt.Printf("   Distance: %.4f\n", results2.Distances[0][i])
	}

	// Cleanup (optional)
	// fmt.Println("\n=== Cleanup ===")
	// err = chromaClient.DeleteCollection(ctx, collection.Name, "", "")
	// if err != nil {
	// 	log.Fatalf("Failed to delete collection: %v", err)
	// }
	// fmt.Printf("Deleted collection: %s\n", collection.Name)

	fmt.Println("\n=== Example completed successfully! ===")
}
