package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const (
	CONFIG_PATH  = "/.config/feedr/config.yml"
	SOURCES_PATH = "/.config/feedr/sources.yml"

	PERMISSION_DENIED_ERROR = "permission denied while trying to access config file"
	HOME_DIR_ERROR          = "could not retrieve home directory"
	READ_ERROR              = "could not read config file"
	WRITE_ERROR             = "could not write to config file"
	YAML_FORMAT_ERROR       = "could not parse yaml data"
	UNKNOWN_ERROR           = "unknown"
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

// Opens file with specified path for reading and writing.
// If any of the parental directories do not exist, creates them
func OpenConfigFile(path string) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, ConfigErrorNew(HOME_DIR_ERROR)
	}

	path = filepath.Join(homeDir, path)
	appPath, _ := filepath.Split(path)
	configPath, _ := filepath.Split(appPath)

	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ConfigErrorNew(HOME_DIR_ERROR)
		}
		if os.IsPermission(err) {
			return nil, ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return nil, ConfigErrorNew(UNKNOWN_ERROR)
	}

	_, err = os.Stat(appPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(homeDir, "/.config/feedr"), 0666)
	}
	if err != nil {
		if os.IsPermission(err) {
			return nil, ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return nil, ConfigErrorNew(UNKNOWN_ERROR)
	}

	var file *os.File

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0666)
	}

	if err != nil {
		if os.IsPermission(err) {
			return nil, ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return nil, ConfigErrorNew(UNKNOWN_ERROR)
	}

	return file, nil
}

// Reads and unmarshalls the contents of the file into "out" structure
func ReadConfigFile(path string, out any) error {
	file, err := OpenConfigFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	data := make([]byte, st.Size())
	_, err = file.Read(data)
	if err != nil {
		if os.IsPermission(err) {
			return ConfigErrorNew(PERMISSION_DENIED_ERROR)
		}
		return ConfigErrorNew(READ_ERROR)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return ConfigErrorNew(YAML_FORMAT_ERROR)
	}

	return nil
}

// Marshalls and writes the contents of "in" structure into file
func WriteConfigFile(path string, in any) error {
	file, err := OpenConfigFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(in)
	if err != nil {
		return ConfigErrorNew(YAML_FORMAT_ERROR)
	}

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return ConfigErrorNew(WRITE_ERROR)
	}

	err = writer.Flush()
	if err != nil {
		return ConfigErrorNew(UNKNOWN_ERROR)
	}

	return nil
}
