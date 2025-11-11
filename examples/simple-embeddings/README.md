# Simple Embeddings Example

This example demonstrates how to use the go-chroma-client library with a simple, self-contained embedding generator. **This is for demonstration purposes only** - it does not use real semantic embeddings.

## Purpose

This example shows:
1. How to structure an embedding generator
2. How embeddings integrate with ChromaDB
3. The mechanics of the workflow without requiring external API keys
4. That you can use ANY embedding approach with this library

## ⚠️ Important Note

**The embeddings in this example are hash-based and do NOT capture semantic meaning.** They are deterministic and consistent, but won't find semantically similar documents like real embedding models would.

For production use, you should use proper embedding models:
- OpenAI's text-embedding models (see `examples/with-embeddings`)
- Cohere embeddings
- Sentence Transformers
- Hugging Face models

## Prerequisites

1. **ChromaDB Server**: You need a running ChromaDB server
   ```bash
   docker run -p 8000:8000 chromadb/chroma
   ```

2. **No API Keys Required**: This example runs completely locally

## Running the Example

```bash
cd examples/simple-embeddings
go run main.go
```

## What This Example Does

1. **Creates Simple Embeddings**: Uses SHA-256 hash to create deterministic embeddings
2. **Stores Documents**: Adds documents with embeddings to ChromaDB
3. **Performs Searches**: Queries using the embedding vectors
4. **Filters by Metadata**: Demonstrates metadata filtering during search
5. **Shows the Pattern**: Illustrates how you would integrate real embeddings

## Sample Output

```
=== Connecting to ChromaDB ===
ChromaDB Version: 0.5.0

=== Creating Collection ===
Collection: simple_embeddings_demo (ID: 12345-abcde)

=== Generating Simple Embeddings ===
Generating embeddings for 8 documents...
Generated 8 embeddings (dimension: 128)

=== Note About This Demo ===
⚠️  These are hash-based embeddings for demonstration only!
⚠️  They do NOT capture semantic meaning like real embedding models.
⚠️  For production, use OpenAI, Cohere, or Sentence Transformers.

Same text similarity: 1.0000 (should be 1.0)
Similar meaning similarity: 0.1234 (would be high with real embeddings)
...
```

## Why This Example Exists

1. **No API Keys**: You can run this immediately without signing up for services
2. **Understand the Pattern**: See exactly how embeddings flow through the system
3. **Testing**: Useful for testing ChromaDB integration without external dependencies
4. **Educational**: Shows that the library is agnostic to how embeddings are created

## Upgrading to Real Embeddings

To use real semantic embeddings, replace `SimpleEmbeddingGenerator` with:

```go
// Option 1: OpenAI (see examples/with-embeddings)
embeddings := callOpenAIEmbeddingAPI(texts)

// Option 2: Call Python sentence-transformers
embeddings := callSentenceTransformers(texts)

// Option 3: Use any other embedding service
embeddings := yourEmbeddingService.Embed(texts)
```

The rest of your code stays the same!

## Learn More

- See `examples/with-embeddings` for a real OpenAI integration
- [Sentence Transformers](https://www.sbert.net/) for local embedding models
- [OpenAI Embeddings](https://platform.openai.com/docs/guides/embeddings) for cloud-based embeddings
