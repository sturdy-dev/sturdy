package validator

import (
	"context"
	"encoding/json"
	"fmt"
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

	license := &licenses.License{}
	if err := json.NewDecoder(resp.Body).Decode(license); err != nil {
		return nil, fmt.Errorf("failed to decode license: %w", err)
	}

	return license, nil
}
