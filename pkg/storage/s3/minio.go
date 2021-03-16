package s3

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	"k8s.io/klog"
)

type MinioClient struct {
	MinioClient    *minio.Client
	OperatorConfig *config.ConfigFile
}

func NewMinioClient(s3 v1.S3Config, opConfig *config.ConfigFile) (*MinioClient, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(opConfig.Minio.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.Accesskey, s3.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create minio client, %s", err.Error()))
		return nil, err
	}
	return &MinioClient{
		MinioClient:    minioClient,
		OperatorConfig: opConfig,
	}, nil
}

func (mc *MinioClient) CreateBucket(bucketName string) error {
	// 检查bucket是否存在
	found, err := mc.MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		errorMessage := fmt.Sprintf("fail to list buckets, witherr: %s", err.Error())
		klog.Error(fmt.Sprintf(errorMessage))
		return err
	}

	if found {
		message := fmt.Sprintf("bucket: %s already exist, no need to create", bucketName)
		klog.Info(message)
		return nil
	}

	// 调用CreateBucket创建一个新的存储桶。
	err = mc.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: mc.OperatorConfig.Minio.Region, ObjectLocking: true})
	if err != nil {
		// 错误信息
		errorMessage := fmt.Sprintf("fail to create bucket: %s, witherr: %s", bucketName, err.Error())
		klog.Error(fmt.Sprintf(errorMessage))
		return err
	}
	return nil
}

func (mc *MinioClient) RemoveObject(bucketName, objectName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		//VersionID: "myversionid",
	}
	err := mc.MinioClient.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to remove object, witherr:%s", err.Error()))
		return err
	}
	return nil
}

func (mc *MinioClient) RemoveBucket(bucketName string) error {
	err := mc.MinioClient.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to remove bucket, witherr:%s", err.Error()))
		return err
	}
	return nil
}
