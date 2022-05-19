package slackwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	url string
}

func New(webHookUrl string) *Client {
	c := new(Client)
	c.url = webHookUrl
	return c
}

func (c *Client) SendMessage(message string) error {
	params := make(map[string]string, 1)
	params["text"] = message
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	res, err := sendRequest(http.MethodPost, c.url, body, nil)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	buffer.ReadFrom(res.Body)
	if buffer.String() != "ok" {
		return errors.New("return body is not ok")
	}
	return nil
}

func sendRequest(method, path string, body []byte, params *map[string]string) (*http.Response, error) {
	req, err := newRequest(method, path, body, params)
	if err != nil {
		return nil, err
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		s := buf.String()
		return nil, fmt.Errorf("faild to get data. with error: %s", s)
	}
	return res, nil
}

func newRequest(method, path string, body []byte, params *map[string]string) (*http.Request, error) {
	q := url.Values{}
	if params != nil {
		for k, v := range *params {
			q.Add(k, v)
		}
	}
	req, err := http.NewRequest(method, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
