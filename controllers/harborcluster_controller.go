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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	harborCluster "src/github.com/goharbor/harbor-cluster-operator/api/v1"
	"src/github.com/goharbor/harbor-cluster-operator/controllers/cache"

	goharborv1 "src/github.com/goharbor/harbor-cluster-operator/api/v1"
)

// HarborClusterReconciler reconciles a HarborCluster object
type HarborClusterReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	DClient  dynamic.Interface
}

// +kubebuilder:rbac:groups=goharbor.goharbor.io,resources=harborclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=goharbor.goharbor.io,resources=harborclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=databases.spotahome.com,resources=redisfailovers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list

func (r *HarborClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var err error
	ctx := context.Background()
	log := r.Log.WithValues("harborCluster", req.NamespacedName)

	var harborCluster harborCluster.HarborCluster
	if err := r.Get(ctx, req.NamespacedName, &harborCluster); err != nil {
		if errors.IsNotFound(err) {
			log.Info(fmt.Sprintf("HarborCluster %s has been deleted", req.Name))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	results := cache.NewDefaultCache(cache.ReconcileCache{
		Client:   r.Client,
		Recorder: r.Recorder,
		Log:      r.Log,
		CTX:      ctx,
		Request:  req,
		DClient:  r.DClient,
		Harbor:   &harborCluster,
		Scheme:   r.Scheme,
	}).Reconcile()

	return results.WithError(err).Aggregate()
}

func (r *HarborClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goharborv1.HarborCluster{}).
		Complete(r)
}
