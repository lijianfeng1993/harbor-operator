package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"
	"log"
	"os"
)

func getActionConfig(namespace string, kubetoken *KubeToken) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	kubeConfig := &genericclioptions.ConfigFlags{
		APIServer:   &kubetoken.ApiServer,
		BearerToken: &kubetoken.Token,
		Namespace:   &namespace,
		Insecure:    boolptr(true),
	}

	if err := actionConfig.Init(kubeConfig, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		klog.Errorf("Fail to init actionConfig, with error: %s", err.Error())
		return nil, err
	}
	return actionConfig, nil
}

func boolptr(val bool) *bool {
	return &val
}
