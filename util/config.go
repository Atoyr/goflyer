package util

import (
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
