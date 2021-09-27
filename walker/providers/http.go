package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func (h *HTTP) Run(address string, flag string) (status *http.Response, body []byte, err error) {
	var data []byte
	if h.JsonBody != nil {
		data, _ = json.Marshal(&h.JsonBody)
		data = []byte(strings.Replace(string(data), "$flag", flag, 1))
	}
	url := fmt.Sprintf("%s://%s:%d%s", h.Schema, address,h.Port, h.Route)
	req, err := http.NewRequest(
		strings.ToUpper(h.Method),
		url,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, nil, err
	}
	if h.Params != nil {
		q := req.URL.Query()
		for key, value := range h.Params {
			value = strings.Replace(value, "$flag", flag, 1)
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()

	}
	if h.Header != nil {
		for key, value := range h.Header {
			value = strings.Replace(value, "$flag", flag, 1)
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Println(body)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return resp, body,  nil
	}
	return nil, body, fmt.Errorf("%s returned status %d", address, resp.StatusCode)
}
