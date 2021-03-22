package instance

import (
	"context"
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	"harbor-operator/pkg/instance/controller"
	"harbor-operator/pkg/instance/helm"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("HarborInstance syncer")

type InstanceSyncer struct {
	Name              string
	Namespace         string
	Client            client.Client
	Scheme            *runtime.Scheme
	HarborServiceInfo *v1.HarborService
	OperatorConfig    *config.ConfigFile
	EventCli          event.Event
}

func NewIntanceSyncer(name, namespace string, hs *v1.HarborService, c client.Client, scheme *runtime.Scheme, opConfig *config.ConfigFile, eventCli event.Event) syncer.SyncInterface {
	return &InstanceSyncer{
		Name:              name,
		Namespace:         namespace,
		Client:            c,
		Scheme:            scheme,
		HarborServiceInfo: hs,
		OperatorConfig:    opConfig,
		EventCli:          eventCli,
	}
}

func (is *InstanceSyncer) GetName() interface{} {
	return is.Name
}

// Sync 用于对harbor实例的操作，如创建，更新等。
func (is *InstanceSyncer) Sync(c context.Context) (ctrl.Result, error) {
	installInfo := fmt.Sprintf("Deploying instance: %s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)
	klog.Info(installInfo)

	switch is.HarborServiceInfo.Status.Condition.Phase {
	case "":
		// 实现创建逻辑
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "deployIntance", "Start deploy harbor service")
		err := is.InstallHarborRelease()
		if err != nil {
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "deployIntanceFailed", "Deploy harbor service by helm failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Deploy harbor service by helm failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}

		klog.Info("Deploy harbor instance successd")
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "deployIntanceSuccessd", "Deploy harbor service by successd.")

		// 部署完成，需要更新cr的状态
		is.HarborServiceInfo.Status.SetRunningStatus(fmt.Sprintf("http://%s.harbor.com:%d", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, is.HarborServiceInfo.Spec.InstanceInfo.NodePortIndex))
		v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
		return ctrl.Result{}, nil

	default:
		// 实现更新逻辑
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "upgradeInstance", "Start upgrade harbor service")
		err := is.UpgradeHarborRelease()
		if err != nil {
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "upgradeIntanceFailed", "Upgrade harbor service by helm failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Upgrade harbor service by helm failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}

		klog.Info("Upgrade harbor instance successd")
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "upgradeIntanceSuccessd", "Upgrade harbor service by successd.")

		//is.HarborServiceInfo.Status.SetRunningStatus(fmt.Sprintf("http://%s.harbor.com:%d", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, is.HarborServiceInfo.Spec.InstanceInfo.NodePortIndex))
		//v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
		return ctrl.Result{}, nil
	}
}

