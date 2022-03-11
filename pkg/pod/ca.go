package pod

// import (
// 	"github.com/kapkiai/nifiop/api/v1alpha1"
// 	"github.com/kapkiai/nifiop/controllers"
// 	appsv1 "k8s.io/api/apps/v1"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
// )

// func (r *controllers.NifiCAReconciler) CreatePod(cr v1alpha1.NifiCA) appsv1.Deployment {
// 	dep := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "nifi-ca",
// 			Namespace: cr.Namespace,
// 		},
// 		Spec: appsv1.DeploymentSpec{
// 			Template: corev1.PodTemplateSpec{
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{{
// 						Image: cr.Spec.ImageName,
// 						Name:  "nifi-ca",
// 						Ports: []corev1.ContainerPort{{
// 							ContainerPort: 9443,
// 							Name:          "nifi-ca-port",
// 						}},
// 					}},
// 				},
// 			},
// 		},
// 	}

// 	controllerutil.SetControllerReference(&cr, dep, r.Scheme)
// 	return dep
// }
