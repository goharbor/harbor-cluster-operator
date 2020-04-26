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
	"github.com/goharbor/harbor-cluster-operator/controllers/cache"
	"github.com/goharbor/harbor-cluster-operator/controllers/database"
	"github.com/goharbor/harbor-cluster-operator/controllers/harbor"
	"github.com/goharbor/harbor-cluster-operator/controllers/storage"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
)

// HarborClusterReconciler reconciles a HarborCluster object
type HarborClusterReconciler struct {
	client.Client
	ServiceGetter
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=goharbor.goharbor.io,resources=harborclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=goharbor.goharbor.io,resources=harborclusters/status,verbs=get;update;patch

func (r *HarborClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("harborcluster", req.NamespacedName)

	var harborCluster *goharborv1.HarborCluster
	if err := r.Get(ctx, req.NamespacedName, harborCluster); err != nil {
		log.Error(err, "unable to fetch HarborCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if harborCluster == nil {
		log.Info("can not get HarborCluster")
		return ctrl.Result{}, nil
	}

	// harborCluster will be gracefully deleted by server when DeletionTimestamp is non-null
	if harborCluster.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	cacheStatus, err := r.Cache(harborCluster).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile cache service.")
		return ctrl.Result{}, err
	}

	dbStatus, err := r.Database(harborCluster).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile database service.")
		return ctrl.Result{}, err
	}

	storageStatus, err := r.Storage(harborCluster).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile storage service.")
		return ctrl.Result{}, err
	}

	// if components is not all ready, reenqueue the HarborCluster
	if !r.ServicesAreAllReady(cacheStatus, dbStatus, storageStatus) {
		return ctrl.Result{
			Requeue: true,
			// TODO: config requeue time when operator started.
			RequeueAfter: time.Second * 1,
		}, nil
	}

	harborStatus, err := r.Harbor(harborCluster).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile harbor service.")
		return ctrl.Result{}, err
	}

	err = r.UpdateHarborClusterStatus(ctx, harborCluster, cacheStatus, dbStatus, storageStatus, harborStatus)
	if err != nil {
		log.Error(err, "error when update harbor cluster status.")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// check whether these services(includes cache, db, storage) are all ready.
func (r *HarborClusterReconciler) ServicesAreAllReady(statuses ...*lcm.CRStatus) bool {
	// TODO
	return false
}

// Update HarborCluster CR status, according the services reconcile result.
func (r *HarborClusterReconciler) UpdateHarborClusterStatus(ctx context.Context, harborCluster *goharborv1.HarborCluster, statuses ...*lcm.CRStatus) error {
	// TODO
	return nil
}

func (r *HarborClusterReconciler) Cache(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &cache.RedisReconciler{
		HarborCluster: harborCluster,
	}
}

func (r *HarborClusterReconciler) Database(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &database.PostgreSQLReconciler{
		HarborCluster: harborCluster,
	}
}

func (r *HarborClusterReconciler) Storage(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &storage.MinIOReconciler{
		HarborCluster: harborCluster,
	}
}

func (r *HarborClusterReconciler) Harbor(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &harbor.HarborReconciler{
		HarborCluster: harborCluster,
	}
}

func (r *HarborClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goharborv1.HarborCluster{}).
		Complete(r)
}
