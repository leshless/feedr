package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const (
	HOME_DIR_ERROR          = "could not retrieve home directory"
	NOT_EXIST_ERROR         = "missing ~/.config directory"
	PERMISSION_DENIED_ERROR = "permission denied while trying to access config files"
	READ_ERROR              = "count not read config file"
	YAML_FORMAT_ERROR       = "bad config file format"
	UNKNOWN_ERROR           = "unknown"
)

var (
	configFile  *os.File
	sourcesFile *os.File
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

// Create if not exist and open config.yml and sourcess.yml in ~/.config/feedr folder
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

	_, err = os.Stat(filepath.Join(homeDir, "/.config/feedr/config.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			configFile, err = os.Create(filepath.Join(homeDir, "/.config/feedr/config.yml"))
			if err != nil && !os.IsExist(err) {
				if os.IsPermission(err) {
					return ConfigErrorNew(PERMISSION_DENIED_ERROR)
				}
				return ConfigErrorNew(UNKNOWN_ERROR)
			}
		}
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	} else {
		configFile, err = os.Open(filepath.Join(homeDir, "/.config/feedr/config.yml"))
		if err != nil && !os.IsExist(err) {
			if os.IsPermission(err) {
				return ConfigErrorNew(PERMISSION_DENIED_ERROR)
			}
			return ConfigErrorNew(UNKNOWN_ERROR)
		}
	}

	_, err = os.Stat(filepath.Join(homeDir, "/.config/feedr/sources.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			sourcesFile, err = os.Create(filepath.Join(homeDir, "/.config/feedr/sources.yml"))
			if err != nil && !os.IsExist(err) {
				if os.IsPermission(err) {
					return ConfigErrorNew(PERMISSION_DENIED_ERROR)
				}
				return ConfigErrorNew(UNKNOWN_ERROR)
			}
		}
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	} else {
		sourcesFile, err = os.Open(filepath.Join(homeDir, "/.config/feedr/sources.yml"))
		if err != nil && !os.IsExist(err) {
			if os.IsPermission(err) {
				return ConfigErrorNew(PERMISSION_DENIED_ERROR)
			}
			return ConfigErrorNew(UNKNOWN_ERROR)
		}
	}

	return nil
}

// Close config.yml and sources.yml
func CloseConfigFiles() {
	configFile.Close()
	sourcesFile.Close()
}

// Load sources from sources.yml
func ReadSources() ([]Source, error) {
	data, err := io.ReadAll(sourcesFile)
	if err != nil {
		return nil, ConfigErrorNew(READ_ERROR)
	}

	var sources []Source
	err = yaml.Unmarshal(data, &sources)
	if err != nil {
		return nil, ConfigErrorNew(YAML_FORMAT_ERROR)
	}

	fmt.Println(sources)

	return sources, nil
}

// Write to sources.yml file
func WriteSources(sources []Source) error {
	data, err := yaml.Marshal(sources)
	if err != nil {
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	err = sourcesFile.Truncate(0)
	if err != nil {
		fmt.Println(err)
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}
	_, err = sourcesFile.Seek(0, 0)
	if err != nil {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	_, err = sourcesFile.Write(data)
	if err != nil {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	return nil
}
