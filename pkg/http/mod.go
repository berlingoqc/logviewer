package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type JsonGetClient struct {
	client http.Client
	url    string
}

func (c JsonGetClient) Get(path string, queryParams ty.MS, body interface{}, responseData interface{}) error {

	var buf bytes.Buffer
	encErr := json.NewEncoder(&buf).Encode(body)
	if encErr != nil {
		return encErr
	}

    fmt.Println(path)
    fmt.Println(string(buf.Bytes()))

	req, err := http.NewRequest("GET", c.url+path, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, getErr := c.client.Do(req)
	if getErr != nil {
		return getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(resBody, &responseData)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func GetClient(url string) JsonGetClient {

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	spaceClient := http.Client{Transport: customTransport}

	return JsonGetClient{
		client: spaceClient,
		url:    url,
	}
}
