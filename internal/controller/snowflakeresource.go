package controller

import (
	"github.com/nikhil-thomas/snowflake-resource/api/v1alpha1"
	"github.com/nikhil-thomas/snowflake-resource/internal/controller/resources"
	"github.com/nikhil-thomas/snowflake-resource/pkg/clients/snowflake"
)

type ResourceReconciler interface {
	Create(resource *v1alpha1.SnowflakeResource) error
	Delete(resource *v1alpha1.SnowflakeResource) error
	Update(resource *v1alpha1.SnowflakeResource) error
	List() (*v1alpha1.SnowflakeResourceList, error)
	Get() (*v1alpha1.SnowflakeResource, error)
	Status() v1alpha1.SnowflakeResourceStatus
	Reconcile(resource *v1alpha1.SnowflakeResource) error
}

func TypeSpecificReconciler(cr *v1alpha1.SnowflakeResource, c *snowflake.Client) ResourceReconciler {
	switch cr.Spec.ObjectType {
	case "database":
		return resources.NewDatabaseReconciler(c)
	}
	return nil
}
