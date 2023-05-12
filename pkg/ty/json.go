package ty

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const LB = "\n"

func ReadJsonFile(path string, object interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, object)
}

func ToJsonString(data any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJsonString(data string, placeholder any) error {
	return json.Unmarshal([]byte(data), placeholder)
}
