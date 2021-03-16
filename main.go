/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"harbor-operator/config"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	harborv1 "harbor-operator/api/v1"
	"harbor-operator/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme     = runtime.NewScheme()
	setupLog   = ctrl.Log.WithName("setup")
	ConfigInfo *config.ConfigFile
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = harborv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme

	config, err := config.ParseConfigFile()
	if err != nil {
		klog.Error(fmt.Sprintf("fail to parse config file, witherr: %s", err.Error()))
		return
	}
	ConfigInfo = config
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "e2a5dbc2.example.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(fmt.Errorf("fail to get kubeclient"), "fail to get kubeclient")
		os.Exit(1)
	}

	if err = (&controllers.HarborServiceReconciler{
		KubeClient: kubeClient,
		Client:     mgr.GetClient(),
		Log:        ctrl.Log.WithName("controllers").WithName("HarborService"),
		Scheme:     mgr.GetScheme(),
		EventsCli:  event.NewEvent(mgr.GetEventRecorderFor("harbor-operator")),
		ConfigInfo: ConfigInfo,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HarborService")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
