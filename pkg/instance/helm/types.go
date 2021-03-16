package helm

import (
	helmtime "helm.sh/helm/v3/pkg/time"
	"time"
)

type KubeToken struct {
	ApiServer string
	Token     string
}

// helm List struct
type ReleaseListOptions struct {
	// All ignores the limit/offset
	All       bool   `json:"all"`
	Namespace string `json:"namespace"`
	// AllNamespaces searches across namespaces
	AllNamespaces bool `json:"all_namespaces"`
	// Overrides the default lexicographic sorting
	ByDate      bool `json:"by_date"`
	SortReverse bool `json:"sort_reverse"`
	// Limit is the number of items to return per Run()
	Limit int `json:"limit"`
	// Offset is the starting index for the Run() call
	Offset int `json:"offset"`
	// Filter is a filter that is applied to the results
	Filter       string `json:"filter"`
	Uninstalled  bool   `json:"uninstalled"`
	Superseded   bool   `json:"superseded"`
	Uninstalling bool   `json:"uninstalling"`
	Deployed     bool   `json:"deployed"`
	Failed       bool   `json:"failed"`
	Pending      bool   `json:"pending"`
}

// helm release options
type ReleaseOptions struct {
	// common
	DryRun          bool          `json:"dry_run"`
	DisableHooks    bool          `json:"disable_hooks"`
	Wait            bool          `json:"wait"`
	Devel           bool          `json:"devel"`
	Description     string        `json:"description"`
	Atomic          bool          `json:"atomic"`
	SkipCRDs        bool          `json:"skip_crds"`
	SubNotes        bool          `json:"sub_notes"`
	Timeout         time.Duration `json:"timeout"`
	Values          string        `json:"values"`
	SetValues       []string      `json:"set"`
	SetStringValues []string      `json:"set_string"`
	ChartPathOptions

	// only install
	CreateNamespace  bool `json:"create_namespace"`
	DependencyUpdate bool `json:"dependency_update"`

	// only upgrade
	Force         bool `json:"force"`
	Install       bool `json:"install"`
	Recreate      bool `json:"recreate"`
	CleanupOnFail bool `json:"cleanup_on_fail"`
}

// ChartPathOptions captures common options used for controlling chart paths
type ChartPathOptions struct {
	CaFile                string `json:"ca_file"`              // --ca-file
	CertFile              string `json:"cert_file"`            // --cert-file
	KeyFile               string `json:"key_file"`             // --key-file
	InsecureSkipTLSverify bool   `json:"insecure_skip_verify"` // --insecure-skip-verify
	Keyring               string `json:"keyring"`              // --keyring
	Password              string `json:"password"`             // --password
	RepoURL               string `json:"repo"`                 // --repo
	Username              string `json:"username"`             // --username
	Verify                bool   `json:"verify"`               // --verify
	Version               string `json:"version"`              // --version
}

type ReleaseElement struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Revision     string `json:"revision"`
	Updated      string `json:"updated"`
	Status       string `json:"status"`
	Chart        string `json:"chart"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`

	Notes string `json:"notes,omitempty"`

	// TODO: Test Suite?
}

type ReleaseInfo struct {
	Revision    int           `json:"revision"`
	Updated     helmtime.Time `json:"updated"`
	Status      string        `json:"status"`
	Chart       string        `json:"chart"`
	AppVersion  string        `json:"app_version"`
	Description string        `json:"description"`
}

type ReleaseHistory []ReleaseInfo
