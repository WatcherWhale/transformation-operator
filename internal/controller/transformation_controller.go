/*
Copyright 2023 WatcherWhale.

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
	"bytes"
	"context"

	"text/template"

	"github.com/Masterminds/sprig"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	transformationsv1alpha1 "github.com/WatcherWhale/transformation-operator/api/v1alpha1"
)

// TransformationReconciler reconciles a Transformation object
type TransformationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=transformations.transformations.go,resources=transformations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=transformations.transformations.go,resources=transformations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=transformations.transformations.go,resources=transformations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Transformation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *TransformationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	tf := transformationsv1alpha1.Transformation{}

	if err := r.Get(ctx, req.NamespacedName, &tf); err != nil {
		logger.Error(err, "Failed to get resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	dataMap := make(map[string]map[string]string)

	for _, source := range tf.Spec.Sources {
		// Load source values
		switch source.Name {
		case "ConfigMap":
			dataMap[source.Name] = make(map[string]string)
		case "Secret":
			dataMap[source.Name] = make(map[string]string)
		}
	}

	templater := template.New("gotpl").Funcs(sprig.FuncMap())

	// Template the shit out of it
	templatedMap := make(map[string]string)

	for key, tplStr := range tf.Spec.Template {
		var buf bytes.Buffer
		tpl, err := templater.Parse(tplStr)

		if err != nil {
			logger.Error(err, "Could not parse template")
			return ctrl.Result{}, err
		}

		err = tpl.Execute(&buf, dataMap)

		if err != nil {
			logger.Error(err, "Could not apply template")
			return ctrl.Result{}, err
		}

		templatedMap[key] = buf.String()
	}

	return ctrl.Result{}, nil
}

func (r *TransformationReconciler) getConfigMapKeys(ctx context.Context) map[string]string {
	return make(map[string]string)
}

func (r *TransformationReconciler) getSecretKeys(ctx context.Context) map[string]string {
	return make(map[string]string)
}

// SetupWithManager sets up the controller with the Manager.
func (r *TransformationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&transformationsv1alpha1.Transformation{}).
		Complete(r)
}
