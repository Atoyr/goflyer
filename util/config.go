package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func CreateConfigDirectoryIfNotExists(appName string) (string, error) {
	var configDir string
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("APPDATA"), appName)
	} else {
		configDir = filepath.Join(home, "."+appName)
	}
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0774); err != nil {
			return "", err
		}
		return configDir, nil
	}
	return configDir, nil
}

func SaveConfig(appName, fileName string, data interface{}) error {
	marshalData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	configDir, err := CreateConfigDirectoryIfNotExists(appName)
	if err != nil {
		return err
	}

	configDir = filepath.Join(configDir, fileName)
	file, err := os.Create(configDir)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(marshalData)
	return nil
}

func LoadConfig(appName, fileName string) (interface{}, error) {
	configDir, err := CreateConfigDirectoryIfNotExists(appName)
	if err != nil {
		return nil, err
	}

	configDir = filepath.Join(configDir, fileName)
	data, err := ioutil.ReadFile(configDir)
	if err != nil {
		return nil, err
	}
	return data, nil
}
