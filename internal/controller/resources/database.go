package resources

import "github.com/nikhil-thomas/snowflake-resource/api/v1alpha1"

// Database represents a database
// Implements the Resource interface
type Database struct {
	Name    string
	Comment string
}

// NewDatabase creates a new database from v1alpha1.SnowflakeResource
func NewDatabase(cr *v1alpha1.SnowflakeResource) *Database {
	db := &Database{
		Name: cr.Name,
	}
	for _, param := range cr.Spec.ObjectConfig {
		switch param.Key {
		case "comment":
			db.Comment = param.Value
		}
	}
	return &Database{}
}

// Create creates a database
func (d *Database) Create() error {
	return nil
}

// Delete deletes a database
func (d *Database) Delete() error {
	return nil
}

// Update updates a database
func (d *Database) Update() error {
	return nil
}

// List lists databases
func (d *Database) List() error {
	return nil
}
