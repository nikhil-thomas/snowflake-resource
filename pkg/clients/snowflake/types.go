package snowflake

import (
	"database/sql"

	"github.com/snowflakedb/gosnowflake"
)

type ClientConfig struct {
	gosnowflake.Config
}

type Client struct {
	conn  *sql.DB
	Close func() error
}

var (
	SnowflakeClient *Client
	Close           func() error
)
