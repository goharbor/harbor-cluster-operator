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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	goharborv1 "src/github.com/goharbor/harbor-cluster-operator/api/v1"
)

// HarborClusterReconciler reconciles a HarborCluster object
type HarborClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=goharbor.goharbor.io,resources=harborclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=goharbor.goharbor.io,resources=harborclusters/status,verbs=get;update;patch

func (r *HarborClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("harborcluster", req.NamespacedName)

	var harborCluster goharborv1.HarborCluster
	if err := r.Get(ctx, req.NamespacedName, &harborCluster); err != nil {
		log.Error(err, "unable to fetch HarborCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// harborCluster will be gracefully deleted by server when DeletionTimestamp is non-null
	if harborCluster.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	if r.RequiredRedis(&harborCluster) {
		if err := r.ReconcileRedis(ctx, req, &harborCluster); err != nil {
			log.Error(err, "error when reconcile redis components.")
			return ctrl.Result{}, err
		}
	}

	if r.RequiredDatabase(&harborCluster) {
		if err := r.ReconcileDatabase(ctx, req, &harborCluster); err != nil {
			log.Error(err, "error when reconcile database components.")
			return ctrl.Result{}, err
		}
	}

	if r.RequiredStorage(&harborCluster) {
		if err := r.ReconcileStorage(ctx, req, &harborCluster); err != nil {
			log.Error(err, "error when reconcile storage components.")
			return ctrl.Result{}, err
		}
	}

	// if components is not all ready, reenqueue the HarborCluster
	if !r.ComponentsIsAllReady(&harborCluster) {
		return ctrl.Result{
			Requeue:      true,
			// TODO: config requeue time when operator started.
			RequeueAfter: time.Second * 1,
		}, nil
	}

	if err := r.ReconcileHarborCore(ctx, req, &harborCluster); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// check whether required deploy redis.
func (r *HarborClusterReconciler) RequiredRedis(harborCluster *goharborv1.HarborCluster) bool {
	// TODO
	return false
}

// check whether required deploy storage(MinIO)
func (r *HarborClusterReconciler) RequiredStorage(harborCluster *goharborv1.HarborCluster) bool {
	// TODO
	return false
}

// check whether required deploy database(PostgreSQL)
func (r *HarborClusterReconciler) RequiredDatabase(harborCluster *goharborv1.HarborCluster) bool {
	// TODO
	return false
}

func (r *HarborClusterReconciler) ComponentsIsAllReady(harborCluster *goharborv1.HarborCluster) bool {
	// TODO
	return false
}

func (r *HarborClusterReconciler) ReconcileRedis(ctx context.Context, req ctrl.Request, harborCluster *goharborv1.HarborCluster) error {
	// TODO
	return nil
}

func (r *HarborClusterReconciler) ReconcileStorage(ctx context.Context, req ctrl.Request, harborCluster *goharborv1.HarborCluster) error {
	// TODO
	return nil
}

func (r *HarborClusterReconciler) ReconcileDatabase(ctx context.Context, req ctrl.Request, harborCluster *goharborv1.HarborCluster) error {
	// TODO
	return nil
}

func (r *HarborClusterReconciler) ReconcileHarborCore(ctx context.Context, req ctrl.Request, harborCluster *goharborv1.HarborCluster) error {
	// TODO
	return nil
}

func (r *HarborClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goharborv1.HarborCluster{}).
		Complete(r)
}
