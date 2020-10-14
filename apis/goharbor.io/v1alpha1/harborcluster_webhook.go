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

package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var harborclusterlog = logf.Log.WithName("harborcluster-resource")

func (r *HarborCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-goharbor-io-v1-harborcluster,mutating=true,failurePolicy=fail,groups=goharbor.io,resources=harborclusters,verbs=create;update,versions=v1alpha1,name=mharborcluster.kb.io

var _ webhook.Defaulter = &HarborCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HarborCluster) Default() {
	harborclusterlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-goharbor-io-v1-harborcluster,mutating=false,failurePolicy=fail,groups=goharbor.io,resources=harborclusters,versions=v1alpha1,name=vharborcluster.kb.io

var _ webhook.Validator = &HarborCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HarborCluster) ValidateCreate() error {
	harborclusterlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HarborCluster) ValidateUpdate(old runtime.Object) error {
	harborclusterlog.Info("validate update", "name", r.Name)

	return r.ValidateComponentKind(old)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HarborCluster) ValidateDelete() error {
	harborclusterlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return r.ValidateCertificateIssuerRef()
}

func (r *HarborCluster) ValidateComponentKind(old runtime.Object) error {
	oldHarbor := old.(*HarborCluster)
	if r.Spec.Redis.Kind != oldHarbor.Spec.Redis.Kind ||
		r.Spec.Database.Kind != oldHarbor.Spec.Database.Kind ||
		r.Spec.Storage.Kind != oldHarbor.Spec.Storage.Kind {
		return errors.New("service kind switching is not supported")
	}
	return nil
}

func (r *HarborCluster) ValidateCertificateIssuerRef() error {
	if len(r.Spec.CertificateIssuerRef.Name) < 1 {
		return errors.New("CertificateIssuerRef name required")
	}
	return nil
}
