package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/berlingoqc/logviewer/pkg/ty"
)

type JsonPostClient struct {
	client http.Client
	url    string
}

func (c JsonPostClient) Post(path string, headers ty.MS, body interface{}, responseData interface{}) error {

	var buf bytes.Buffer
	encErr := json.NewEncoder(&buf).Encode(body)
	if encErr != nil {
		return encErr
	}

	path = c.url + path

	log.Printf("[POST]%s %s"+ty.LB, path, buf.String())

	req, err := http.NewRequest("POST", path, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		log.Printf("error %d  %s"+ty.LB, res.StatusCode, string(resBody))
		return errors.New(string(resBody))
	}

	return json.Unmarshal(resBody, &responseData)
}

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

	path = c.url + path

	log.Printf("[GET]%s %s\n", path, buf.String())

	req, err := http.NewRequest("GET", path, &buf)
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

func PostClient(url string) JsonPostClient {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	spaceClient := http.Client{Transport: customTransport}

	return JsonPostClient{
		client: spaceClient,
		url:    url,
	}
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
