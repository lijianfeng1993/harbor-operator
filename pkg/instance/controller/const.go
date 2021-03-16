package controller

const (
	ImagePullPolicy                 = "IfNotPresent"
	UpdateStrategy                  = "RollingUpdate"
	LogLevel                        = "debug"
	DefaultHarborAdminPassword      = "8cDcos11"
	DefaultHarborInternalDbPassword = "8cDcos11"

	S3DefaultEncrypt       = false
	S3DefaultSecure        = false
	S3DefaultV4auth        = true
	S3DefaultChunksize     = "5242880"
	S3DefaultRootDirectory = "/"

	NginxRepoName        = "goharbor/harbor-portal"
	NginxRepoTag         = "v2.1.3"
	PortalRepoName       = "goharbor/harbor-portal"
	PortalRepoTag        = "v2.1.3"
	CoreRepoName         = "goharbor/harbor-core"
	CoreRepoTag          = "v2.1.3"
	JobserviceRepoName   = "goharbor/harbor-jobservice"
	JobserviceRepoTag    = "v2.1.3"
	RegistryRepoName     = "goharbor/registry-photon"
	RegistryRepoTag      = "v2.1.3"
	RegistryCtlRepoName  = "goharbor/harbor-registryctl"
	RegistryCtlRepoTag   = "v2.1.3"
	ChartmuseumRepoName  = "goharbor/chartmuseum-photon"
	ChartmuseumRepoTag   = "v2.1.3"
	ClairRepoName        = "goharbor/clair-photon"
	ClairRepoTag         = "v2.1.3"
	ClairAdapterRepoName = "goharbor/clair-adapter-photon"
	ClairAdapterRepoTag  = "v2.1.3"
	TrivyRepoName        = "goharbor/trivy-adapter-photon"
	TrivyRepoTag         = "v2.1.3"
	NotaryServerRepoName = "goharbor/notary-server-photon"
	NotaryServerRepoTag  = "v2.1.3"
	NotarySignerRepoName = "goharbor/notary-signer-photon"
	NotarySignerRepoTag  = "v2.1.3"
	HarborDbRepoName     = "goharbor/harbor-db"
	HarborDbRepoTag      = "v2.1.3"
	HarborRedisRepoName  = "goharbor/redis-photon"
	HarborRedisRepoTag   = "v2.1.3"
)
