# ChromaDB with Embeddings Example

This example demonstrates how to use the go-chroma-client library with a real embedding service. It shows how to:

1. Generate embeddings using OpenAI's API
2. Store documents with embeddings in ChromaDB
3. Perform semantic search queries

## Prerequisites

1. **ChromaDB Server**: You need a running ChromaDB server
   ```bash
   docker run -p 8000:8000 chromadb/chroma
   ```

2. **OpenAI API Key**: Sign up at [OpenAI](https://platform.openai.com/) and get an API key

## Setup

1. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Example

```bash
cd examples/with-embeddings
go run main.go
```

## What This Example Does

1. **Connects to ChromaDB**: Verifies connection to a local ChromaDB instance
2. **Creates a Collection**: Sets up a new collection for storing documents
3. **Generates Embeddings**: Uses OpenAI's `text-embedding-3-small` model to create embeddings for sample documents
4. **Stores Documents**: Adds documents with their embeddings to ChromaDB
5. **Semantic Search**: Performs similarity searches using query embeddings
6. **Displays Results**: Shows the most similar documents with their distances

## Sample Output

```
=== Connecting to ChromaDB ===
ChromaDB Version: 0.5.0

=== Creating Collection ===
Collection: articles_with_embeddings (ID: 12345-abcde)

=== Generating Embeddings ===
Generating embeddings for 5 documents using OpenAI...
Generated 5 embeddings (dimension: 1536)

=== Adding Documents to ChromaDB ===
Successfully added 5 documents with embeddings
Total documents in collection: 5

=== Performing Semantic Search ===
Query: "understanding images and visual data"

Top 3 Results:

1. ID: doc_5
   Document: Computer vision allows machines to interpret and understand visual information.
   Distance: 0.2341
   Metadata: map[index:4 source:example topic:AI/ML]

2. ID: doc_3
   Document: Natural language processing enables computers to understand human language.
   Distance: 0.3567
   Metadata: map[index:2 source:example topic:AI/ML]
...
```

## Using Different Embedding Services

This example uses OpenAI, but you can easily adapt it to use other services:

### Cohere
```go
// Use Cohere's embedding API
// https://docs.cohere.com/docs/embeddings
```

### Hugging Face
```go
// Use Hugging Face's inference API
// https://huggingface.co/docs/api-inference/
```

### Local Models
```go
// Use sentence-transformers via Python bindings
// Or use a Go-native embedding solution
```

## Cost Considerations

- OpenAI's `text-embedding-3-small` model costs approximately $0.02 per 1M tokens
- This example generates embeddings for 5 short documents (minimal cost)
- For production use, consider:
  - Caching embeddings to avoid regenerating them
  - Batch processing to reduce API calls
  - Using smaller or local models for cost savings

## Architecture Notes

**Important**: This library (go-chroma-client) does NOT include embedding generation. You are responsible for:

- Choosing an embedding model
- Generating embeddings for your documents
- Managing API keys and costs
- Caching and optimization strategies

This separation of concerns gives you full control over your embedding pipeline while keeping the ChromaDB client focused on database operations.

## Learn More

- [OpenAI Embeddings Guide](https://platform.openai.com/docs/guides/embeddings)
- [ChromaDB Documentation](https://docs.trychroma.com/)
- [Choosing an Embedding Model](https://www.sbert.net/docs/pretrained_models.html)
