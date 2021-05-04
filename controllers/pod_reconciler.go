package controllers

import (
	"context"
	"path/filepath"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PloverReconciler reconciles a Plover object
type PodReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *PodReconciler) Reconcile(pod corev1.Pod,
	containerName string,
	ploverType string) (ctrl.Result, error) {

	ctx := context.Background()
	r.Log.Info("in pdReconciler",
		"Pod", pod.Name,
		"ns", pod.Namespace,
		"container", containerName,
		"plover type", ploverType)

	if ploverType == "ImagePullBackOff" {
		//Try to find a secret
		for _, c := range pod.Spec.Containers {
			r.Log.Info(
				"found container",
				"name", c.Name,
			)
			if c.Name == containerName {
				imageRepo := filepath.Dir(c.Image)
				r.Log.Info(
					"the image repo",
					"path", imageRepo,
				)

				var secrets corev1.SecretList
				if err := r.List(ctx, &secrets); err != nil {
					r.Log.Error(err, "unable to list Secrets")
					return ctrl.Result{Requeue: true}, err
				}
				for _, secret := range secrets.Items {
					if secret.Type == corev1.SecretTypeDockercfg {
						r.Log.Info(
							"found docker secret",
							"secret", secret.Name,
							"namespace", secret.Namespace,
						)
					}
				}
			}
		}

	}

	return ctrl.Result{}, nil
}
