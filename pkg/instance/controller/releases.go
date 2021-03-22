package controller

import (
	"harbor-operator/pkg/instance/helm"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/klog"
)

func InstallRelease(kubeToken *helm.KubeToken, releaseListOps *helm.ReleaseListOptions, installOps *helm.ReleaseOptions, settings *cli.EnvSettings, releaseName, chartPath string) error {
	releaseManager := helm.NewHelmReleaseManager(kubeToken, releaseListOps, installOps, settings)
	err := releaseManager.InstallRelease(releaseName, chartPath)
	if err != nil {
		klog.Errorf("Fail to install release of %s, with error: %s", releaseName, err.Error())
		return err
	}
	return nil
}

func UpgradeRelease(kubeToken *helm.KubeToken, releaseListOps *helm.ReleaseListOptions, installOps *helm.ReleaseOptions, settings *cli.EnvSettings, releaseName, chartPath string) error {
	releaseManager := helm.NewHelmReleaseManager(kubeToken, releaseListOps, installOps, settings)
	err := releaseManager.UpgradeRelease(releaseName, chartPath)
	if err != nil {
		klog.Errorf("Fail to install release of %s, with error: %s", releaseName, err.Error())
		return err
	}
	return nil
}

func ListReleases(kubeToken *helm.KubeToken, releaseListOps *helm.ReleaseListOptions, installOps *helm.ReleaseOptions, settings *cli.EnvSettings) ([]helm.ReleaseElement, error) {
	var targetNamespace string
	if releaseListOps.Namespace != "" {
		targetNamespace = releaseListOps.Namespace
	} else if releaseListOps.AllNamespaces == true {
		targetNamespace = ""
	} else {
		targetNamespace = "default"
	}
	releaseListOps.Namespace = targetNamespace

	releaseManager := helm.NewHelmReleaseManager(kubeToken, releaseListOps, installOps, settings)

	releases, err := releaseManager.ListReleases()
	if err != nil {
		klog.Errorf("Fail to list releases, with error: %s", err.Error())
		return nil, err
	}
	return releases, nil
}

func GetReleaseStatus(kubeToken *helm.KubeToken, releaseListOps *helm.ReleaseListOptions, installOps *helm.ReleaseOptions, settings *cli.EnvSettings, name string) (helm.ReleaseElement, error) {
	releaseManager := helm.NewHelmReleaseManager(kubeToken, releaseListOps, installOps, settings)
	release, err := releaseManager.GetReleaseStatus(name)
	if err != nil {
		klog.Errorf("Fail to get release status, with error: %s", err.Error())
		return helm.ReleaseElement{}, err
	}
	return release, nil
}

func GetReleaseHistry(kubeToken *helm.KubeToken, releaseListOps *helm.ReleaseListOptions, installOps *helm.ReleaseOptions, settings *cli.EnvSettings, name string) (helm.ReleaseHistory, error) {
	releaseManager := helm.NewHelmReleaseManager(kubeToken, releaseListOps, installOps, settings)
	releaseHistrory, err := releaseManager.ListReleaseHistories(name)
	if err != nil {
		klog.Errorf("Fail to get release histrory, with error: %s", err.Error())
		return helm.ReleaseHistory{}, err
	}
	return releaseHistrory, nil
}

func UninstallRelease(kubeToken *helm.KubeToken, releaseListOps *helm.ReleaseListOptions, installOps *helm.ReleaseOptions, settings *cli.EnvSettings, name string) error {
	releaseManager := helm.NewHelmReleaseManager(kubeToken, releaseListOps, installOps, settings)
	err := releaseManager.UninstallRelease(name)
	if err != nil {
		klog.Errorf("Fail to uninstall release of %s, with error: %s", name, err.Error())
		return err
	}
	return nil
}
