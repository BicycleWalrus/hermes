package kubeeye

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const watchNamespaceEnv = "WATCH_NAMESPACE"

type inClusterConfigLoader func() (*rest.Config, error)
type kubeconfigBuilder func(masterURL, kubeconfigPath string) (*rest.Config, error)

// Run starts namespace-scoped informers and logs creation events.
func Run(ctx context.Context, logger *slog.Logger) error {
	if logger == nil {
		return errors.New("logger is required")
	}

	namespace, err := watchNamespaceFromEnv(os.Getenv)
	if err != nil {
		return err
	}

	config, err := loadRESTConfig(rest.InClusterConfig, clientcmd.BuildConfigFromFlags, os.Getenv, os.UserHomeDir)
	if err != nil {
		return fmt.Errorf("build kubernetes client config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create kubernetes clientset: %w", err)
	}

	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace(namespace),
	)

	deploymentInformer := factory.Apps().V1().Deployments().Informer()
	statefulSetInformer := factory.Apps().V1().StatefulSets().Informer()
	daemonSetInformer := factory.Apps().V1().DaemonSets().Informer()
	podInformer := factory.Core().V1().Pods().Informer()

	var ready atomic.Bool
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if !ready.Load() {
				return
			}
			logCreation(logger, obj)
		},
	}

	if _, err := deploymentInformer.AddEventHandler(handler); err != nil {
		return fmt.Errorf("register deployment handler: %w", err)
	}
	if _, err := statefulSetInformer.AddEventHandler(handler); err != nil {
		return fmt.Errorf("register statefulset handler: %w", err)
	}
	if _, err := daemonSetInformer.AddEventHandler(handler); err != nil {
		return fmt.Errorf("register daemonset handler: %w", err)
	}
	if _, err := podInformer.AddEventHandler(handler); err != nil {
		return fmt.Errorf("register pod handler: %w", err)
	}

	logger.Info("starting kubeeye watcher", "namespace", namespace)
	factory.Start(ctx.Done())

	if !cache.WaitForCacheSync(
		ctx.Done(),
		deploymentInformer.HasSynced,
		statefulSetInformer.HasSynced,
		daemonSetInformer.HasSynced,
		podInformer.HasSynced,
	) {
		return errors.New("timed out waiting for informer caches to sync")
	}

	ready.Store(true)
	logger.Info("kubeeye watcher ready", "namespace", namespace)

	<-ctx.Done()
	logger.Info("kubeeye watcher stopping", "reason", ctx.Err())
	return nil
}

func watchNamespaceFromEnv(getenv func(string) string) (string, error) {
	namespace := getenv(watchNamespaceEnv)
	if namespace == "" {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnv)
	}
	return namespace, nil
}

func loadRESTConfig(inCluster inClusterConfigLoader, buildConfig kubeconfigBuilder, getenv func(string) string, homeDir func() (string, error)) (*rest.Config, error) {
	config, err := inCluster()
	if err == nil {
		return config, nil
	}

	kubeconfig, pathErr := kubeconfigPath(getenv, homeDir)
	if pathErr != nil {
		return nil, fmt.Errorf("in-cluster config failed: %w; kubeconfig fallback failed: %w", err, pathErr)
	}

	config, buildErr := buildConfig("", kubeconfig)
	if buildErr != nil {
		return nil, fmt.Errorf("in-cluster config failed: %w; kubeconfig %q failed: %w", err, kubeconfig, buildErr)
	}
	return config, nil
}

func kubeconfigPath(getenv func(string) string, homeDir func() (string, error)) (string, error) {
	if kubeconfig := getenv(clientcmd.RecommendedConfigPathEnvVar); kubeconfig != "" {
		return kubeconfig, nil
	}

	home, err := homeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home directory: %w", err)
	}
	return filepath.Join(home, ".kube", "config"), nil
}

func logCreation(logger *slog.Logger, obj interface{}) {
	var (
		kind string
		meta metav1.Object
	)

	switch resource := obj.(type) {
	case *appsv1.Deployment:
		kind = "Deployment"
		meta = resource
	case *appsv1.StatefulSet:
		kind = "StatefulSet"
		meta = resource
	case *appsv1.DaemonSet:
		kind = "DaemonSet"
		meta = resource
	case *corev1.Pod:
		kind = "Pod"
		meta = resource
	default:
		logger.Warn("received unsupported add event object", "type", fmt.Sprintf("%T", obj))
		return
	}

	logger.Info(
		"resource created",
		"kind", kind,
		"namespace", meta.GetNamespace(),
		"name", meta.GetName(),
		"uid", string(meta.GetUID()),
	)
}
