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

	query := "show users"
	rows, err := r.SnowflakeClient.Query(query, true)
	printRows(rows)
	if err != nil {
		log.Log.Error(err, "unable to query Snowflake")
		return ctrl.Result{}, err

	}

	log.Log.Info("reconciling SnowflakeResource", "name", res.Name)

	return ctrl.Result{}, nil
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
