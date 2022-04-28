/*
Copyright 2022.

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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kapkiai/nifiop/api/v1alpha1"
	nifitoolkitv1alpha1 "github.com/kapkiai/nifiop/api/v1alpha1"
)

// NifiCAReconciler reconciles a NifiCA object
type NifiCAReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nifitoolkit.safaricom.et,resources=nificas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nifitoolkit.safaricom.et,resources=nificas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nifitoolkit.safaricom.et,resources=nificas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NifiCA object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *NifiCAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)
	// TODO(user): your logic here
	reqLogger := log.Log.WithValues("Request.Namespace", req.Namespace, "Req.Name", req.Name)
	reqLogger.Info("Reconcilling NiFiCA...")

	ca := &v1alpha1.NifiCA{}
	err := r.Client.Get(ctx, req.NamespacedName, ca)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Object not found, returning")
			return ctrl.Result{}, nil
		}
		//Error reading the object - requeue the object
		reqLogger.Error(err, "Error reading the object - requeing the object")
		return ctrl.Result{}, nil
	}

	var result *ctrl.Result

	result, err = r.ensureDeployment(ctx, req, ca, r.caPod(ca))

	if result != nil {
		reqLogger.Info("Deployment found")
		return *result, err
	}

	nifiCaRunning := r.isCaUp(ctx, ca)

	if !nifiCaRunning {
		// If NiFi CA isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)
		reqLogger.Info(fmt.Sprintf("Nifi CA not running,waiting for %s", delay))
		return ctrl.Result{RequeueAfter: delay}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NifiCAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nifitoolkitv1alpha1.NifiCA{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
