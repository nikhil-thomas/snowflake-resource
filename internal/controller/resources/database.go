package resources

import (
	"github.com/nikhil-thomas/snowflake-resource/api/v1alpha1"
	"github.com/nikhil-thomas/snowflake-resource/pkg/clients/snowflake"
)

// Database represents a database
// Implements the Resource interface
type DatabaseReconciler struct {
	C *snowflake.Client
}

// NewDatabaseReconciler creates a new database from v1alpha1.SnowflakeResource
func NewDatabaseReconciler(client *snowflake.Client) *DatabaseReconciler {
	return &DatabaseReconciler{
		C: client,
	}
}

// Create creates a database
func (d *DatabaseReconciler) Create(resource *v1alpha1.SnowflakeResource) error {
	query := "CREATE DATABASE" + resource.Name
	return d.C.Execute(query, true)
}

// Delete deletes a database
func (d *DatabaseReconciler) Delete(resource *v1alpha1.SnowflakeResource) error {
	query := "DROP DATABASE" + resource.Name
	return d.C.Execute(query, true)
}

// Update updates a database
func (d *DatabaseReconciler) Update(resource *v1alpha1.SnowflakeResource) error {
	return nil
}

// List lists databases
func (d *DatabaseReconciler) List() (*v1alpha1.SnowflakeResourceList, error) {
	return nil, nil
}

// Get gets a database
func (d *DatabaseReconciler) Get() (*v1alpha1.SnowflakeResource, error) {
	return nil, nil
}

// Reconcile reconciles a database
func (d *DatabaseReconciler) Reconcile(resource *v1alpha1.SnowflakeResource) error {
	return nil
}

// Status returns the status of a database
func (d *DatabaseReconciler) Status() v1alpha1.SnowflakeResourceStatus {
	return v1alpha1.SnowflakeResourceStatus{}
}
