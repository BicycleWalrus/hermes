# Agent Guidelines - Secure Kubernetes Tooling Monorepo

This repository is a playground for developing secure, lightweight, and performant Go utilities (such as **KubeEye**) designed to run within Kubernetes clusters. 

All AI agents (Hermes, Codex, and others) collaborating on this repository must strictly adhere to these architectural and behavioral guidelines.

---

## 🎯 Architectural Principles

### 1. Minimal Attack Surface (`FROM scratch`)
* All applications must compile into statically-linked, self-contained binaries and run in pure, empty `FROM scratch` containers.
* Zero-vulnerability design: No base OS, package managers, or shells in the final runtime container.

### 2. Least-Privilege RBAC Authorization
* Security is paramount. Do not request wildcard permissions (`*`) for resources or verbs in Kubernetes manifests.
* Every utility requiring Kubernetes API access must have corresponding, precise, read-only RBAC manifests (`ServiceAccount`, `ClusterRole`, `ClusterRoleBinding`) matching the principle of least privilege.

### 3. Monorepo Organization
Maintain a clean, standard Go monorepo layout:
* `/cmd/<tool-name>`: Entrypoints for compileable binaries (e.g., `/cmd/kubeeye`).
* `/pkg/<lib-name>`: Reusable, shared library code.
* `/deployments`: Kubernetes YAML manifests (RBAC, Services, Deployments).
* `/build`: Dockerfiles named `<tool-name>.Dockerfile` (e.g., `kubeeye.Dockerfile`).

---

## 🚀 Execution & Verification Workflow

Agents must verify and test changes iteratively before concluding work:
1. **Compilation Check**: Ensure the statically linked binary compiles successfully.
2. **Quality Gates**: Run Go formatting and static analysis (`go fmt` and `go vet`).
3. **Unit Testing**: Run and maintain the test suite with `go test ./...`.
4. **Local Port-Forward Access**: Utilities (like KubeEye) will be accessed from the BCHD workstation as the `student` user. Verify that port-forwarding instructions match:
   `kubectl port-forward svc/kubeeye 2225:80 --address=0.0.0.0`
