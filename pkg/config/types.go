package config

import "crypto/rsa"

type contextKey string

type DiceConfig struct {
	Environments          []*Environment `mapstructure:"environments"`
	Vault                 *Vault         `mapstructure:"vault"`
	PlatformConfigRootDir string         `mapstructure:"platform-config-root-dir"`
	CurrentEnvironment    *Environment
	CurrentWorkingDir     string
	LogLevel              string
}

type Snowflake struct {
	Account         string `mapstructure:"account"`
	User            string `mapstructure:"user"`
	Role            string `mapstructure:"role"`
	Region          string `mapstructure:"region"`
	DefaultDatabase string `mapstructure:"default-database"`
	Warehouse       string `mapstructure:"warehouse"`
	PrivateKeyPath  string `mapstructure:"private-key-path"`
	PrivateKey      *rsa.PrivateKey
}

type Vault struct {
	Address     string `mapstructure:"address"`
	RoleID      string `mapstructure:"role-id"`
	SecretID    string `mapstructure:"secret-id"`
	SecretsRoot string `mapstructure:"secrets-root"`
}

type Environment struct {
	Name      string `mapstructure:"name"`
	Snowflake *Snowflake
}
