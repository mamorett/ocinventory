// Package config provides helpers to build an OCI ConfigurationProvider
// from a standard ~/.oci/config file and a named profile.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oracle/oci-go-sdk/v65/common"
)

// NewProvider returns a ConfigurationProvider for the given profile.
// configPath may be empty, in which case ~/.oci/config is used.
func NewProvider(profile, configPath string) (common.ConfigurationProvider, error) {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home dir: %w", err)
		}
		configPath = filepath.Join(home, ".oci", "config")
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("OCI config file not found at %s: %w", configPath, err)
	}

	p := common.CustomProfileConfigProvider(configPath, profile)

	// Validate early so the user gets a clear error message.
	if _, err := p.TenancyOCID(); err != nil {
		return nil, fmt.Errorf("profile %q in %s is invalid or missing tenancy: %w", profile, configPath, err)
	}

	return p, nil
}
