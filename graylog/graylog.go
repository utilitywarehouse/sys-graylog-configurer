package graylog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	username   string
	password   string
	url        string
	httpClient *http.Client
}

func NewClient(url, password string) *Client {
	cl := &Client{
		url:        url,
		username:   "admin",
		password:   password,
		httpClient: http.DefaultClient,
	}
	return cl
}

func (cl *Client) callAPI(method, path string, input, output interface{}) error {
	// Prepare request
	path = fmt.Sprintf("%s%s", cl.url, path)
	reqBody := &bytes.Buffer{}
	if input != nil {
		err := json.NewEncoder(reqBody).Encode(input)
		if err != nil {
			return fmt.Errorf("error encoding input into json: %s", err)
		}
	}
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return fmt.Errorf("error creating the request: %s", err)
	}
	req.SetBasicAuth(cl.username, cl.password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-By", "graylog-configurer")

	// Execute request
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing the request: %s", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error calling API: %s on %s %s", resp.Status, method, path)
	}
	if output == nil {
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		return fmt.Errorf("error decoding body: %s", err)
	}
	return nil
}

func (cl *Client) ApiReachable() error {
	return cl.callAPI("GET", "/users", nil, nil)
}
