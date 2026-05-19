# Embedder Tool

Embedder is a standalone gateway that converts input text into vectors via OpenAI and forwards vector operations to Vald.
It provides gRPC and REST endpoints for search, insert, update, upsert, remove, and embedding.

## Required Configuration

Minimum keys:

```yaml
server:
  grpc:
    port: 8081
  http:
    port: 8080
client:
  addrs:
    - "vald-gateway.default.svc.cluster.local:8081"
llm:
  provider: openai
  openai:
    token: "${OPENAI_API_KEY}"
    model: "small3" # adav2 | small3 | large3
```

Optional metadata client:

```yaml
meta:
  client:
    addrs:
      - "vald-meta-gateway.default.svc.cluster.local:8081"
```

## Run

```bash
make cmd/tools/embedder/build
./bin/tools/embedder -c /path/to/embedder.yaml
```

## API Examples

Embedding:

```bash
curl -sS -X POST http://127.0.0.1:8080/embedding \
  -H 'Content-Type: application/json' \
  -d '{"text":"hello vald"}'
```

Insert:

```bash
curl -sS -X POST http://127.0.0.1:8080/insert \
  -H 'Content-Type: application/json' \
  -d '{"document":{"id":"doc-1","text":"hello vald"}}'
```

Search:

```bash
curl -sS -X POST http://127.0.0.1:8080/search \
  -H 'Content-Type: application/json' \
  -d '{"text":"hello"}'
```
