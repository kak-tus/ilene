# Ilene

UI for Redis Cluster Streams (queues).

Also this is Prometheus exporter for Redis Cluster Streams metrics.

## Run

```
docker run --rm -p 8080:8080 -e ILENE_REDIS_ADDRS=redis://:pass@redis-cluster-1.example.com:7000,redis://:pass@redis-cluster-2.example.com:7000 kaktuss/ilene
```

## Configuration

ILENE_REDIS_ADDRS - Comma separated urls (one or more), to connect to cluster nodes. Password is optional.

## API docs

Open http://127.0.0.1:8080/docs/ at browser. It's OpenAPI specification of UI API.

## Metrics

http://127.0.0.1:8080/metrics
