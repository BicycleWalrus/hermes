# Hermes Playground - Secure Kubernetes Tooling

KubeEye is a small Go watcher built with client-go informers. It watches creation events for Deployments, StatefulSets, DaemonSets, and Pods inside the namespace named by `WATCH_NAMESPACE`, then emits structured JSON logs to stdout.

## Development

Local development falls back to `KUBECONFIG` or `$HOME/.kube/config` when in-cluster credentials are unavailable.

Run the quality gates:

```bash
export PATH=/usr/local/go/bin:$PATH
go fmt ./...
go vet ./...
go test ./...
```

Run locally:

```bash
export WATCH_NAMESPACE=default
export KUBECONFIG=$HOME/.kube/config
go run ./cmd/kubeeye
```

## Container

Build the static binary into a scratch image:

```bash
docker build -t bchd.registry/kubeeye:v1 -f build/kubeeye.Dockerfile .
```

## Kubernetes manifests

- `deployments/rbac.yaml`: namespace-scoped read-only access to Pods, Deployments, StatefulSets, and DaemonSets.
- `deployments/deployment.yaml`: single-replica deployment wiring `WATCH_NAMESPACE` to the pod namespace.
