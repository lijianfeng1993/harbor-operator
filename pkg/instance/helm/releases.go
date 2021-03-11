package helm

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/klog"
	"sigs.k8s.io/yaml"
	"strconv"

	"github.com/pkg/errors"
)

type HelmReleaseManager struct {
	kubeToken    *KubeToken
	listOps      *ReleaseListOptions
	releaseOps   *ReleaseOptions
	helmSettings *cli.EnvSettings
}

func NewHelmReleaseManager(cfg *KubeToken, listOps *ReleaseListOptions, releaseOps *ReleaseOptions, settings *cli.EnvSettings) *HelmReleaseManager {
	return &HelmReleaseManager{
		kubeToken:    cfg,
		listOps:      listOps,
		releaseOps:   releaseOps,
		helmSettings: settings,
	}
}

func (hRelease *HelmReleaseManager) ListReleases() ([]ReleaseElement, error) {
	// 获取actionConfig
	actionConfig, err := getActionConfig(hRelease.listOps.Namespace, hRelease.kubeToken)
	if err != nil {
		klog.Errorf("Fail to get actionConfig, with error: %s", err.Error())
		return nil, err
	}

	// 创建list client
	client := action.NewList(actionConfig)
	client.ByDate = hRelease.listOps.ByDate
	client.SortReverse = hRelease.listOps.SortReverse
	client.Limit = hRelease.listOps.Limit
	client.Offset = hRelease.listOps.Offset
	client.Filter = hRelease.listOps.Filter
	client.Uninstalled = hRelease.listOps.Uninstalled
	client.Superseded = hRelease.listOps.Superseded
	client.Uninstalling = hRelease.listOps.Uninstalling
	client.Deployed = hRelease.listOps.Deployed
	client.Failed = hRelease.listOps.Failed
	client.Pending = hRelease.listOps.Pending
	client.SetStateMask()

	results, err := client.Run()
	if err != nil {
		klog.Errorf("Fail to run helm release list client, with error: %s", err.Error())
		return nil, err
	}

	// Initialize the array so no results returns an empty array instead of null
	elements := make([]ReleaseElement, 0, len(results))
	for _, r := range results {
		elements = append(elements, constructReleaseElement(r, false))
	}

	return elements, nil
}

func (hRelease *HelmReleaseManager) GetReleaseStatus(name string) (ReleaseElement, error) {
	// 获取actionConfig
	actionConfig, err := getActionConfig(hRelease.listOps.Namespace, hRelease.kubeToken)
	if err != nil {
		klog.Errorf("Fail to get actionConfig, with error: %s", err.Error())
		return ReleaseElement{}, err
	}
	client := action.NewStatus(actionConfig)
	result, err := client.Run(name)
	if err != nil {
		klog.Errorf("Fail to run helm release getstatus client, with error: %s", err.Error())
		return ReleaseElement{}, err
	}
	element := constructReleaseElement(result, true)
	return element, nil
}

func (hRelease *HelmReleaseManager) ListReleaseHistories(name string) (ReleaseHistory, error) {
	actionConfig, err := getActionConfig(hRelease.listOps.Namespace, hRelease.kubeToken)
	if err != nil {
		klog.Errorf("Fail to get actionConfig, with error: %s", err.Error())
		return ReleaseHistory{}, err
	}

	client := action.NewHistory(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		klog.Errorf("Fail to get release history of %s, with error: %s", name, err.Error())
		return ReleaseHistory{}, err
	}

	return getReleaseHistory(results), nil
}

