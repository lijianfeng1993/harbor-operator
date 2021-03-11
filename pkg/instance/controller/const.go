package controller

const (
	ImagePullPolicy = "IfNotPresent"
	UpdateStrategy = "RollingUpdate"
	LogLevel = "debug"
	DefaultHarborAdminPassword = "8cDcos11"
	DefaultHarborInternalDbPassword = "8cDcos11"

	S3DefaultEncrypt = false
	S3DefaultSecure = false
	S3DefaultV4auth = true
	S3DefaultChunksize = "5242880"
	S3DefaultRootDirectory = "/"


	NginxRepoName = "registry.paas/cmss/harbor-portal"
	NginxRepoTag = "v2.1.3"
	PortalRepoName = "registry.paas/cmss/harbor-portal"
	PortalRepoTag = "v2.1.3"
	CoreRepoName = "registry.paas/cmss/harbor-core"
	CoreRepoTag = "v2.1.3"
	JobserviceRepoName = "registry.paas/cmss/harbor-jobservice"
	JobserviceRepoTag = "v2.1.3"
	RegistryRepoName = "registry.paas/cmss/registry-photon"
	RegistryRepoTag = "v2.1.3"
	RegistryCtlRepoName = "registry.paas/cmss/harbor-registryctl"
	RegistryCtlRepoTag = "v2.1.3"
	ChartmuseumRepoName = "registry.paas/cmss/chartmuseum-photon"
	ChartmuseumRepoTag = "v2.1.3"
	ClairRepoName = "registry.paas/cmss/clair-photon"
	ClairRepoTag = "v2.1.3"
	ClairAdapterRepoName = "registry.paas/cmss/clair-adapter-photon"
	ClairAdapterRepoTag = "v2.1.3"
	TrivyRepoName = "registry.paas/cmss/trivy-adapter-photon"
	TrivyRepoTag = "v2.1.3"
	NotaryServerRepoName = "registry.paas/cmss/notary-server-photon"
	NotaryServerRepoTag = "v2.1.3"
	NotarySignerRepoName = "registry.paas/cmss/notary-signer-photon"
	NotarySignerRepoTag = "v2.1.3"
	HarborDbRepoName = "registry.paas/cmss/harbor-db"
	HarborDbRepoTag = "v2.1.3"
	HarborRedisRepoName = "registry.paas/cmss/redis-photon"
	HarborRedisRepoTag = "v2.1.3"

)
