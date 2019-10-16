package util

import (
	"bytes"
	"os"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
)

func SaveJsonMarshalIndent(value interface{}, path string) error {
	marshalJson, err := json.Marshal(value)
	if err != nil {
		return err
	}
	jsonIndent := new(bytes.Buffer)
	err = json.Indent(jsonIndent, marshalJson, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, jsonIndent.Bytes(), 0777)

	if err != nil {
		return err
	}
	return nil
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func Float64ToBytes(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
