package controllers

import (
	"context"

	"github.com/kapkiai/nifiop/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *NifiCAReconciler) ensureDeployment(ctx context.Context, req ctrl.Request, instance *v1alpha1.NifiCA, dep *appsv1.Deployment) (*ctrl.Result, error) {
	found := &appsv1.Deployment{}

	err := r.Client.Get(ctx, types.NamespacedName{
		Name:      dep.Name,
		Namespace: dep.Namespace,
	}, found)

	logger := log.FromContext(ctx)

	if err != nil && errors.IsNotFound(err) {
		//Create the deployment
		logger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(ctx, dep)
		if err != nil {
			// Deployment failed
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &ctrl.Result{}, err
		}
		// Deployment was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		logger.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	return nil, nil

}
