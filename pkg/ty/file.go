package ty

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const lineBreak = "\n"
const lineRegex = "(.*)[:=](.*)"

/*
* Supported separator: : , = , json map
 */
func (ms *MS) LoadMS(path string) error {

	r, _ := regexp.Compile(lineRegex)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	value, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	strValue := strings.Trim(string(value), " ")

	if strValue[0] == '{' {
		return json.Unmarshal(value, ms)
	}

	lines := strings.Split(strValue, lineBreak)

	for _, v := range lines {
		matches := r.FindAllStringSubmatch(v, len(v))

		(*ms)[strings.Trim(matches[0][1], " ")] = strings.Trim(matches[0][2], " ")
	}

	return nil
}

func (mi *MI) Load(path string) error {
	err := ReadJsonFile(path, mi)
	if err != nil {
		return err
	}

	return nil
}
