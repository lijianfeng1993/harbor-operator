package helm

//define install harbor request info
type InstallInfo struct {
	InstanceName    string    `json:"instanceName"`
	InstanceVersion string    `json:"instanceVersion"`
	NodePortIndex   int64     `json:"nodePortIndex"`
	S3              S3Config  `json:"s3"`
	PgInfo          PgInfo    `json:"pginfo"`
	RedisInfo       RedisInfo `json:"redisinfo"`
	Jobservicepvc   string    `json:"jobservicepvc"`
}

type RedisInfo struct {
	RedisDbIndex  int64 `json:"redisDbIndex"`
	RedisAddr     string
	RedisPassword string
}

type PgInfo struct {
	Host                 string
	Port                 string
	UserName             string
	Password             string
	CoreDatabase         string
	ClairDatabase        string
	NotaryServerDatabase string
	NotarySignerDatabase string
	SSLMode              string
}

// define values.yaml
type Values struct {
	Expose              Expose            `yaml:"expose"`
	ExternalUrl         string            `yaml:"externalURL"`
	InternalTls         InternalTls       `yaml:"internalTLS"`
	Persistence         PersistenceConfig `yaml:"persistence"`
	ImagePullPolicy     string            `yamlm:"imagePullPolicy"`
	UpdateStrategy      UpdateStrategy    `yaml:"updateStrategy"`
	LogLevel            string            `yaml:"logLevel"`
	HarborAdminPassword string            `yaml:"harborAdminPassword"`
	CaSecretName        string            `yaml:"caSecretName"`
	SecretKey           string            `yaml:"secretKey"`
	Proxy               Proxy             `yaml:"proxy"`
	Nginx               Nginx             `yaml:"nginx"`
	Portal              Portal            `yaml:"portal"`
	Core                Core              `yaml:"core"`
	Jobservice          Jobservice        `yaml:"jobservice"`
	Registry            Registry          `yaml:"registry"`
	Chartmuseum         Chartmuseum       `yaml:"chartmuseum"`
	Clair               Clair             `yaml:"clair"`
	Trivy               Trivy             `yaml:"trivy"`
	Notary              Notary            `yaml:"notary"`
	Database            Database          `yaml:"database"`
	Redis               Redis             `yaml:"redis"`
}

type Expose struct {
	Type     string         `yaml:"type"`
	Tls      TlsConfig      `yaml:"tls"`
	NodePort NodePortConfig `yaml:"nodePort"`
	Ingress  IngressConfig  `yaml:"ingress"`
}

type TlsConfig struct {
	Enable     bool            `yaml:"enabled"`
	CertSource string          `yaml:"certSource"`
	Auto       TlsConfigAuto   `yaml:"auto"`
	Secret     TlsConfigSecret `yaml:"secret"`
}

type TlsConfigAuto struct {
	CommonName string `yaml:"commonName"`
}

type TlsConfigSecret struct {
	SecretName       string `yaml:"secretName"`
	NotarySecretName string `yaml:"notarySecretName"`
}

type IngressConfig struct {
	Controller string `yaml:"controller"`
}

type NodePortConfig struct {
	Name  string      `yaml:"name"`
	Ports PortsConfig `yaml:"ports"`
}

type PortsConfig struct {
	Http   PortNodePort `yaml:"http"`
	Https  PortNodePort `yaml:"https"`
	Notary PortNodePort `yaml:"notary"`
}

type PortNodePort struct {
	Port     int64 `yaml:"port"`
	NodePort int64 `yaml:"nodePort"`
}

type InternalTls struct {
	Enable     bool   `yaml:"enabled"`
	CertSource string `yaml:"certSource"`
}

type PersistenceConfig struct {
	Enable                bool                        `yaml:"enabled"`
	ResourcePolicy        string                      `yaml:"resourcePolicy"`
	PersistentVolumeClaim PersistentVolumeClaimConfig `yaml:"persistentVolumeClaim"`
	ImageChartStorage     ImageChartStorage           `yaml:"imageChartStorage"`
}

type PersistentVolumeClaimConfig struct {
	RegistryPVC    RegistryPVC   `yaml:"registry"`
	JobservicePVC  JobservicePVC `yaml:"jobservice"`
	ChartmuseumPVC RegistryPVC   `yaml:"chartmuseum"`
}

type RegistryPVC struct {
	ExistingClaim string `yaml:"existingClaim"`
	StorageClass  string `yaml:"storageClass"`
	SubPath       string `yaml:"subPath"`
	AccessMode    string `yaml:"accessMode"`
	Size          string `yaml:"size"`
}

