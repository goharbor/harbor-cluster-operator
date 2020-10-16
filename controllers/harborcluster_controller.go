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
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/goharbor/harbor-cluster-operator/controllers/image"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	goharborv1 "github.com/goharbor/harbor-cluster-operator/apis/goharbor.io/v1alpha2"
)

// HarborClusterReconciler reconciles a HarborCluster object
type HarborClusterReconciler struct {
	client.Client
	ServiceGetter
	Log          logr.Logger
	Scheme       *runtime.Scheme
	RequeueAfter time.Duration
	Recorder     record.EventRecorder
}

var (
	ComponentToConditionType = map[goharborv1.Component]goharborv1.HarborClusterConditionType{
		goharborv1.ComponentHarbor:   goharborv1.ServiceReady,
		goharborv1.ComponentCache:    goharborv1.CacheReady,
		goharborv1.ComponentStorage:  goharborv1.StorageReady,
		goharborv1.ComponentDatabase: goharborv1.DatabaseReady,
	}
	ReconcileWaitResult = reconcile.Result{RequeueAfter: 30 * time.Second}
)

// +kubebuilder:rbac:groups=goharbor.io,resources=harborclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=goharbor.io,resources=harborclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=goharbor.io,resources=harbors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=databases.spotahome.com,resources=redisfailovers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=acid.zalan.do,resources=postgresqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list
// +kubebuilder:rbac:groups=minio.min.io,resources=tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets;deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods;configmaps;services;events;secrets;ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;update

func (r *HarborClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("harborcluster", req.NamespacedName)

	log.Info("start to reconcile.")

	var harborCluster goharborv1.HarborCluster
	if err := r.Get(ctx, req.NamespacedName, &harborCluster); err != nil {
		log.Error(err, "unable to fetch HarborCluster")
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}

	// harborCluster will be gracefully deleted by server when DeletionTimestamp is non-null
	if harborCluster.DeletionTimestamp != nil {
		return ReconcileWaitResult, nil
	}

	dClient, err := k8s.NewDynamicClient()
	if err != nil {
		log.Error(err, "unable to create dynamic client")
		return ReconcileWaitResult, err
	}

	option := &GetOptions{
		Client:   k8s.WrapClient(ctx, r.Client),
		Recorder: r.Recorder,
		Log:      r.Log,
		DClient:  k8s.WrapDClient(dClient),
		Scheme:   r.Scheme,
	}

	componentToStatus := r.DefaultComponentStatus()
	cacheStatus, err := r.Cache(ctx, &harborCluster, option).Reconcile()
	componentToStatus[goharborv1.ComponentCache] = cacheStatus
	if err != nil {
		log.Error(err, "error when reconcile cache component.")
		updateErr := r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
		if updateErr != nil {
			log.Error(updateErr, "update harbor cluster status")
		}
		return ReconcileWaitResult, err
	}

	dbStatus, err := r.Database(ctx, &harborCluster, option).Reconcile()
	componentToStatus[goharborv1.ComponentDatabase] = dbStatus
	if err != nil {
		log.Error(err, "error when reconcile database component.")
		updateErr := r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
		if updateErr != nil {
			log.Error(updateErr, "update harbor cluster status")
		}
		return ReconcileWaitResult, err
	}

	storageStatus, err := r.Storage(ctx, &harborCluster, option).Reconcile()
	componentToStatus[goharborv1.ComponentStorage] = storageStatus
	if err != nil {
		log.Error(err, "error when reconcile storage component.")
		updateErr := r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
		if updateErr != nil {
			log.Error(updateErr, "update harbor cluster status")
		}
		return ReconcileWaitResult, err
	}

	// if components is not all ready, requeue the HarborCluster
	if !r.ComponentsAreAllReady(componentToStatus) {
		log.Info("components not all ready.",
			string(goharborv1.ComponentCache), cacheStatus,
			string(goharborv1.ComponentDatabase), dbStatus,
			string(goharborv1.ComponentStorage), storageStatus)
		err = r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
		return ReconcileWaitResult, err
	}

	getRegistry := func() *string {
		if harborCluster.Spec.ImageSource != nil && harborCluster.Spec.ImageSource.Registry != "" {
			return &harborCluster.Spec.ImageSource.Registry
		}
		return nil
	}
	var imageGetter image.Getter
	if imageGetter, err = image.NewImageGetter(getRegistry(), harborCluster.Spec.Version); err != nil {
		log.Error(err, "error when create Getter.")
		return ReconcileWaitResult, err
	}
	option.ImageGetter = imageGetter
	harborStatus, err := r.Harbor(ctx, &harborCluster, componentToStatus, option).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile harbor service.")
		return ReconcileWaitResult, err
	}
	componentToStatus[goharborv1.ComponentHarbor] = harborStatus

	err = r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
	if err != nil {
		log.Error(err, "error when update harbor cluster status.")
		return ReconcileWaitResult, err
	}
	// wait to resync to update status.
	return ReconcileWaitResult, nil
}

