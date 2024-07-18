package snowflake

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/nikhil-thomas/snowflake-resource/pkg/config"
	"github.com/snowflakedb/gosnowflake"
)

var fncCount = 0

func NewClient(opts ClientConfig) (*Client, error) {
	if SnowflakeClient != nil {
		return SnowflakeClient, nil
	}
	opts.setDefaults()
	dsn, err := gosnowflake.DSN(&opts.Config)
	if err != nil {
		return nil, err
	}
	conn, err := sql.Open("snowflake", dsn)
	if err != nil {
		return nil, err
	}
	SnowflakeClient = &Client{
		conn:  conn,
		Close: conn.Close,
	}
	err = conn.Ping()
	if err != nil {
		return nil, err

	}
	return SnowflakeClient, nil
}

func NewClientFromPlatformConfig(cfg *config.Snowflake) (*Client, error) {
	cc := NewClientConfigFromPlatformConfig(cfg)
	return NewClient(*cc)
}

func SnowflakeClientOrPanic(cfg *config.Snowflake) *Client {
	c, err := NewClientFromPlatformConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Snowflake client: %v", err))
	}
	return c
}

func NewClientConfig() *ClientConfig {
	cc := &ClientConfig{}
	cc.setDefaults()
	return cc
}

func NewClientConfigFromPlatformConfig(cfg *config.Snowflake) *ClientConfig {
	cc := NewClientConfig()
	cc.Account = cfg.Account
	cc.User = cfg.User
	cc.Role = cfg.Role
	cc.Region = cfg.Region
	cc.Database = cfg.DefaultDatabase
	cc.Warehouse = cfg.Warehouse
	cc.PrivateKey = cfg.PrivateKey
	return cc
}

func (c *ClientConfig) setDefaults() {
	c.Authenticator = gosnowflake.AuthTypeJwt
}

func (c *Client) Query(query string, verbose bool) (*sql.Rows, error) {
	fncCount++
	colorReset := "\033[0m"
	color := "\033[34m"
	if verbose {
		fmt.Println("\n", strings.Repeat("-", 20), "\n Step: ", fncCount, "Query : ", color, mask(query), colorReset, "\n", strings.Repeat("-", 20))
	}
	rows, err := c.conn.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *Client) ExecuteBatch(queries []string, verbose bool) error {
	start := time.Now()
	for _, query := range queries {
		err := c.Execute(query, verbose)
		if err != nil {
			return err
		}
	}
	fmt.Println("Time elapsed: ", time.Since(start))
	return nil
}

func (c *Client) Execute(query string, verbose bool) error {
	fncCount++
	colorReset := "\033[0m"
	color := "\033[34m"
	if verbose {
		fmt.Println("\n", strings.Repeat("-", 20), "\n Step: ", fncCount, "Query : ", color, mask(query), colorReset, "\n", strings.Repeat("-", 20))
	}
	RowsAffected, err := execute(c, query)
	if err != nil {
		return err
	}
	fmt.Println("Rows Affected :", RowsAffected)
	return nil
}

func (c *Client) ExecuteWithResponse(query string, verbose bool) (int64, error) {
	fncCount++
	colorReset := "\033[0m"
	color := "\033[34m"
	if verbose {
		fmt.Println("\n", strings.Repeat("-", 20), "\n Step: ", fncCount, "Query : ", color, mask(query), colorReset, "\n", strings.Repeat("-", 20))
	}
	rowsAffected, err := execute(c, query)
	if err != nil {
		return 0, err
	}
	fmt.Println("Rows Affected :", rowsAffected)
	return rowsAffected, nil
}

func mask(query string) string {
	return regexp.MustCompile("RSA_PUBLIC_KEY = '(.*?)'").ReplaceAllString(query, "RSA_PUBLIC_KEY = <REDACTED>")
}

func execute(c *Client, query string) (int64, error) {
	result, err := c.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	RowsAffected, _ := result.RowsAffected()
	return RowsAffected, nil
}

func (c *Client) Ping() error {
	return c.conn.Ping()
}
