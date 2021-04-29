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

	canarianv1alpha1 "canarian.io/plover/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PloverReconciler reconciles a Plover object
type PloverReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=canarian.canarian.io,resources=plovers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=canarian.canarian.io,resources=plovers/status,verbs=get;update;patch

func (r *PloverReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("plover", req.NamespacedName)

	var plover canarianv1alpha1.Plover
	if err := r.Get(ctx, req.NamespacedName, &plover); err != nil {
		log.Info("error getting object")
		return ctrl.Result{}, err
	}
	log.Info(
		"reconciling",
		"plover", req.NamespacedName,
		"active", plover.Spec.Active,
	)

	var pods corev1.PodList

	// Find all Pods

	if err := r.List(ctx, &pods); err != nil {
		log.Error(err, "unable to list  Pods")
		return ctrl.Result{Requeue: true}, err
	}
	// report failed and pending pods
	for _, item := range pods.Items {
		if item.Status.Phase != "Running" && item.Status.Phase != "Succeeded" {
			r.Log.Info(
				"Unhealthy Pod",
				"Namespace", item.Namespace,
				"Name", item.Name,
				"Phase", item.Status.Phase,
			)
		}

	}

	return ctrl.Result{}, nil
}

func (r *PloverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&canarianv1alpha1.Plover{}).
		Complete(r)
}
