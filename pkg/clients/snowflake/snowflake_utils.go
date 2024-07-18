package snowflake

import (
	"github.com/DATA-DOG/go-sqlmock"
)

func NewMockClient() (*Client, sqlmock.Sqlmock, error) {
	if SnowflakeClient != nil {
		return SnowflakeClient, nil, nil
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	SnowflakeClient = &Client{
		conn:  db,
		Close: db.Close,
	}
	return SnowflakeClient, mock, nil
}
