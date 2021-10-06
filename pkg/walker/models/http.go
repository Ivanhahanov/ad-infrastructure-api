package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type HTTP struct {
	Route    string                 `yaml:"route"`
	Schema   string                 `yaml:"schema"`
	Method   string                 `yaml:"method"`
	Port     int                    `yaml:"port"`
	Params   map[string]string      `yaml:"params"`
	Header   map[string]string      `yaml:"header"`
	JsonBody map[string]interface{} `yaml:"json_body"`
}

func (h *HTTP) Run(address string, flag string) (status *http.Response, body []byte, err error) {
	var data []byte
	if h.JsonBody != nil {
		data, _ = json.Marshal(&h.JsonBody)
		data = []byte(strings.Replace(string(data), "$flag", flag, 1))
	}

	method := strings.ToUpper(h.Method)
	url := fmt.Sprintf("%s://%s:%d%s", h.Schema, address,h.Port, h.Route)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {return nil, nil, err}

	if h.Params != nil {
		putFlagInQuery(h.Params, req, flag)}
	if h.Header != nil {
		putFlagInHeaders(h.Header, req, flag)}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {return nil, nil, err}

	body, bodyErr := io.ReadAll(res.Body)
	if bodyErr != nil {log.Println(body)}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, body, fmt.Errorf("%s returned status %d", address, res.StatusCode)
	}

	return res, body,  nil
}

func putFlag(value string, flag string) string {return strings.Replace(value, "$flag", flag, 1)}

func putFlagInQuery(params map[string]string, req *http.Request, flag string)  {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, putFlag(value, flag))
	}
	req.URL.RawQuery = q.Encode()
}

func putFlagInHeaders(header map[string]string, req *http.Request, flag string)  {
	for key, value := range header {
		req.Header.Set(key, putFlag(value, flag))
	}
}