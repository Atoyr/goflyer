package util

import (
	"encoding/json"
	"io/ioutil"
	"bytes"
)

func JsonMarshalIndent(value interface{}, path string) error {
marshalJson ,err := json.Marshal(value)
	if err != nil {
		return err
	}
	jsonIndent := new (bytes.Buffer)
	err =json.Indent(jsonIndent, marshalJson, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path,jsonIndent.Bytes(), 0777)

	if err != nil {
		return err
	}
	return nil
}
