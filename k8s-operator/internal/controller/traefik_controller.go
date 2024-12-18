/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"myoperator/internal/util"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ingressv1beta1 "myoperator/api/v1beta1"
)

// TraefikReconciler reconciles a Traefik object
type TraefikReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ingress.demo.morawskim.pl,resources=traefiks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ingress.demo.morawskim.pl,resources=traefiks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ingress.demo.morawskim.pl,resources=traefiks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Traefik object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *TraefikReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here
	list, err := fetchTraefikResources(r.Client)
	if err != nil {
		logger.Error(err, "Failed to fetch Traefik resources")
	}

	resource := ingressv1beta1.Traefik{}
	err = r.Client.Get(context.Background(), req.NamespacedName, &resource)
	if err != nil {
		logger.Info("Requested resource not found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	err = createIndexFile(list, r.Client, resource.Spec.TargetNamespace, resource.Spec.TargetConfigMapName, resource.Spec.TargetDeploymentName)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Index file created")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TraefikReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetLogger().Info("setting up manager")
	err := AddToScheme(mgr.GetScheme())

	if err != nil {
		return errors.Wrap(err, "failed during adding traefik scheme")
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&ingressv1beta1.Traefik{}).
		Complete(r)
}

func fetchTraefikResources(k8sClient client.Client) ([]util.TraefikItem, error) {
	ctx := context.Background()
	traefikRouteList := &IngressRouteList{}
	err := k8sClient.List(ctx, traefikRouteList, &client.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to list TraefikRoute: %w", err)
	}

	list := make([]util.TraefikItem, 0, len(traefikRouteList.Items))

	for _, traefikRoute := range traefikRouteList.Items {
		list = append(list, util.TraefikItem{
			Name:      traefikRoute.Name,
			Namespace: traefikRoute.Namespace,
			Url:       "https://todo.example.com?name=" + traefikRoute.Name,
		})
	}

	return list, nil
}

func createIndexFile(data []util.TraefikItem, client client.Client, targetNamespace, targetConfigMapName, targetDeploymentName string) error {
	if len(targetNamespace) == 0 || len(targetConfigMapName) == 0 || len(targetDeploymentName) == 0 {
		return nil
	}

	indexFile, err := util.GenerateIndexFile(data)

	if err != nil {
		return err
	}

	err = util.CreateConfigMap(client, targetNamespace, targetConfigMapName, indexFile)
	if err != nil {
		return err
	}

	err = util.RestartDeployment(client, targetNamespace, targetDeploymentName)
	if err != nil {
		return err
	}

	return nil
}
