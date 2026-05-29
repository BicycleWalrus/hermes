package kubeeye

import (
	"errors"
	"testing"

	"k8s.io/client-go/rest"
)

func TestWatchNamespaceFromEnvRequiresValue(t *testing.T) {
	_, err := watchNamespaceFromEnv(func(string) string { return "" })
	if err == nil {
		t.Fatal("expected error when WATCH_NAMESPACE is unset")
	}
}

func TestWatchNamespaceFromEnvReturnsValue(t *testing.T) {
	got, err := watchNamespaceFromEnv(func(string) string { return "team-a" })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "team-a" {
		t.Fatalf("expected namespace team-a, got %q", got)
	}
}

func TestKubeconfigPathUsesEnvironmentOverride(t *testing.T) {
	got, err := kubeconfigPath(func(string) string { return "/tmp/custom-kubeconfig" }, func() (string, error) {
		return "/home/student", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/tmp/custom-kubeconfig" {
		t.Fatalf("expected override path, got %q", got)
	}
}

func TestKubeconfigPathFallsBackToHomeDirectory(t *testing.T) {
	got, err := kubeconfigPath(func(string) string { return "" }, func() (string, error) {
		return "/home/student", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/home/student/.kube/config" {
		t.Fatalf("expected home fallback path, got %q", got)
	}
}

func TestLoadRESTConfigPrefersInClusterConfig(t *testing.T) {
	inClusterCalls := 0
	buildCalls := 0

	cfg, err := loadRESTConfig(
		func() (*rest.Config, error) {
			inClusterCalls++
			return &rest.Config{Host: "https://cluster.example"}, nil
		},
		func(string, string) (*rest.Config, error) {
			buildCalls++
			return nil, errors.New("should not be called")
		},
		func(string) string { return "" },
		func() (string, error) { return "/home/student", nil },
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "https://cluster.example" {
		t.Fatalf("expected in-cluster config host, got %q", cfg.Host)
	}
	if inClusterCalls != 1 {
		t.Fatalf("expected 1 in-cluster call, got %d", inClusterCalls)
	}
	if buildCalls != 0 {
		t.Fatalf("expected 0 kubeconfig builder calls, got %d", buildCalls)
	}
}

func TestLoadRESTConfigFallsBackToKubeconfig(t *testing.T) {
	cfg, err := loadRESTConfig(
		func() (*rest.Config, error) {
			return nil, errors.New("not in cluster")
		},
		func(_ string, path string) (*rest.Config, error) {
			return &rest.Config{Host: path}, nil
		},
		func(string) string { return "" },
		func() (string, error) { return "/home/student", nil },
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "/home/student/.kube/config" {
		t.Fatalf("expected fallback kubeconfig path, got %q", cfg.Host)
	}
}
