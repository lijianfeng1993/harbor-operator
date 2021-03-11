package s3

import (
	"context"
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	operator_s3 "harbor-operator/pkg/storage/s3"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MinioSyncer struct {
	Name          string
	BucketName    string
	HarborService *v1.HarborService
	Client        client.Client
	MinioClient   *operator_s3.MinioClient
	EventCli      event.Event
}

func NewMinioSyncer(name string, hs *v1.HarborService, c client.Client, opConfig *config.ConfigFile, eventCli event.Event) syncer.SyncInterface {
	hs.Spec.InstanceInfo.S3Config.RegionEndpoint = syncer.DefaultEndpoint
	minioClient, _ := operator_s3.NewMinioClient(hs.Spec.InstanceInfo.S3Config, opConfig)
	return &MinioSyncer{
		Name:          name,
		BucketName:    hs.Spec.InstanceInfo.S3Config.Bucket,
		HarborService: hs,
		Client:        c,
		MinioClient:   minioClient,
		EventCli:      eventCli,
	}
}

func (mi *MinioSyncer) GetName() interface{} {
	return mi.Name
}

func (mi *MinioSyncer) Sync(c context.Context) (ctrl.Result, error) {
	switch mi.HarborService.Status.Condition.Phase {
	case "":
		// 实现创建逻辑
		// 在minio中创建bucket
		installInfo := fmt.Sprintf("create bucket: %s", mi.BucketName)
		klog.Info(installInfo)
		mi.EventCli.NewEventAdd(mi.HarborService, "createS3Bucket", "Start create minio bucket for harborservice")

		bucketName := mi.BucketName
		err := mi.MinioClient.CreateBucket(bucketName)
		if err != nil {
			// 数据库初始化失败，更新cr status和event
			mi.EventCli.NewEventAdd(mi.HarborService, "createS3BucketFailed", "Create minio bucket for harborservice failed")
			mi.HarborService.Status.SetFailedStatus("Create minio bucket failed")
			v1.FlushInstanceStatus(mi.Client, mi.HarborService)
			return ctrl.Result{}, err
		}

		// 同database,创建成功无需更新cr status
		klog.Info("Create minio bucket successd.")
		mi.EventCli.NewEventAdd(mi.HarborService, "createS3BucketSuccess", "Sync pgsql database success")
	}

	return ctrl.Result{}, nil
}

func (mi *MinioSyncer) Delete(c context.Context) error {
	/*
		errObj := mi.MinioClient.RemoveObject(mi.BucketName, "docker")
		if errObj != nil {
			mi.Log.Error(errObj, "fail to remove object in bucket")
		}

		errBucket := mi.MinioClient.RemoveBucket(mi.BucketName)
		if errBucket != nil {
			mi.Log.Error(errObj, "fail to remove bucket in minio")
		}

		if errObj != nil || errBucket != nil {
			return errors.New("fail to clear minio resources")
		}

		mi.Log.Info("Remove minio resource success")
	*/

	// 对象存储资源保留，不删除，用户不需要的话可以到minio自行删除
	return nil
}
