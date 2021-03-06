package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"getsturdy.com/api/pkg/licenses"
)

type Validator struct {
	httpClient *http.Client
}

func New() *Validator {
	return &Validator{
		httpClient: &http.Client{},
	}
}

func (v *Validator) Validate(ctx context.Context, licenseKey string) (*licenses.License, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.getsturdy.com/v3/licenses/%s", licenseKey), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to validate license key: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	license := &licenses.License{}
	if err := json.Unmarshal(body, &license); err != nil {
		return nil, fmt.Errorf("failed to decode license: %w", err)
	}

	return license, nil
}
