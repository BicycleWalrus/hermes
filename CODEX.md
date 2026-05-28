# Codex Developer Guide - Hermes Playground

This guide defines the development commands, build flags, and environment standards for **ChatGPT Codex 5.5** (and other AI code-writing tools) working in this repository.

---

## 🛠️ Codex Developer Commands

### Environment
Always use the installed Go compiler at `/usr/local/go/bin/go`. Prepend `/usr/local/go/bin` to your `PATH` if it is not already available:
```bash
export PATH=/usr/local/go/bin:$PATH
```

### Build & Compilation (Static Linked)
To compile a Go utility (e.g., `kubeeye`) as a statically-linked binary with CGO disabled:
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubeeye ./cmd/kubeeye
```

### Code Quality & Testing
* **Format**: `go fmt ./...`
* **Static Analysis**: `go vet ./...`
* **Run Tests**: `go test ./...`

---

## 🐳 Containerization & Registry Standards

### Multi-Stage Dockerfile Pattern
All Dockerfiles must utilize multi-stage builds and run on top of an empty `scratch` image.
```dockerfile
FROM golang:1.21.4 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-binary ./cmd/kubeeye

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app-binary /app-binary
EXPOSE 8888
CMD ["/app-binary"]
```

### Registry Publishing
We compile and push container images to our local private registry (`bchd.registry`):
```bash
sudo docker build -t bchd.registry/kubeeye:v1 -f build/kubeeye.Dockerfile .
sudo docker push bchd.registry/kubeeye:v1
```

---

## ☸️ Kubernetes Deployment & Access

All deployments must operate securely on the cluster:
1. **RBAC**: Define precise, read-only ServiceAccounts, ClusterRoles, and ClusterRoleBindings in `/deployments/rbac.yaml`.
2. **Access**: Expose via Service, and access from the BCHD workstation using:
   ```bash
   kubectl port-forward svc/kubeeye 2225:80 --address=0.0.0.0
   ```