type JobservicePVC struct {
	ExistingClaim string `yaml:"existingClaim"`
	StorageClass  string `yaml:"storageClass"`
	SubPath       string `yaml:"subPath"`
	AccessMode    string `yaml:"accessMode"`
	Size          string `yaml:"size"`
}

type ImageChartStorage struct {
	Disableredirect bool       `yaml:"disableredirect"`
	Type            string     `yaml:"type"`
	FileSystem      FileSystem `yaml:"filesystem"`
	S3              S3Config   `yaml:"s3"`
}

type FileSystem struct {
	Rootdirectory string `yaml:"rootdirectory"`
}

type S3Config struct {
	Region         string `yaml:"region"`
	Bucket         string `yaml:"bucket"`
	Accesskey      string `yaml:"accesskey"`
	SecretKey      string `yaml:"secretkey"`
	RegionEndpoint string `yaml:"regionendpoint"`
	Encrypt        bool   `yaml:"encrypt"`
	Secure         bool   `yaml:"secure"`
	V4auth         bool   `yaml:"v4auth"`
	Chunksize      string `yaml:"chunksize"`
	RootDirectory  string `yaml:"rootdirectory"`
	/*
		Storageclass string `yaml:"storageclass"`
		Multipartcopychunksize string 	`yaml:"multipartcopychunksize"`
		Multipartcopymaxconcurrency string `yaml:"multipartcopymaxconcurrency"`
		Multipartcopythresholdsize string `yaml:"multipartcopythresholdsize"`
	*/
}

type UpdateStrategy struct {
	Type string `yaml:"type"`
}

type Proxy struct {
	HttpProxy  string   `yaml:"httpProxy"`
	HttpsProxy string   `yaml:"httpsProxy"`
	NoProxy    string   `yaml:"noProxy"`
	Componemts []string `yaml:"components"`
}

type Nginx struct {
	Image              Image  `yaml:"image"`
	ServiceAccountName string `yaml:"serviceAccountName"`
	Replicase          int64  `yaml:"replicas"`
	//Resources          Resources `yaml:"resources"`
}

type Portal struct {
	Image              Image  `yaml:"image"`
	ServiceAccountName string `yaml:"serviceAccountName"`
	Replicase          int64  `yaml:"replicas"`
	//Resources          Resources `yaml:"resources"`
}

type Core struct {
	Image              Image  `yaml:"image"`
	ServiceAccountName string `yaml:"serviceAccountName"`
	Replicase          int64  `yaml:"replicas"`
	//Resources          Resources `yaml:"resources"`
	RepositoryQuota int64        `yaml:"repositoryquota"`
	HelmChartQuota  int64        `yaml:"helmchartquota"`
	StartupProbe    StartupProbe `yaml:"startupProbe"`
	Secret          string       `yaml:"secret"`
	SecretName      string       `yaml:"secretName"`
	XsrfKey         string       `yaml:"xsrfKey"`
}

type Jobservice struct {
	Image              Image  `yaml:"image"`
	ServiceAccountName string `yaml:"serviceAccountName"`
	Replicase          int64  `yaml:"replicas"`
	MaxJobWorkers      int64  `yaml:"maxJobWorkers"`
	JobLogger          string `yaml:"jobLogger"`
	//Resources          Resources `yaml:"resources"`
	Secret string `yaml:"secret"`
}

type StartupProbe struct {
	Enable              bool  `yaml:"enabled"`
	InitialDelaySeconds int64 `yaml:"initialDelaySeconds"`
}

type Registry struct {
	ServiceAccountName string              `yaml:"serviceAccountName"`
	Registry           RegistryContainer   `yaml:"registry"`
	Controller         RegistryContainer   `yaml:"controller"`
	Replicase          int64               `yaml:"replicas"`
	Secret             string              `yaml:"secret"`
	Relativeurls       bool                `yaml:"relativeurls"`
	Credentials        RegistryCredentials `yaml:"credentials"`
	Middleware         RegistryMiddleware  `yaml:"middleware"`
}

type RegistryContainer struct {
	Image Image `yaml:"image"`
	//Resources Resources `yaml:"resources"`
}

type RegistryCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Htpasswd string `yaml:"htpasswd"`
}

type RegistryMiddleware struct {
	Enable     bool       `yaml:"enabled"`
	Type       string     `yaml:"type"`
	CloudFront CloudFront `yaml:"middleware"`
}

