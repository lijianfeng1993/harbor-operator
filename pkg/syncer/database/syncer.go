package database

import (
	"context"
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	"harbor-operator/pkg/dao"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DatabaseSyncer struct {
	Name              string
	HarborServiceInfo *v1.HarborService
	Client            client.Client
	OperatorConfig    *config.ConfigFile
	EventCli          event.Event
}

// NewDatabaseSyner 创建数据库的同步syncer
func NewDatabaseSyncer(name string, hs *v1.HarborService, c client.Client, opConfig *config.ConfigFile, eventCli event.Event) syncer.SyncInterface {
	return &DatabaseSyncer{
		Name:              name,
		HarborServiceInfo: hs,
		Client:            c,
		OperatorConfig:    opConfig,
		EventCli:          eventCli,
	}
}

func (is *DatabaseSyncer) GetName() interface{} {
	return is.Name
}

// Sync 用于对harbor实例的操作，如创建，更新等。
func (is *DatabaseSyncer) Sync(c context.Context) (ctrl.Result, error) {
	switch is.HarborServiceInfo.Status.Condition.Phase {
	case "":
		// 实现初始化pg数据库
		installInfo := fmt.Sprintf("sync database: %s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)
		klog.Info(installInfo)
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "syncDatabase", "Start sync pgsql database and table")

		// 参考https://github.com/goharbor/harbor/tree/master/make/photon/db中的*.sql， 完成pgsql数据的初始化
		err := dao.SyncDatabases(is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, is.OperatorConfig)
		if err != nil {
			// 数据库初始化失败，更新cr status和event
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "syncDatabaseFailed", "Sync pgsql database failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Sync pgsql database failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}

		err = dao.SyncTables(is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, is.OperatorConfig)
		if err != nil {
			// 表初始化失败，更新cr status和event
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "syncDatabaseFailed", "Sync pgsql table failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Sync pgsql table failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}

		// 打印日志并且记录event
		// 数据库创建成功不会更新cr的status,只有最后一步的instance的创建成功才会更新cr的status
		klog.Info("Sync pgsql successed.")
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "syncDatabaseSuccess", "Sync pgsql database success")
	}

	return ctrl.Result{}, nil
}

func (is *DatabaseSyncer) Delete(c context.Context) error {
	err := dao.DeleteDatabases(is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, is.OperatorConfig)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to clear databases, witherr: %s", err.Error()))
		return err
	}

	klog.Info("Remove pgsql database success")
	return nil
}
