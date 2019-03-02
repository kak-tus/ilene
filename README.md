# Ilene

UI for Redis Cluster Streams (queues).

Also this is Prometheus exporter for Redis Cluster Streams metrics.

## Run

```
docker run --rm -p 8080:8080 ILENE_REDIS_ADDRS=redis-cluster.example.com:7000 kaktuss/ilene
```

## API docs

Open http://127.0.0.1:8080/docs/ at browser. It's OpenAPI specification of UI API.

## Metrics

http://127.0.0.1:8080/metrics
