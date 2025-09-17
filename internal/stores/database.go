package stores

import (
	"github.com/nikhil-thomas/snowflake-resource/internal/models"
	"github.com/nikhil-thomas/snowflake-resource/pkg/clients/snowflake"
)

type Database struct {
	client *snowflake.Client
}

// compile-time assertion to ensure that the *Database type
// correctly implements the DatabaseStore interface
var _ models.DatabaseStore = (*Database)(nil)

func NewDatabase(client *snowflake.Client) *Database {
	return &Database{
		client: client,
	}
}

// implement DatabaseStore interface
func (d *Database) Get(name string) (*models.Database, error) {
	return nil, nil
}

func (d *Database) List() ([]*models.Database, error) {
	return nil, nil
}

func (d *Database) Create(database *models.Database) error {
	return nil
}

func (d *Database) Update(database *models.Database) error {
	return nil
}

func (d *Database) Delete(database *models.Database) error {
	return nil
}
