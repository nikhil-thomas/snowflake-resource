package controller

import (
	"github.com/nikhil-thomas/snowflake-resource/api/v1alpha1"
	"github.com/nikhil-thomas/snowflake-resource/internal/controller/resources"
)

type Resource interface {
	Create() error
	Delete() error
	Update() error
	List() error
}

func ResourceFromCR(cr *v1alpha1.SnowflakeResource) Resource {
	switch cr.Spec.ObjectType {
	case "database":
		return resources.NewDatabase(cr)
	}
	return nil
}
