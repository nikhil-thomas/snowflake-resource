package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*DiceConfig, error) {
	v := viper.New()

	v.SetConfigName("config")   // Base name of your config file (without extension)
	v.SetConfigType("yaml")     // File type
	v.AddConfigPath(configPath) // Default path to look for the file

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg DiceConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	//env := viper.GetString("environment")
	env := "snowflake-hacks"
	cfg.setCurrentWorkingDirOrPanic()
	cfg.setCurrentEnvironmentOrPanic(env)
	logLevel := viper.GetString("log-level")
	cfg.LogLevel = logLevel

	return &cfg, nil
}

func (d *DiceConfig) setCurrentWorkingDirOrPanic() {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	d.CurrentWorkingDir = currentWorkingDir
}

func (d *DiceConfig) setCurrentEnvironmentOrPanic(env string) {
	for _, e := range d.Environments {
		if e.Name == env {
			d.CurrentEnvironment = e
		}
	}
	if d.CurrentEnvironment == nil {
		panic(fmt.Sprintf("environment %s not found in config", env))
	}
	d.CurrentEnvironment.Snowflake.loadPrivateKeyOrPanic()
}

func (d *DiceConfig) PlatformConfigPath() string {
	configPath := filepath.Join(d.PlatformConfigRootDir, d.CurrentEnvironment.Name)
	if filepath.IsAbs(configPath) {
		return configPath
	}
	return filepath.Join(d.CurrentWorkingDir, configPath)
}

func dir() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(b), "..", "..")
}

func (s *Snowflake) loadPrivateKeyOrPanic() {
	privateKeyPath := s.PrivateKeyPath
	if !filepath.IsAbs(privateKeyPath) {
		privateKeyPath = filepath.Join(dir(), privateKeyPath)
	}
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read private key: %v, %v", s.PrivateKey, err))
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(fmt.Sprintf("No PEM data found in private key file, %v", s.PrivateKeyPath))
	}

	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse private key: %v, %v", s.PrivateKeyPath, err))
	}

	rsaPrivate, ok := private.(*rsa.PrivateKey)
	if !ok {
		panic("Key is not an RSA private key")
	}

	s.PrivateKey = rsaPrivate
}

func FromCtx(ctx context.Context) *DiceConfig {
	return ctx.Value(CtxConfigKey).(*DiceConfig)
}

func SnowflakeConfigFromCtx(ctx context.Context) *Snowflake {
	return FromCtx(ctx).CurrentEnvironment.Snowflake
}

func VaultConfigFromCtx(ctx context.Context) *Vault {
	return FromCtx(ctx).Vault
}