func (hRelease *HelmReleaseManager) InstallRelease(name, chart string) error {
	vals, err := mergeValues(hRelease.releaseOps)
	if err != nil {
		klog.Errorf("Fail to mergeValues, with error: %s", err.Error())
		return err
	}

	actionConfig, err := getActionConfig(hRelease.listOps.Namespace, hRelease.kubeToken)
	if err != nil {
		klog.Errorf("Fail to get actionConfig, with error: %s", err.Error())
		return err
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = name
	client.Namespace = hRelease.listOps.Namespace

	// merge install options
	client.DryRun = hRelease.releaseOps.DryRun
	client.DisableHooks = hRelease.releaseOps.DisableHooks
	client.Wait = hRelease.releaseOps.Wait
	client.Devel = hRelease.releaseOps.Devel
	client.Description = hRelease.releaseOps.Description
	client.Atomic = hRelease.releaseOps.Atomic
	client.SkipCRDs = hRelease.releaseOps.SkipCRDs
	client.SubNotes = hRelease.releaseOps.SubNotes
	client.Timeout = hRelease.releaseOps.Timeout
	client.CreateNamespace = hRelease.releaseOps.CreateNamespace
	client.DependencyUpdate = hRelease.releaseOps.DependencyUpdate

	// merge chart path options
	client.ChartPathOptions.CaFile = hRelease.releaseOps.ChartPathOptions.CaFile
	client.ChartPathOptions.CertFile = hRelease.releaseOps.ChartPathOptions.CertFile
	client.ChartPathOptions.KeyFile = hRelease.releaseOps.ChartPathOptions.KeyFile
	client.ChartPathOptions.InsecureSkipTLSverify = hRelease.releaseOps.ChartPathOptions.InsecureSkipTLSverify
	client.ChartPathOptions.Keyring = hRelease.releaseOps.ChartPathOptions.Keyring
	client.ChartPathOptions.Password = hRelease.releaseOps.ChartPathOptions.Password
	client.ChartPathOptions.RepoURL = hRelease.releaseOps.ChartPathOptions.RepoURL
	client.ChartPathOptions.Username = hRelease.releaseOps.ChartPathOptions.Username
	client.ChartPathOptions.Verify = hRelease.releaseOps.ChartPathOptions.Verify
	client.ChartPathOptions.Version = hRelease.releaseOps.ChartPathOptions.Version

	chartPath, err := client.ChartPathOptions.LocateChart(chart, hRelease.helmSettings)
	if err != nil {
		klog.Errorf("Fail to get chartPath, with error: %s", err.Error())
		return err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		klog.Errorf("Fail to load chartPath into chartRequested, with error: %s", err.Error())
		return err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		klog.Errorf("Fail to check if chart is installable, with error: %s", err.Error())
		return err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					ChartPath:        chartPath,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          getter.All(hRelease.helmSettings),
					RepositoryConfig: hRelease.helmSettings.RepositoryConfig,
					RepositoryCache:  hRelease.helmSettings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					klog.Errorf("Fail to update chart, with error: %s", err.Error())
					return err
				}
			} else {
				klog.Errorf("Fail to check dependencies, with error: %s", err.Error())
				return err
			}
		}
	}

	_, err = client.Run(chartRequested, vals)
	if err != nil {
		klog.Errorf("")
		return err
	}

	return nil
}

func (hRelease *HelmReleaseManager) UninstallRelease(name string) error {
	actionConfig, err := getActionConfig(hRelease.listOps.Namespace, hRelease.kubeToken)
	if err != nil {
		klog.Errorf("Fail to get actionConfig, with error: %s", err.Error())
		return err
	}

	client := action.NewUninstall(actionConfig)
	_, err = client.Run(name)
	if err != nil {
		klog.Errorf("Fail to uninstall release of %s, with error: %s", name, err.Error())
		return err
	}
	return nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}

	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func mergeValues(options *ReleaseOptions) (map[string]interface{}, error) {
	vals := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(options.Values), &vals)
	if err != nil {
		return vals, fmt.Errorf("failed parsing values")
	}

	for _, value := range options.SetValues {
		if err := strvals.ParseInto(value, vals); err != nil {
			return vals, fmt.Errorf("failed parsing set data")
		}
	}

	for _, value := range options.SetStringValues {
		if err := strvals.ParseIntoString(value, vals); err != nil {
			return vals, fmt.Errorf("failed parsing set_string data")
		}
	}

	return vals, nil
}

// 格式化release记录
func constructReleaseElement(r *release.Release, showStatus bool) ReleaseElement {
	element := ReleaseElement{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Revision:     strconv.Itoa(r.Version),
		Status:       r.Info.Status.String(),
		Chart:        r.Chart.Metadata.Name,
		ChartVersion: r.Chart.Metadata.Version,
		AppVersion:   r.Chart.Metadata.AppVersion,
	}
	if showStatus {
		element.Notes = r.Info.Notes
	}
	t := "-"
	if tspb := r.Info.LastDeployed; !tspb.IsZero() {
		t = tspb.String()
	}
	element.Updated = t
	return element
}

func getReleaseHistory(rls []*release.Release) (history ReleaseHistory) {
	for i := len(rls) - 1; i >= 0; i-- {
		r := rls[i]
		c := formatChartname(r.Chart)
		s := r.Info.Status.String()
		v := r.Version
		d := r.Info.Description
		a := formatAppVersion(r.Chart)

		rInfo := ReleaseInfo{
			Revision:    v,
			Status:      s,
			Chart:       c,
			AppVersion:  a,
			Description: d,
		}
		if !r.Info.LastDeployed.IsZero() {
			rInfo.Updated = r.Info.LastDeployed

		}
		history = append(history, rInfo)
	}

	return history
}

func formatChartname(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return fmt.Sprintf("%s-%s", c.Name(), c.Metadata.Version)
}

func formatAppVersion(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return c.AppVersion()
}
