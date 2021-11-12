package controllers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	rand "k8s.io/apimachinery/pkg/util/rand"
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
	r.Log.Info("in podReconciler",
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
				secretFound := false
				var nsSecrets, allSecrets corev1.SecretList
				//first - look for secrets in pod namespace
				if err := r.List(ctx, &nsSecrets, client.InNamespace(pod.Namespace)); err != nil {
					r.Log.Error(err, "unable to list Secrets in namespace")
					return ctrl.Result{RequeueAfter: 60 * time.Second}, err
				}
				for _, secret := range nsSecrets.Items {
					if r.checkSecretCandidate(secret, imageRepo) {
						r.Log.Info(
							"found appropriate docker secret in namespace",
							"secret", secret.Name,
							"namespace", secret.Namespace,
							"image repo", imageRepo,
							"resolution", "don't know yet",
						)
						secretFound = true
						r.Log.Info("Going to patch pod",
							"pod name", pod.Name,
							"secret", secret.Name,
						)
						//try to patch pod
						patch := []byte(fmt.Sprintf(`{"spec":{"imagePullSecrets":[{"name": "%s"}]}}`,
							secret.Name))
						if err := r.Patch(ctx, &corev1.Pod{
							ObjectMeta: metav1.ObjectMeta{
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
						},
							client.RawPatch(types.StrategicMergePatchType, patch)); err != nil {
							r.Log.Error(err, "Couldn't patch pod")
						}
					}
				}
				if !(secretFound) {
					r.Log.Info("Didn't find appropriate secrets in ns")
					if err := r.List(ctx, &allSecrets); err != nil {
						r.Log.Error(err, "unable to list Secrets")
						return ctrl.Result{RequeueAfter: 60 * time.Second}, err
					}
					for _, secret := range allSecrets.Items {
						if r.checkSecretCandidate(secret, imageRepo) {
							r.Log.Info(
								"found appropriate docker secret",
								"secret", secret.Name,
								"namespace", secret.Namespace,
								"image repo", imageRepo,
							)

							newSecret := &corev1.Secret{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: pod.Namespace,
									Name:      "plover-" + rand.String(8),
									Labels: map[string]string{
										"plover": "true",
									},
								},
								Data: map[string][]byte{
									".dockerconfigjson": secret.Data[".dockerconfigjson"],
								},
								Type: secret.Type,
							}
							if err := r.Create(ctx, newSecret); err != nil {
								r.Log.Error(err, "unable to create Secret")
								return ctrl.Result{Requeue: true}, err
							} //if secret creation failed
							r.Log.Info(
								"Created docker secret",
								"secret", newSecret.Name,
								"namespace", newSecret.Namespace,
								"image repo", imageRepo,
							)
						}
						//if it's a relevant secret
					} //loop on all secrets
				}
			} //if relevant container
		} //loop on containers
	} //if ImagePullBackofff
	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

func (r *PodReconciler) checkSecretCandidate(secret corev1.Secret, imageRepo string) bool {
	if secret.Type == "kubernetes.io/dockerconfigjson" || secret.Type == "kubernetes.io/dockercfg" {
		//trim repo url to first slash
		// because secrets are defined for base url
		// e.g for https://myrepo.jfrog.io and not https://myprepo.jfrog.io/docker-virtual
		baseRepo := imageRepo
		if idx := strings.IndexByte(imageRepo, '/'); idx >= 0 {
			baseRepo = imageRepo[:idx]
		}
		r.Log.Info("Checking secret", "secret name:", secret.Name, "image repo", baseRepo)
		return strings.Contains(string(secret.Data[".dockerconfigjson"]), baseRepo)
	}
	return false
}
