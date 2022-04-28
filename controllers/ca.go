package controllers

import (
	"context"

	"github.com/kapkiai/nifiop/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NifiCAReconciler) caPod(cr *v1alpha1.NifiCA) *appsv1.Deployment {

	labels := map[string]string{
		"app": cr.Name,
	}
	matchlabels := map[string]string{
		"app": cr.Name,
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nifi-ca",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: matchlabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchlabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: cr.Spec.ImageName,
						// ImagePullPolicy: "",
						Name: "nifi-ca",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 9443,
							Name:          "nifi-ca-port",
						}},
						Command: []string{"bash", "-ce", "./bin/tls-toolkit.sh server -c localhost -t 123456789123456778 -p 9443"},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.Scheme)
	return dep
}

func (r *NifiCAReconciler) isCaUp(ctx context.Context, v *v1alpha1.NifiCA) bool {
	dep := &appsv1.Deployment{}

	err := r.Client.Get(ctx, types.NamespacedName{
		Name:      "nifi-ca",
		Namespace: v.Namespace,
	}, dep)
	logger := log.FromContext(ctx)

	if err != nil {
		logger.Error(err, "Deployment nifi-ca not found")
		return false
	}

	if dep.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
