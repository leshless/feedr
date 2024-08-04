package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	HOME_DIR_ERROR          = "could not retrieve home directory"
	NOT_EXIST_ERROR         = "missing ~/.config directory"
	PERMISSION_DENIED_ERROR = "permission denied while trying to access config files"
	UNKNOWN_ERROR           = "unknown"
)

var (
	configFile *os.File
	hostsFile  *os.File
)

type ConfigError struct {
	What string
}

func ConfigErrorNew(what string) ConfigError {
	return ConfigError{what}
}

func (err ConfigError) Error() string {
	return fmt.Sprintf("config error: %s", err.What)
}

// Create if not exist and open config.yaml and hosts.yaml in ~/.config/feedr folder
func OpenConfigFiles() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ConfigErrorNew(HOME_DIR_ERROR)
	}

	_, err = os.Stat(filepath.Join(homeDir, "/.config"))
	if err != nil {
		if os.IsNotExist(err) {
			return ConfigErrorNew(NOT_EXIST_ERROR)
		}
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	err = os.Mkdir(filepath.Join(homeDir, "/.config/feedr"), 0777)
	if err != nil && !os.IsExist(err) {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	configFile, err = os.Create(filepath.Join(homeDir, "/.config/feedr/config.yaml"))
	if err != nil && !os.IsExist(err) {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	hostsFile, err = os.Create(filepath.Join(homeDir, "/.config/feedr/hosts.yaml"))
	if err != nil && !os.IsExist(err) {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	return nil
}

// Close config.yaml and hosts.yaml
func CloseConfigFiles() {
	configFile.Close()
	hostsFile.Close()
}
