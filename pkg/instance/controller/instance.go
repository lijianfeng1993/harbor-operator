package controller

import (
	"fmt"
	"harbor-operator/pkg/instance/helm"
	"k8s.io/klog"

	"gopkg.in/yaml.v2"
	"io/ioutil"

	"strconv"
)

func GenerageInstanceHelmChart(installInfo *helm.InstallInfo) (values helm.Values) {
	intanceValues := helm.Values{
		Expose: helm.Expose{
			Type: "nodePort",
			Tls: helm.TlsConfig{
				Enable:     false,
				CertSource: "secret",
				Auto:       helm.TlsConfigAuto{CommonName: ""},
				Secret:     helm.TlsConfigSecret{SecretName: "harbor-tls", NotarySecretName: ""},
			},
			NodePort: helm.NodePortConfig{
				Name: installInfo.InstanceName,
				Ports: helm.PortsConfig{
					Http: helm.PortNodePort{
						Port:     80,
						NodePort: installInfo.NodePortIndex,
					},
					Https: helm.PortNodePort{
						Port:     443,
						NodePort: installInfo.NodePortIndex + 1,
					},
					Notary: helm.PortNodePort{
						Port:     4443,
						NodePort: installInfo.NodePortIndex + 2,
					},
				},
			},
			Ingress: helm.IngressConfig{
				Controller: "default",
			},
		},

		ExternalUrl: fmt.Sprintf("http://%s.harbor.com:%d", installInfo.InstanceName, installInfo.NodePortIndex),
		InternalTls: helm.InternalTls{Enable: false, CertSource: "auto"},
		Persistence: helm.PersistenceConfig{
			Enable:         true,
			ResourcePolicy: "keep",
			PersistentVolumeClaim: helm.PersistentVolumeClaimConfig{
				JobservicePVC: helm.JobservicePVC{
					ExistingClaim: installInfo.Jobservicepvc,
					StorageClass:  "",
					SubPath:       "",
					AccessMode:    "ReadWriteOnce",
					Size:          "1Gi",
				},
				RegistryPVC: helm.RegistryPVC{
					SubPath: "",
				},
				ChartmuseumPVC: helm.RegistryPVC{
					SubPath: "",
				},
			},
			ImageChartStorage: helm.ImageChartStorage{
				Disableredirect: false,
				Type:            "s3",
				S3: helm.S3Config{
					Region:         installInfo.S3.Region,
					Bucket:         installInfo.S3.Bucket,
					Accesskey:      installInfo.S3.Accesskey,
					SecretKey:      installInfo.S3.SecretKey,
					RegionEndpoint: installInfo.S3.RegionEndpoint,
					Encrypt:        S3DefaultEncrypt,
					Secure:         S3DefaultSecure,
					V4auth:         S3DefaultV4auth,
					Chunksize:      S3DefaultChunksize,
					RootDirectory:  S3DefaultRootDirectory,
				},
				FileSystem: helm.FileSystem{
					Rootdirectory: "/storage",
				},
			},
		},
		ImagePullPolicy:     ImagePullPolicy,
		UpdateStrategy:      helm.UpdateStrategy{Type: UpdateStrategy},
		LogLevel:            LogLevel,
		HarborAdminPassword: DefaultHarborAdminPassword,
		CaSecretName:        "",
		SecretKey:           "not-a-secure-key",
		Proxy: helm.Proxy{
			NoProxy:    "127.0.0.1,localhost,.local,.internal",
			Componemts: []string{"core", "jobservice", "clair", "trivy"},
		},
		Nginx: helm.Nginx{
			Image: helm.Image{
				Repository: NginxRepoName,
				Tag:        NginxRepoTag,
			},
			ServiceAccountName: "",
			Replicase:          1,
		},
		Portal: helm.Portal{
			Image: helm.Image{
				Repository: PortalRepoName,
				Tag:        PortalRepoTag,
			},
			ServiceAccountName: "",
			Replicase:          1,
		},
		Core: helm.Core{
			Image: helm.Image{
				Repository: CoreRepoName,
				Tag:        CoreRepoTag,
			},
			ServiceAccountName: "",
			Replicase:          1,
			StartupProbe: helm.StartupProbe{
				Enable:              true,
				InitialDelaySeconds: 10,
			},
			SecretName: "",
			Secret:     "",
			XsrfKey:    "",
		},
		Jobservice: helm.Jobservice{
			Image: helm.Image{
				Repository: JobserviceRepoName,
				Tag:        JobserviceRepoTag,
			},
			Replicase:          1,
			ServiceAccountName: "",
			MaxJobWorkers:      10,
			JobLogger:          "file",
			Secret:             "",
		},
		Registry: helm.Registry{
			ServiceAccountName: "",
			Registry: helm.RegistryContainer{
				Image: helm.Image{
					Repository: RegistryRepoName,
					Tag:        RegistryRepoTag,
				},
			},
			Controller: helm.RegistryContainer{
				Image: helm.Image{
					Repository: RegistryCtlRepoName,
					Tag:        RegistryCtlRepoTag,
				},
			},
			Replicase:    1,
			Secret:       "",
			Relativeurls: false,
			Credentials: helm.RegistryCredentials{
				Username: "harbor_registry_user",
				Password: "harbor_registry_password",
				Htpasswd: "harbor_registry_user:$2y$10$9L4Tc0DJbFFMB6RdSCunrOpTHdwhid4ktBJmLD00bYgqkkGOvll3m",
			},
			Middleware: helm.RegistryMiddleware{
				Enable: false,
				Type:   "cloudFront",
				CloudFront: helm.CloudFront{
					Baseurl:          "example.cloudfront.net",
					Keypairid:        "KEYPAIRID",
					Duration:         "3000s",
					Ipfilteredby:     "none",
					PrivateKeySecret: "my-secret",
				},
			},
		},
		Chartmuseum: helm.Chartmuseum{
			Enable:             true,
			ServiceAccountName: "",
			AbsoluteUrl:        false,
			Image: helm.Image{
				Repository: ChartmuseumRepoName,
				Tag:        ChartmuseumRepoTag,
			},
			Replicase: 1,
		},
		Clair: helm.Clair{
			Enable:             true,
			ServiceAccountName: "",
			Clair: helm.ClairContainer{
				Image: helm.Image{
					Repository: ClairRepoName,
					Tag:        ClairRepoTag,
				},
			},
			Adapter: helm.ClairContainer{
				Image: helm.Image{
					Repository: ClairAdapterRepoName,
					Tag:        ClairAdapterRepoTag,
				},
			},
			Replicase:        1,
			UpdatersInterval: 12,
		},
		Trivy: helm.Trivy{
			Enable: false,
			Image: helm.Image{
				Repository: TrivyRepoName,
				Tag:        TrivyRepoTag,
			},
			ServiceAccountName: "",
			Replicase:          1,
			DebugMode:          false,
			VulnType:           "os,library",
			Severity:           "UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL",
			IgnoreUnfixed:      false,
			Insecure:           false,
			GitHubToken:        "",
			SkipUpdate:         false,
		},
		Notary: helm.Notary{
			Enable: true,
			Server: helm.NotaryContainer{
				ServiceAccountName: "",
				Image: helm.Image{
					Repository: NotaryServerRepoName,
					Tag:        NotaryServerRepoTag,
				},
				Replicase: 1,
			},
			Signer: helm.NotaryContainer{
				ServiceAccountName: "",
				Image: helm.Image{
					Repository: NotarySignerRepoName,
					Tag:        NotarySignerRepoTag,
				},
			},
			SecretName: "",
		},
		Database: helm.Database{
			Type: "external",
			Internal: helm.DatabaseInternal{
				ServiceAccountName: "",
				Image: helm.Image{
					Repository: HarborDbRepoName,
					Tag:        HarborDbRepoTag,
				},
				Password: DefaultHarborInternalDbPassword,
			},
			External: helm.DatabaseExternal{
				Host:                 installInfo.PgInfo.Host,
				Port:                 installInfo.PgInfo.Port,
				Username:             installInfo.PgInfo.UserName,
				Password:             installInfo.PgInfo.Password,
				CoreDatabase:         installInfo.PgInfo.CoreDatabase,
				ClairDatabase:        installInfo.PgInfo.ClairDatabase,
				NotaryServerDatabase: installInfo.PgInfo.NotaryServerDatabase,
				NotarySignerDatabase: installInfo.PgInfo.NotarySignerDatabase,
				Sslmode:              installInfo.PgInfo.SSLMode,
			},
			MaxIdleConns: 50,
			MaxOpenConns: 1000,
		},
		Redis: helm.Redis{
			Type: "external",
			Internal: helm.RedisInternal{
				ServiceAccountName: "",
				Image: helm.Image{
					Repository: HarborRedisRepoName,
					Tag:        HarborRedisRepoTag,
				},
			},
			External: helm.RedisExternal{
				Addr:                     installInfo.RedisInfo.RedisAddr,
				SentinelMasterSet:        "",
				CoreDatabaseIndex:        strconv.FormatInt(installInfo.RedisInfo.RedisDbIndex, 10),
				JobserviceDatabaseIndex:  strconv.FormatInt(installInfo.RedisInfo.RedisDbIndex+1, 10),
				RegistryDatabaseIndex:    strconv.FormatInt(installInfo.RedisInfo.RedisDbIndex+2, 10),
				ChartmuseumDatabaseIndex: strconv.FormatInt(installInfo.RedisInfo.RedisDbIndex+3, 10),
				ClairAdapterIndex:        strconv.FormatInt(installInfo.RedisInfo.RedisDbIndex+4, 10),
				TrivyAdapterIndex:        strconv.FormatInt(installInfo.RedisInfo.RedisDbIndex+5, 10),
				Password:                 installInfo.RedisInfo.RedisPassword,
			},
		},
	}
	return intanceValues
}

func SaveValuesToFile(values helm.Values, targetFilePath string) error {
	d, err := yaml.Marshal(values)
	if err != nil {
		klog.Errorf("fail to marshal values struct, witherr: %s", err.Error())
		return err
	}

	err = ioutil.WriteFile(targetFilePath, d, 0644)
	if err != nil {
		klog.Errorf("fail to write values to yaml file, witherr: %s", err.Error())
		return err
	}
	return nil
}
