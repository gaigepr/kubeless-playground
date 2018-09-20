package kubeless

import (
	"encoding/json"
	"os"

	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	"github.com/kubeless/kubeless/pkg/functions"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NOTE: Reliant on the default:default service account having a cluster-admin clusterrolebinding.

var (
	argoNamespace string = os.Getenv("ARGO_NAMESPACE")

	restConfig   *rest.Config
	clientset    *kubernetes.Clientset
	clientConfig clientcmd.ClientConfig
	wfClient     v1alpha1.WorkflowInterface
)

func initKubeAndArgoClients(namespace string) (err error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := clientcmd.ConfigOverrides{}
	_ = clientcmd.RecommendedConfigOverrideFlags("")
	clientConfig = clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)

	restConfig, err = clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	wfClient = wfclientset.NewForConfigOrDie(restConfig).ArgoprojV1alpha1().Workflows(namespace)

	return nil
}

func Handler(event functions.Event, context functions.Context) (string, error) {
	var err error

	err = initKubeAndArgoClients(argoNamespace)
	if err != nil {
		return "error1", err
	}

	wfs, err := wfClient.List(metav1.ListOptions{})
	if err != nil {
		return "error2", err
	}

	buf, err := json.Marshal(wfs.Items)
	if err != nil {
		return "error3", err
	}

	return string(buf), nil
}
