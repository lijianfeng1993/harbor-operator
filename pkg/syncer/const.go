package syncer

const (
	OperationCreate = "create"

	GroupV1   string = "harbor.example.com"
	VersionV1 string = "v1"
	Resource  string = "harborservices"

	StatusRunning string = "running"

	DefaultMinioRegion string = "us-east-1"
	DefaultSSLMode     string = "disable"
	DefaultEndpoint    string = "10.142.113.234:9000"

	DefaultPgHost     string = "10.142.113.234"
	DefaultPgPort     string = "5432"
	DefaultPgUser     string = "postgres"
	DefaultPgPassword string = "8cDcos11"

	DefaultRedisAddr     string = "10.142.113.234:6379"
	DefaultRedisPassword string = "8cDcos11"

	DefaultApiserver           string = "https://10.142.113.231:6443"
	DefaultKubeToken           string = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkVCWDBWRkFSYjJkak04SlNzemN5bXIzVDFobmh0Tzc1cmZZUlJydDh5Y28ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi01a3JsYyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImU5MDcwZTRkLTdkNDItNDYwNC1hMjZmLWYzYTM1Njc0MTc1MSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.UlL5n3dGxPkKRFuH9Nm3wFn_u_7L-col8h7KXcKUGwrPeJ9-s7VYWy7uiQ53Jw7JO8aMzG_a75ymqfAX20X3Q9lrM8dWRW4Y1YMqtb-OeFqPLAFoI197LQFGNyVAXcaYLjyMhHvdrRKlj4LQ_oCXPlaFwkmo6nuxGYpLYmWDxjca0HyBdx1Wn_68YGUI2W7zSsxDqTLDkqc7SULCOYOOaAcGfi3QN_LwyEQa807JzL6FzMODnniSiaKkw_NQsJNxfASi2uPWiUxyQsWpy1N8e5l9vsHPVUfPGOLFLEURW_GNPKBAYouBj5LQuJrpOId6F0tjWwFlsz2dOrZF69OEyw"
	DefaultHarborHelmChartPath string = "/tmp/harbor-helm-1.5.3"
)
