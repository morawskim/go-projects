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
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ingressv1beta1 "myoperator/api/v1beta1"
	"myoperator/internal/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

const ResourceLabel string = "demo.morawskim.pl/traefik-url"
const ResourceManagedLabel string = "demo.morawskim.pl/manged-by"
const Separator string = "/"

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
	logger.Info("Start Traefik reconciliation")
	resource := ingressv1beta1.Traefik{}
	err := r.Client.Get(context.Background(), req.NamespacedName, &resource)
	if err == nil {
		return r.reconcileTraefik(&resource, logger)
	}

	route := IngressRoute{}
	err = r.Client.Get(ctx, req.NamespacedName, &route)
	if err == nil {
		if val, ok := route.Annotations[ResourceManagedLabel]; ok {
			chunks := strings.SplitN(val, Separator, 2)
			if len(chunks) == 2 {
				err = r.Client.Get(context.Background(), client.ObjectKey{
					Namespace: chunks[0],
					Name:      chunks[1],
				}, &resource)

				if err == nil {
					return r.reconcileTraefik(&resource, logger)
				}
			}
		}
	}

	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *TraefikReconciler) reconcileTraefik(resource *ingressv1beta1.Traefik, logger logr.Logger) (ctrl.Result, error) {
	list, err := fetchTraefikResources(r.Client, resource.Spec.LookForLabel)
	if err != nil {
		logger.Error(err, "Failed to fetch Traefik resources")
	}

	err = createIndexFile(list, r.Client, resource.Spec.TargetNamespace, resource.Spec.TargetConfigMapName, resource.Spec.TargetDeploymentName, resource.Name, resource.Namespace)
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
		Watches(
			&IngressRoute{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.NewPredicateFuncs(func(object client.Object) bool {
				labels := object.GetLabels()
				_, ok := labels[ResourceLabel]
				return ok
			})),
		).Complete(r)
}

func fetchTraefikResources(k8sClient client.Client, label string) ([]util.TraefikItem, error) {
	ctx := context.Background()
	traefikRouteList := &IngressRouteList{}
	err := k8sClient.List(ctx, traefikRouteList, &client.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to list TraefikRoute: %w", err)
	}

	list := make([]util.TraefikItem, 0, len(traefikRouteList.Items))

	for _, traefikRoute := range traefikRouteList.Items {
		if value, ok := traefikRoute.Labels[label]; ok {
			list = append(list, util.TraefikItem{
				Name:      traefikRoute.Name,
				Namespace: traefikRoute.Namespace,
				Url:       "https://" + value,
			})
		}
	}

	return list, nil
}

func createIndexFile(data []util.TraefikItem, k8sClient client.Client, targetNamespace, targetConfigMapName, targetDeploymentName, resourceName, resourceNamespace string) error {
	if len(targetNamespace) == 0 || len(targetConfigMapName) == 0 || len(targetDeploymentName) == 0 {
		return nil
	}

	indexFile, err := util.GenerateIndexFile(data)

	if err != nil {
		return err
	}

	ctx := context.Background()
	for _, traefikItem := range data {
		route := IngressRoute{}
		// see https://www.rfc-editor.org/rfc/rfc6901#section-3 for why we replace "/"
		patchData := []byte(fmt.Sprintf(
			`[ {"op": "replace", "path": "/metadata/annotations/%s", "value": "%s"} ]`,
			strings.ReplaceAll(ResourceManagedLabel, "/", "~1"),
			resourceNamespace+Separator+resourceName,
		))
		err = k8sClient.Get(ctx, client.ObjectKey{
			Namespace: traefikItem.Namespace,
			Name:      traefikItem.Name,
		}, &route)

		if err != nil {
			return errors.Wrap(err, "failed to get traefik route resource")
		}

		err = k8sClient.Patch(ctx, &route, client.RawPatch(types.JSONPatchType, patchData))
		if err != nil {
			return errors.Wrap(err, "failed to patch traefik resource")
		}
	}

	err = util.CreateConfigMap(k8sClient, targetNamespace, targetConfigMapName, indexFile)
	if err != nil {
		return err
	}

	err = util.RestartDeployment(k8sClient, targetNamespace, targetDeploymentName)
	if err != nil {
		return err
	}

	return nil
}
