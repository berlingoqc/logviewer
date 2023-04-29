package ty

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type MI map[string]interface{}
type MS map[string]string

const Format = time.RFC3339

func (mi MI) GetString(key string) string {
	if v, b := mi[key]; b {
		return v.(string)
	}
	return ""
}

func (mi MI) GetBool(key string) bool {
	if v, b := mi[key]; b {
		return v.(bool)
	}
	return false
}



func ReadJsonFile(path string, object interface{}) error {
    jsonFile, err := os.Open(path)
    if err != nil { return err }

    defer jsonFile.Close()

    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil { return err }

    return json.Unmarshal([]byte(byteValue), object)
}