type CloudFront struct {
	Baseurl          string `yaml:"baseurl"`
	Keypairid        string `yaml:"keypairid"`
	Duration         string `yaml:"duration"`
	Ipfilteredby     string `yaml:"ipfilteredby"`
	PrivateKeySecret string `yaml:"privateKeySecret"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type Resources struct {
	Requests struct {
		Memory string `yaml:"memory"`
		CPU    string `yaml:"cpu"`
	} `yaml:"requests"`
	Limits struct {
		Memory string `yaml:"memory"`
		CPU    string `yaml:"cpu"`
	} `yaml:"limits"`
}

type Chartmuseum struct {
	Enable             bool   `yaml:"enabled"`
	ServiceAccountName string `yaml:"serviceAccountName"`
	AbsoluteUrl        bool   `yaml:"absoluteUrl"`
	Image              Image  `yaml:"image"`
	Replicase          int64  `yaml:"replicas"`
}

type Clair struct {
	Enable             bool           `yaml:"enabled"`
	ServiceAccountName string         `yaml:"serviceAccountName"`
	Clair              ClairContainer `yaml:"clair"`
	Adapter            ClairContainer `yaml:"adapter"`
	Replicase          int64          `yaml:"replicas"`
	UpdatersInterval   int64          `yaml:"updatersInterval"`
}

type ClairContainer struct {
	Image Image `yaml:"image"`
	//Resources Resources `yaml:"resources"`
}

type Trivy struct {
	Enable             bool   `yaml:"enabled"`
	Image              Image  `yaml:"image"`
	ServiceAccountName string `yaml:"serviceAccountName"`
	Replicase          int64  `yaml:"replicas"`
	DebugMode          bool   `yaml:"debugMode"`
	VulnType           string `yaml:"vulnType"`
	Severity           string `yaml:"severity"`
	IgnoreUnfixed      bool   `yaml:"ignoreUnfixed"`
	Insecure           bool   `yaml:"insecure"`
	GitHubToken        string `yaml:"gitHubToken"`
	SkipUpdate         bool   `yaml:"skipUpdate"`
	//Resources          Resources `yaml:"resources"`
}

type Notary struct {
	Enable     bool            `yaml:"enabled"`
	Server     NotaryContainer `yaml:"server"`
	Signer     NotaryContainer `yaml:"signer"`
	SecretName string          `yaml:"secretName"`
}

type NotaryContainer struct {
	ServiceAccountName string `yaml:"serviceAccountName"`
	Image              Image  `yaml:"image"`
	//Resources          Resources `yaml:"resources"`
	Replicase int64 `yaml:"replicas"`
}

type Database struct {
	Type           string           `yaml:"type"`
	Internal       DatabaseInternal `yaml:"internal"`
	External       DatabaseExternal `yaml:"external"`
	MaxIdleConns   int64            `yaml:"maxIdleConns"`
	MaxOpenConns   int64            `yaml:"maxOpenConns"`
	PodAnnotations int64            `yaml:"podAnnotations"`
}

type DatabaseInternal struct {
	ServiceAccountName string `yaml:"serviceAccountName"`
	Image              Image  `yaml:"image"`
	Password           string `yaml:"password"`
	//Resources          Resources `yaml:"resources"`
}

type DatabaseExternal struct {
	Host                 string `yaml:"host"`
	Port                 string `yaml:"port"`
	Username             string `yaml:"username"`
	Password             string `yaml:"password"`
	CoreDatabase         string `yaml:"coreDatabase"`
	ClairDatabase        string `yaml:"clairDatabase"`
	NotaryServerDatabase string `yaml:"notaryServerDatabase"`
	NotarySignerDatabase string `yaml:"notarySignerDatabase"`
	Sslmode              string `yaml:"sslmode"`
}

type Redis struct {
	Type     string        `yaml:"type"`
	Internal RedisInternal `yaml:"internal"`
	External RedisExternal `yaml:"external"`
}

type RedisInternal struct {
	ServiceAccountName string `yaml:"serviceAccountName"`
	Image              Image  `yaml:"image"`
	//Resources          Resources `yaml:"resources"`
}

type RedisExternal struct {
	Addr                     string `yaml:"addr"`
	SentinelMasterSet        string `yaml:"sentinelMasterSet"`
	CoreDatabaseIndex        string `yaml:"coreDatabaseIndex"`
	JobserviceDatabaseIndex  string `yaml:"jobserviceDatabaseIndex"`
	RegistryDatabaseIndex    string `yaml:"registryDatabaseIndex"`
	ChartmuseumDatabaseIndex string `yaml:"chartmuseumDatabaseIndex"`
	ClairAdapterIndex        string `yaml:"clairAdapterIndex"`
	TrivyAdapterIndex        string `yaml:"trivyAdapterIndex"`
	Password                 string `yaml:"password"`
}
