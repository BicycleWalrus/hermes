# Hermes Playground - Secure Kubernetes Tooling

Welcome to the **Hermes Playground**. This repository is dedicated to building secure, lightweight, and performant Go utilities designed to run within Kubernetes clusters. 

This playground focuses on developing security systems and native cluster utilities compiled into minimal, zero-dependency `FROM scratch` container images. This approach ensures an extremely small attack surface and aligns with secure-by-default philosophies (such as bootstrapping clusters the "Alta3 Way").

---

## 🎯 Repository Objectives

1. **Secure-by-Design Go Applications**: Compile highly-optimized, statically-linked Go binaries that run inside pure `scratch` containers.
2. **Kubernetes-Native Integration**: Implement clean, secure communications with the internal Kubernetes API Server and resources (such as the `metrics-server`).
3. **Rigorous RBAC Governance**: Ensure all utilities operate under the principle of least privilege using precise `ServiceAccount`, `ClusterRole`, and `ClusterRoleBinding` configurations.
4. **Hands-on Demonstration**: Showcase autonomous agent-driven development using Hermes to solve complex containerization, deployment, and security-compliance workflows.

---

## 🏗️ Monorepo Directory Structure

The repository follows a standard Go monorepo design, allowing us to manage multiple tools, shared libraries, and Kubernetes manifests in a clean structure:

```text
hermes/
├── cmd/                  # Main entry points for each Go application
│   └── kubeeye/          # KubeEye topology and metrics backend
├── pkg/                  # Reusable, shared library code
│   └── k8sclient/        # Common Kubernetes API helper packages
├── deployments/          # Kubernetes manifests (YAML files)
│   ├── rbac.yaml         # Least-privilege roles and bindings
│   └── deployment.yaml   # Service and Deployment specs
├── build/                # Dockerfiles and build automation scripts
│   └── kubeeye.Dockerfile
├── LICENSE               # MIT License
└── README.md             # This documentation
```

---

## 🚀 Key Feature: KubeEye (Interactive Cluster Topology Map)

During the playground demo, we will build **KubeEye**, a real-time cluster visualization tool that acts as a secure dashboard showing cluster topology.

### Concept & Design
* **Real-Time Map**: An interactive frontend dashboard displaying the topology of Nodes, Namespaces, Pods, and Services.
* **Live Performance Metrics**: Queries the cluster's internal `metrics-server` to stream live CPU and Memory utilization gauge bars per pod and node.
* **Secure RBAC Boundary**: Accesses the cluster API securely using read-only RBAC credentials mapped to its running Pod.

### Deployment & Access
To interact with KubeEye from the BCHD workstation as the `student` user, run:
```bash
kubectl port-forward svc/kubeeye 2225:80 --address=0.0.0.0
```
Then, access the dashboard in your web browser at: `http://localhost:2225` (or the workstation's IP on port 2225).

---

## 🐳 Containerization & Registry Standards

All binaries compiled in this playground are packaged into zero-vulnerability container images utilizing multi-stage builds.

### Standard Go Scratch Dockerfile
Below is the standard Dockerfile pattern for Go utilities in this workspace:

```dockerfile
FROM golang:1.21.4 AS builder

WORKDIR /app
COPY . .

# Compile the static, self-contained binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-binary ./cmd/kubeeye

# Final minimal stage
FROM scratch

# Copy CA certificates from builder for SSL/TLS verification
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy compiled static binary
COPY --from=builder /app/app-binary /app-binary

# Expose target listening port
EXPOSE 8888

# Run binary as the container entrypoint
CMD ["/app-binary"]
```

### Publishing to Registry
We publish compiled images to the local **BCHD registry** (`bchd.registry`), which is pre-configured on our cluster.

To push an image to the private registry:
```bash
sudo docker build -t bchd.registry/kubeeye:v1 -f build/kubeeye.Dockerfile .
sudo docker push bchd.registry/kubeeye:v1
```

---

## 🛠️ Local Workstation & Development Environment

* **Workstation**: Alta3 BCHD Workstation (logged in as the `student` user).
* **Target Cluster**: Active Kubernetes cluster configured with native API access.
* **Docker Registry**: `bchd.registry`

Let the building begin! 🚀
