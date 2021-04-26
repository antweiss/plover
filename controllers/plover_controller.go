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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	canarianv1alpha1 "canarian.io/plover/api/v1alpha1"
)

// +kubebuilder:rbac
type PloverController struct {
	client.Client
	scheme *runtime.Scheme
}

func Add(mgr manager.Manager) error {
	// Create a new Controller
	c, err := controller.New("plover-controller", mgr,
		controller.Options{Reconciler: &PloverController{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		}})
	if err != nil {
		return err
	}

	// Watch for changes to Plover
	err = c.Watch(
		&source.Kind{Type: &canarianv1alpha1.Plover{}},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to Deployments created by a ContainerSet and trigger a Reconcile for the owner
	err = c.Watch(
		&source.Kind{Type: &v1.Pod{}},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// PloverReconciler reconciles a Plover object
// type PloverReconciler struct {
// 	client.Client
// 	Log    logr.Logger
// 	Scheme *runtime.Scheme
// }

// +kubebuilder:rbac:groups=canarian.canarian.io,resources=plovers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=canarian.canarian.io,resources=plovers/status,verbs=get;update;patch

var _ reconcile.Reconciler = &PloverController{}

func (r *ReconcilePlover) Reconcile(req ctrl.Request) (ctrl.Result, error) {
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

	return ctrl.Result{}, nil
}

// func (r *PloverReconciler) SetupWithManager(mgr ctrl.Manager) error {
// 	return ctrl.NewControllerManagedBy(mgr).
// 		For(&canarianv1alpha1.Plover{}).
// 		Complete(r)
// }
