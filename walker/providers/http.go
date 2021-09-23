package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *HTTP) Run(address string) (status int, err error) {
	var data []byte
	if h.JsonBody != nil {
		data, _ = json.Marshal(&h.JsonBody)
		data = []byte(strings.Replace(string(data), "$flag", GenerateFlag(16), 1))
	}
	url := fmt.Sprintf("%s://%s:%d%s", h.Schema, address,h.Port, h.Route)
	req, err := http.NewRequest(
		strings.ToUpper(h.Method),
		url,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return 0, err
	}
	if h.Params != nil {
		q := req.URL.Query()
		for key, value := range h.Params {
			value = strings.Replace(value, "$flag", GenerateFlag(16), 1)
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()

	}
	if h.Header != nil {
		for key, value := range h.Header {
			value = strings.Replace(value, "$flag", GenerateFlag(16), 1)
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return 1, nil
	}
	return 0, fmt.Errorf("%s returned status %d", address, resp.StatusCode)
}
