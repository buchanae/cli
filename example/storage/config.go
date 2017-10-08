package storage

import (
  "fmt"
  "os"
)

func DefaultConfig() Config {
	cwd, _ := os.Getwd()
  return Config{
    Local: LocalConfig{
      AllowedDirs: []string{cwd},
    },
  }
}

// Config describes configuration for all storage types
type Config struct {
	Local LocalConfig
	S3    S3Config
	GS    []GSConfig
	Swift SwiftConfig
}

// LocalConfig describes the directories Funnel can read from and write to
type LocalConfig struct {
	AllowedDirs []string
}

// Valid validates the LocalConfig configuration
func (l LocalConfig) Valid() bool {
	return len(l.AllowedDirs) > 0
}

// GSConfig describes configuration for the Google Cloud storage backend.
type GSConfig struct {
	AccountFile string
	FromEnv     bool
}

// Valid validates the GSConfig configuration.
func (g GSConfig) Valid() bool {
	return g.FromEnv || g.AccountFile != ""
}

func (g GSConfig) Validate() (errs []error) {
  fmt.Println("GS VALIDATE")
  return
}

// S3Config describes the directories Funnel can read from and write to
type S3Config struct {
	Key     string
	Secret  string
	FromEnv bool
}

// Validate validates the LocalConfig configuration
func (l S3Config) Validate() (errs []error) {
  if l.Key == "" {
    errs = append(errs, fmt.Errorf("key is empty"))
  }
  if l.Secret == "" {
    errs = append(errs, fmt.Errorf("secret is empty"))
  }
  return
}

// SwiftConfig configures the OpenStack Swift object storage backend.
type SwiftConfig struct {
	UserName   string
	Password   string
	AuthURL    string
	TenantName string
	TenantID   string
	RegionName string
}

// Valid validates the SwiftConfig configuration.
func (s SwiftConfig) Valid() bool {
	return s.UserName != "" && s.Password != "" && s.AuthURL != "" &&
		s.TenantName != "" && s.TenantID != "" && s.RegionName != ""
}
