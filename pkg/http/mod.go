package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/berlingoqc/logviewer/pkg/ty"
)

type Auth interface {
	Login(req *http.Request) error
}

type CookieAuth struct {
	Cookie string
}

func (c CookieAuth) Login(req *http.Request) error {

	req.Header.Set("Cookie", c.Cookie)

	return nil
}

type HttpClient struct {
	client http.Client
	url    string
}

func (c HttpClient) post(path string, headers ty.MS, buf *bytes.Buffer, responseData interface{}, auth Auth) error {
	path = c.url + path

	log.Printf("[POST]%s %s"+ty.LB, path, buf.String())

	req, err := http.NewRequest("POST", path, buf)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if auth != nil {
		if err = auth.Login(req); err != nil {
			log.Printf("%s", err.Error())
		}
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

	//println(string(resBody))

	if res.StatusCode >= 400 {
		log.Printf("error %d  %s"+ty.LB, res.StatusCode, string(resBody))
		return errors.New(string(resBody))
	}

	return json.Unmarshal(resBody, &responseData)
}

func (c HttpClient) PostData(path string, headers ty.MS, body ty.MS, responseData interface{}, auth Auth) error {

	headers["Content-Type"] = "application/x-www-form-urlencoded"

	data := ""

	for k, v := range body {
		data += k + "=" + v + "&"
	}

	buf := bytes.NewBuffer([]byte(data))

	return c.post(path, headers, buf, responseData, auth)

}

func (c HttpClient) PostJson(path string, headers ty.MS, body interface{}, responseData interface{}, auth Auth) error {

	headers["Content-Type"] = "application/json"

	var buf bytes.Buffer
	encErr := json.NewEncoder(&buf).Encode(body)
	if encErr != nil {
		return encErr
	}

	return c.post(path, headers, &buf, responseData, auth)

}

func (c HttpClient) Get(path string, queryParams ty.MS, body interface{}, responseData interface{}, auth Auth) error {

	var buf bytes.Buffer

	if body != nil {
		encErr := json.NewEncoder(&buf).Encode(body)
		if encErr != nil {
			return encErr
		}

	}
	path = c.url + path

	q := url.Values{}

	for k, v := range queryParams {
		q.Add(k, v)
	}

	queryParamString := q.Encode()

	if queryParamString != "" {
		path += "?" + queryParamString
	}

	log.Printf("[GET]%s %s\n", path, buf.String())

	req, err := http.NewRequest("GET", path, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if auth != nil {
		if err = auth.Login(req); err != nil {
			log.Printf("%s", err.Error())
		}
	}

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

	//println(string(resBody))

	jsonErr := json.Unmarshal(resBody, &responseData)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func GetClient(url string) HttpClient {
	spaceClient := getSpaceClient()

	return HttpClient{
		client: spaceClient,
		url:    url,
	}
}

func getSpaceClient() http.Client {
	switch v := http.DefaultTransport.(type) {
	case (*http.Transport):
		customTransport := v.Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		return http.Client{Transport: customTransport}
	default:
		return http.Client{}

	}

}
