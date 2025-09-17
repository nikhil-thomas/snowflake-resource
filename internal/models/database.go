package models

type Database struct {
	DatabaseName   string `db:"database_name"`
	Created        string `db:"created"`
	DATABASE_OWNER string `db:"database_owner"`
	Comment        string `db:"comment"`
}

type DatabaseStore interface {
	Get(name string) (*Database, error)
	List() ([]*Database, error)
	Create(database *Database) error
	Update(database *Database) error
	Delete(database *Database) error
}
