package util

import (
	"bytes"
	"os"
	"math"

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
	bits := math.Float64bits(f)
	bytes := make([]byte,8)
	binary.LittleEndian.PutUint64(bytes,bits)
	return bytes
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	f := math.Float64frombits(bits)
	return f
}
