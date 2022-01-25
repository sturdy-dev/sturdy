package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Request(host, method, path, authToken string, request, response interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	req, err := http.NewRequest(method, host+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	req.AddCookie(&http.Cookie{Name: "auth", Value: authToken})
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	err = json.Unmarshal(respContent, response)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}
