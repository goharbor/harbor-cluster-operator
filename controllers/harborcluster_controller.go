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
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/controllers/image"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
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
	Log          logr.Logger
	Scheme       *runtime.Scheme
	RequeueAfter time.Duration
	Recorder     record.EventRecorder
}

// +kubebuilder:rbac:groups=cluster.goharbor.io,resources=harborclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.goharbor.io,resources=harborclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=databases.spotahome.com,resources=redisfailovers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list

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
	
	dClient, err := k8s.NewDynamicClient()
	if err != nil {
		log.Error(err, "unable to create dynamic client")
		return ctrl.Result{}, err
	}

	option := &GetOptions{
		Client:   k8s.WrapClient(r.Client),
		Recorder: r.Recorder,
		Log:      r.Log,
		DClient:  k8s.WrapDClient(dClient),
		Scheme:   r.Scheme,
	}

	cacheStatus, err := r.Cache(ctx, &harborCluster, option).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile cache component.")
		return ctrl.Result{}, err
	}

	dbStatus, err := r.Database(ctx, &harborCluster, nil).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile database component.")
		return ctrl.Result{}, err
	}

	storageStatus, err := r.Storage(ctx, &harborCluster, nil).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile storage component.")
		return ctrl.Result{}, err
	}

	componentToStatus := make(map[goharborv1.Component]*lcm.CRStatus)
	componentToStatus[goharborv1.ComponentCache] = cacheStatus
	componentToStatus[goharborv1.ComponentDatabase] = dbStatus
	componentToStatus[goharborv1.ComponentStorage] = storageStatus
	// if components is not all ready, requeue the HarborCluster
	if !r.ComponentsAreAllReady(componentToStatus) {
		err = r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second * r.RequeueAfter,
		}, err
	}

	getRegistry := func() *string {
		if harborCluster.Spec.ImageSource != nil && harborCluster.Spec.ImageSource.Registry != "" {
			return &harborCluster.Spec.ImageSource.Registry
		}
		return nil
	}
	var imageGetter image.ImageGetter
	if imageGetter, err = image.NewImageGetter(getRegistry(), harborCluster.Spec.Version); err != nil {
		log.Error(err, "error when create ImageGetter.")
		return ctrl.Result{}, err
	}
	harborStatus, err := r.Harbor(ctx, &harborCluster, componentToStatus, &GetOptions{
		ImageGetter: imageGetter,
	}).Reconcile()
	if err != nil {
		log.Error(err, "error when reconcile harbor service.")
		return ctrl.Result{}, err
	}
	componentToStatus[goharborv1.ComponentHarbor] = harborStatus

	err = r.UpdateHarborClusterStatus(ctx, &harborCluster, componentToStatus)
	if err != nil {
		log.Error(err, "error when update harbor cluster status.")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// ServicesAreAllReady check whether these components(includes cache, db, storage) are all ready.
func (r *HarborClusterReconciler) ComponentsAreAllReady(serviceToMap map[goharborv1.Component]*lcm.CRStatus) bool {
	for _, status := range serviceToMap {
		if status == nil {
			return false
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
	for service, status := range componentToCRStatus {
		if status == nil {
			continue
		}
		harborClusterCondition, defaulted := r.getHarborClusterCondition(harborCluster, string(service))
		r.updateHarborClusterCondition(harborClusterCondition, status)
		if defaulted {
			harborCluster.Status.Conditions = append(harborCluster.Status.Conditions, *harborClusterCondition)
		}
	}
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
	conditionType string) (condition *goharborv1.HarborClusterCondition, defaulted bool) {
	for i := range harborCluster.Status.Conditions {
		condition = &harborCluster.Status.Conditions[i]
		if string(condition.Type) == conditionType {
			return condition, false
		}
	}
	return &goharborv1.HarborClusterCondition{
		Type:   goharborv1.HarborClusterConditionType(conditionType),
		Status: corev1.ConditionUnknown,
	}, true
}

func (r *HarborClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goharborv1.HarborCluster{}).
		Complete(r)
}
