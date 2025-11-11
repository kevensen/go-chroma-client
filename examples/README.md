# Examples Overview

This directory contains three examples showing different approaches to using the go-chroma-client library with embeddings.

## Quick Comparison

| Example | Use Case | Embedding Quality | API Key Required | Production Ready |
|---------|----------|-------------------|------------------|------------------|
| **basic** | Learning the client API | N/A (no search) | No | Yes (for non-search use cases) |
| **simple-embeddings** | Understanding the pattern | ⚠️ Hash-based (not semantic) | No | ❌ Demo only |
| **with-embeddings** | Real semantic search | ✅ OpenAI embeddings | Yes (OpenAI) | ✅ Yes |

## Choosing an Example

### Start with `basic/`
- Learn ChromaDB operations (collections, add, update, delete)
- Understand the client API
- Test your ChromaDB setup
- **No embeddings involved** - just document storage

### Try `simple-embeddings/`  
- See how embeddings integrate with ChromaDB
- Understand the workflow without external dependencies
- Test locally without API keys
- ⚠️ **Not for real semantic search** - uses hash-based "fake" embeddings

### Use `with-embeddings/` for production
- Real semantic search capabilities
- Production-ready embedding generation
- Shows best practices
- Requires OpenAI API key (or adapt to other services)

## Example Workflows

### Learning Path
1. `basic/` - Understand the client
2. `simple-embeddings/` - See the embedding pattern
3. `with-embeddings/` - Implement real search

### Quick Start for Production
Go directly to `with-embeddings/` if you:
- Already understand ChromaDB concepts
- Need semantic search
- Have an embedding service ready (OpenAI, Cohere, etc.)

## Running the Examples

Each example has its own README with specific instructions. General pattern:

```bash
# Make sure ChromaDB is running
docker run -p 8000:8000 chromadb/chroma

# Run an example
cd examples/with-embeddings
go run main.go
```

## Adapting for Your Embedding Service

All examples follow the same pattern:

```go
// 1. Generate embeddings (YOUR CHOICE)
embeddings := generateEmbeddings(documents)

// 2. Add to ChromaDB (same for all)
client.Add(ctx, collectionID, chromaclient.AddEmbedding{
    IDs:        ids,
    Documents:  documents,
    Embeddings: embeddings,  // Your embeddings here
})

// 3. Query (YOUR CHOICE)
queryEmbedding := generateEmbedding(query)
results := client.Query(ctx, collectionID, chromaclient.QueryEmbedding{
    QueryEmbeddings: [][]float64{queryEmbedding},
})
```

### Popular Embedding Services

Replace the embedding generation in `with-embeddings/` with:

**Cohere**
```go
// Use Cohere Go SDK
// https://github.com/cohere-ai/cohere-go
```

**Hugging Face**
```go
// Call Hugging Face Inference API
// https://huggingface.co/docs/api-inference/
```

**Sentence Transformers (via Python)**
```go
// Call Python script with sentence-transformers
// exec.Command("python", "embed.py", text)
```

**Voyage AI**
```go
// Use Voyage AI's embedding API
// https://docs.voyageai.com/
```

## Key Principle

**This library (go-chroma-client) is embedding-agnostic.** It:
- ✅ Stores embeddings you provide
- ✅ Queries using embeddings you generate
- ❌ Does NOT generate embeddings for you

This design gives you complete control over:
- Which embedding model to use
- How to optimize costs
- When to cache embeddings
- How to batch process

## Questions?

- See the main [README.md](../README.md) for API documentation
- Each example has its own README with details
- Check [ChromaDB docs](https://docs.trychroma.com/) for more on vector databases