func (is *InstanceSyncer) InstallHarborRelease() error {
	// 根据cr信息，渲染出harbor部署需要的helmchart values数据
	// 目前支持持部署http的harbor
	install := &helm.InstallInfo{
		InstanceName:    is.HarborServiceInfo.Spec.InstanceInfo.InstanceName,
		InstanceVersion: is.HarborServiceInfo.Spec.InstanceInfo.InstanceVersion,
		NodePortIndex:   is.HarborServiceInfo.Spec.InstanceInfo.NodePortIndex,
		S3: helm.S3Config{
			Region:         syncer.DefaultMinioRegion,
			Bucket:         is.HarborServiceInfo.Spec.InstanceInfo.S3Config.Bucket,
			Accesskey:      is.HarborServiceInfo.Spec.InstanceInfo.S3Config.Accesskey,
			SecretKey:      is.HarborServiceInfo.Spec.InstanceInfo.S3Config.SecretKey,
			RegionEndpoint: fmt.Sprintf("%s://%s", is.OperatorConfig.Minio.Proto, is.OperatorConfig.Minio.EndPoint),
		},
		PgInfo: helm.PgInfo{
			Host:                 is.OperatorConfig.DB.PGSQL.Host,
			Port:                 is.OperatorConfig.DB.PGSQL.Port,
			UserName:             is.OperatorConfig.DB.PGSQL.Username,
			Password:             is.OperatorConfig.DB.PGSQL.Password,
			CoreDatabase:         fmt.Sprintf("%s_registry", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			ClairDatabase:        fmt.Sprintf("%s_clair", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			NotaryServerDatabase: fmt.Sprintf("%s_notaryserver", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			NotarySignerDatabase: fmt.Sprintf("%s_notarysigner", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			SSLMode:              is.OperatorConfig.DB.PGSQL.SSLMode,
		},
		RedisInfo: helm.RedisInfo{
			RedisAddr:     is.OperatorConfig.Redis.Addr,
			RedisPassword: is.OperatorConfig.Redis.Password,
			RedisDbIndex:  is.HarborServiceInfo.Spec.InstanceInfo.RedisDbIndex,
		},
		Jobservicepvc: fmt.Sprintf("pvc-%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
	}
	values := controller.GenerageInstanceHelmChart(install)

	// 将values 更新到最新的helm-chart中
	var targetPath = is.OperatorConfig.HarborHelmPath.HarborV213
	err := controller.SaveValuesToFile(values, targetPath+"/values.yaml")
	if err != nil {
		klog.Error(fmt.Sprintf("fail to save values.yaml file"))
		return err
	}

	// 调用helm 方法，创建实例
	listOps := &helm.ReleaseListOptions{
		Namespace: is.HarborServiceInfo.Spec.InstanceInfo.InstanceName,
	}
	installOps := &helm.ReleaseOptions{}
	settings := cli.New()

	kubeToken := &helm.KubeToken{
		ApiServer: is.OperatorConfig.Kubernetes.Apiserver,
		Token:     is.OperatorConfig.Kubernetes.Token,
	}

	err = controller.InstallRelease(kubeToken, listOps, installOps, settings, is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, targetPath)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to install release:%s, witherr:%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, err.Error()))
		return err
	}

	return nil
}

func (is *InstanceSyncer) UpgradeHarborRelease() error {
	// 根据cr信息，渲染出harbor部署需要的helmchart values数据
	// 目前支持持部署http的harbor
	install := &helm.InstallInfo{
		InstanceName:    is.HarborServiceInfo.Spec.InstanceInfo.InstanceName,
		InstanceVersion: is.HarborServiceInfo.Spec.InstanceInfo.InstanceVersion,
		NodePortIndex:   is.HarborServiceInfo.Spec.InstanceInfo.NodePortIndex,
		S3: helm.S3Config{
			Region:         syncer.DefaultMinioRegion,
			Bucket:         is.HarborServiceInfo.Spec.InstanceInfo.S3Config.Bucket,
			Accesskey:      is.HarborServiceInfo.Spec.InstanceInfo.S3Config.Accesskey,
			SecretKey:      is.HarborServiceInfo.Spec.InstanceInfo.S3Config.SecretKey,
			RegionEndpoint: fmt.Sprintf("%s://%s", is.OperatorConfig.Minio.Proto, is.OperatorConfig.Minio.EndPoint),
		},
		PgInfo: helm.PgInfo{
			Host:                 is.OperatorConfig.DB.PGSQL.Host,
			Port:                 is.OperatorConfig.DB.PGSQL.Port,
			UserName:             is.OperatorConfig.DB.PGSQL.Username,
			Password:             is.OperatorConfig.DB.PGSQL.Password,
			CoreDatabase:         fmt.Sprintf("%s_registry", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			ClairDatabase:        fmt.Sprintf("%s_clair", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			NotaryServerDatabase: fmt.Sprintf("%s_notaryserver", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			NotarySignerDatabase: fmt.Sprintf("%s_notarysigner", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
			SSLMode:              is.OperatorConfig.DB.PGSQL.SSLMode,
		},
		RedisInfo: helm.RedisInfo{
			RedisAddr:     is.OperatorConfig.Redis.Addr,
			RedisPassword: is.OperatorConfig.Redis.Password,
			RedisDbIndex:  is.HarborServiceInfo.Spec.InstanceInfo.RedisDbIndex,
		},
		Jobservicepvc: fmt.Sprintf("pvc-%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName),
	}
	values := controller.GenerageInstanceHelmChart(install)

	// 将values 更新到最新的helm-chart中
	var targetPath = is.OperatorConfig.HarborHelmPath.HarborV213
	err := controller.SaveValuesToFile(values, targetPath+"/values.yaml")
	if err != nil {
		klog.Error(fmt.Sprintf("fail to save values.yaml file"))
		return err
	}

	// 调用helm 方法，创建实例
	listOps := &helm.ReleaseListOptions{
		Namespace: is.HarborServiceInfo.Spec.InstanceInfo.InstanceName,
	}
	installOps := &helm.ReleaseOptions{}
	settings := cli.New()

	kubeToken := &helm.KubeToken{
		ApiServer: is.OperatorConfig.Kubernetes.Apiserver,
		Token:     is.OperatorConfig.Kubernetes.Token,
	}

	err = controller.UpgradeRelease(kubeToken, listOps, installOps, settings, is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, targetPath)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to install release:%s, witherr:%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, err.Error()))
		return err
	}

	return nil
}

func (is *InstanceSyncer) Delete(c context.Context) error {
	// 调用helm 方法，删除实例
	listOps := &helm.ReleaseListOptions{
		Namespace: is.HarborServiceInfo.Spec.InstanceInfo.InstanceName,
	}
	installOps := &helm.ReleaseOptions{}
	settings := cli.New()

	kubeToken := &helm.KubeToken{
		ApiServer: is.OperatorConfig.Kubernetes.Apiserver,
		Token:     is.OperatorConfig.Kubernetes.Token,
	}

	err := controller.UninstallRelease(kubeToken, listOps, installOps, settings, is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to uninstall release:%s, witherr: %s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, err.Error()))
		return err
	}

	klog.Info("Remove harborservice release success")
	return nil
}
