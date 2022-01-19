package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"client/cmd/sturdy/version"
)

var ErrUnauthorized = errors.New("unexpected response code 401")

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
	req.Header.Set("X-Client-Name", "sturdy-cli")
	req.Header.Set("X-Client-Version", version.Version)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return ErrUnauthorized
	}
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
