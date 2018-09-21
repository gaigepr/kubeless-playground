package kubeless

import (
	"os"

	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"

	"github.com/kubeless/kubeless/pkg/functions"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	wfNamespace      string = os.Getenv("WORKFLOW_NAMESPACE")
	wfServiceAccount string = os.Getenv("WORKFLOW_SERVICE_ACCOUNT")

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

func decodeWorkflowYaml(data string) (wf *wfv1.Workflow) {
	wf = new(wfv1.Workflow)
	// TODO
	return
}

func Handler(event functions.Event, context functions.Context) (string, error) {
	var err error

	err = initKubeAndArgoClients(argoNamespace)
	if err != nil {
		return "error1", err
	}

	var wf *wfv1.Workflow
	wf, err = decodeWorkflowYaml(event.Data)
	if err != nil {
		return "error2", err
	}

	wf, err = wfClient.Create(wf)
	if err != nil {
		return "error3", err
	}

	return "ok", nil
}
