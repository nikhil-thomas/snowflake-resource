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
	"database/sql"
	"fmt"
	"strings"

	"github.com/nikhil-thomas/snowflake-resource/pkg/clients/snowflake"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	snowflakev1alpha1 "github.com/nikhil-thomas/snowflake-resource/api/v1alpha1"
)

// SnowflakeResourceReconciler reconciles a SnowflakeResource object
type SnowflakeResourceReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	SnowflakeClient *snowflake.Client
}

// +kubebuilder:rbac:groups=snowflake.my.domain,resources=snowflakeresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=snowflake.my.domain,resources=snowflakeresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=snowflake.my.domain,resources=snowflakeresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SnowflakeResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *SnowflakeResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	res := &snowflakev1alpha1.SnowflakeResource{}
	if err := r.Get(ctx, req.NamespacedName, res); err != nil {
		log.Log.Error(err, "unable to fetch SnowflakeResource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	objectReconciler := TypeSpecificReconciler(res, r.SnowflakeClient)
	if objectReconciler == nil {
		log.Log.Error(fmt.Errorf("unable to find reconciler for resource"), "unable to find reconciler for resource")
		return ctrl.Result{}, errors.New("unable to find reconciler for resource")
	}

	result, err := r.handleFinalizer(ctx, res, objectReconciler)
	if err != nil {
		log.Log.Error(err, "unable to handle finalizer")
		return ctrl.Result{}, err
	}
	if result.Requeue {
		// requeue
		return result, nil
	}

	// handle reconciliation logic

	// check if resource exists
	result, err = r.createObjectIfNotExists(err, objectReconciler, res)
	if err != nil {
		log.Log.Error(err, "unable to create object")
		return ctrl.Result{}, err
	}

	if result.Requeue {
		// requeue
		return result, nil
	}

	res.Status = objectReconciler.Status()

	// update status
	err = r.Status().Update(ctx, res)
	if err != nil {
		log.Log.Error(err, "unable to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SnowflakeResourceReconciler) createObjectIfNotExists(err error, objectReconciler ResourceReconciler, res *snowflakev1alpha1.SnowflakeResource) (ctrl.Result, error) {
	_, err = objectReconciler.Get()
	if err == nil {
		return ctrl.Result{}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Log.Error(err, "unable to get resource")
		return ctrl.Result{}, err
	}
	// resource not found
	err = objectReconciler.Create(res)
	if err != nil {
		log.Log.Error(err, "unable to create resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{
		Requeue: true,
	}, nil
}

func (r *SnowflakeResourceReconciler) handleFinalizer(ctx context.Context, res *snowflakev1alpha1.SnowflakeResource, objectReconciler ResourceReconciler) (ctrl.Result, error) {
	// if deletion timestamp not set on res handle finalizer logic
	// implies not a deletion request
	if res.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.addFinalizerIfNotExists(ctx, res)
	}

	return r.handleDeletion(ctx, res, objectReconciler)
}

func (r *SnowflakeResourceReconciler) handleDeletion(ctx context.Context, res *snowflakev1alpha1.SnowflakeResource, objectReconciler ResourceReconciler) (ctrl.Result, error) {
	if !containsString(res.ObjectMeta.Finalizers, "finalizer.snowflake.my.domain") {
		// if finalizer not present on res return
		// implies deletion request but finalizer already removed
		return ctrl.Result{}, nil
	}

	// handle deletion logic
	err := objectReconciler.Delete(res)
	if err != nil {
		log.Log.Error(err, "unable to delete resource")
		return ctrl.Result{}, err
	}

	// remove finalizer
	res.ObjectMeta.Finalizers = removeString(res.ObjectMeta.Finalizers, "finalizer.snowflake.my.domain")
	err = r.Update(ctx, res)
	if err != nil {
		log.Log.Error(err, "unable to update resource")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *SnowflakeResourceReconciler) addFinalizerIfNotExists(ctx context.Context, res *snowflakev1alpha1.SnowflakeResource) (ctrl.Result, error) {
	if containsString(res.ObjectMeta.Finalizers, "finalizer.snowflake.my.domain") {
		return ctrl.Result{
			// no requeue as the finalizer already exists
			Requeue: false,
		}, nil
	}
	res.ObjectMeta.Finalizers = append(res.ObjectMeta.Finalizers, "finalizer.snowflake.my.domain")
	err := r.Update(ctx, res)
	if err != nil {
		log.Log.Error(err, "unable to update resource")
		return ctrl.Result{}, err
	}
	return ctrl.Result{
		// requeue as the finalizer was added and the object was updated
		Requeue: true,
	}, nil
}

func printRows(rows *sql.Rows) {
	columns, err := rows.Columns()
	if err != nil {
		log.Log.Error(err, "unable to get columns")
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(sql.NullString)
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Log.Error(err, "unable to scan row")
			return
		}

		for i, col := range values {
			if col.(*sql.NullString).Valid {
				fmt.Printf("%s: %s\n", columns[i], col.(*sql.NullString).String)
			} else {
				fmt.Printf("%s: NULL\n", columns[i])
			}
		}
		fmt.Println(strings.Repeat("-", 20))
	}

	if err := rows.Err(); err != nil {
		log.Log.Error(err, "error in rows")
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SnowflakeResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&snowflakev1alpha1.SnowflakeResource{}).
		Complete(r)
}