func (r *HarborClusterReconciler) DefaultComponentStatus() map[goharborv1.Component]*lcm.CRStatus {
	return map[goharborv1.Component]*lcm.CRStatus{
		goharborv1.ComponentCache:    lcm.New(goharborv1.CacheReady).WithStatus(corev1.ConditionUnknown),
		goharborv1.ComponentDatabase: lcm.New(goharborv1.DatabaseReady).WithStatus(corev1.ConditionUnknown),
		goharborv1.ComponentStorage:  lcm.New(goharborv1.CacheReady).WithStatus(corev1.ConditionUnknown),
		goharborv1.ComponentHarbor:   lcm.New(goharborv1.ServiceReady).WithStatus(corev1.ConditionUnknown),
	}
}

// ServicesAreAllReady check whether these components(includes cache, db, storage) are all ready.
func (r *HarborClusterReconciler) ComponentsAreAllReady(serviceToMap map[goharborv1.Component]*lcm.CRStatus) bool {
	for _, status := range serviceToMap {
		if status == nil {
			return false
		}

		if status.Condition.Type == goharborv1.ServiceReady {
			continue
		}
		if status.Condition.Status != corev1.ConditionTrue {
			return false
		}
	}
	return true
}

// UpdateHarborClusterStatus will Update HarborCluster CR status, according the services reconcile result.
func (r *HarborClusterReconciler) UpdateHarborClusterStatus(
	ctx context.Context,
	harborCluster *goharborv1.HarborCluster,
	componentToCRStatus map[goharborv1.Component]*lcm.CRStatus) error {
	for component, status := range componentToCRStatus {
		if status == nil {
			continue
		}
		var conditionType goharborv1.HarborClusterConditionType
		var ok bool
		if conditionType, ok = ComponentToConditionType[component]; !ok {
			r.Log.Info(fmt.Sprintf("can not found the condition type for %s", component))
		}
		harborClusterCondition, defaulted := r.getHarborClusterCondition(harborCluster, conditionType)
		r.updateHarborClusterCondition(harborClusterCondition, status)
		if defaulted {
			harborCluster.Status.Conditions = append(harborCluster.Status.Conditions, *harborClusterCondition)
		}
	}
	r.Log.Info("update harbor cluster.", "harborcluster", harborCluster)
	return r.Update(ctx, harborCluster)
}

// updateHarborClusterCondition update condition according to status.
func (r *HarborClusterReconciler) updateHarborClusterCondition(condition *goharborv1.HarborClusterCondition, crStatus *lcm.CRStatus) {
	if condition.Type != crStatus.Condition.Type {
		return
	}

	if condition.Status != crStatus.Condition.Status ||
		condition.Message != crStatus.Condition.Message ||
		condition.Reason != crStatus.Condition.Reason {
		condition.Status = crStatus.Condition.Status
		condition.Message = crStatus.Condition.Message
		condition.Reason = crStatus.Condition.Reason
		condition.LastTransitionTime = metav1.Now()
	}
}

// getHarborClusterCondition will get HarborClusterCondition by conditionType
func (r *HarborClusterReconciler) getHarborClusterCondition(
	harborCluster *goharborv1.HarborCluster,
	conditionType goharborv1.HarborClusterConditionType) (condition *goharborv1.HarborClusterCondition, defaulted bool) {
	for i := range harborCluster.Status.Conditions {
		condition = &harborCluster.Status.Conditions[i]
		if condition.Type == conditionType {
			return condition, false
		}
	}
	return &goharborv1.HarborClusterCondition{
		Type:               conditionType,
		LastTransitionTime: metav1.Now(),
		Status:             corev1.ConditionUnknown,
	}, true
}

func (r *HarborClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goharborv1.HarborCluster{}).
		Complete(r)
}
