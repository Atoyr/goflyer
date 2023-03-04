package config

import (
  "os"
	"bytes"
  "path/filepath"
  "runtime"
	"encoding/json"
	"io/ioutil"
)

func createConfigDirectoryIfNotExists(appName string) (string, error) {
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

func saveJsonMarshalIndent(value interface{}, path string) error {
	marshalJson, err := json.Marshal(value)
	if err != nil {
		return err
	}
	jsonIndent := new(bytes.Buffer)
	err = json.Indent(jsonIndent, marshalJson, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, jsonIndent.Bytes(), 0770)

	if err != nil {
		return err
	}
	return nil
}

