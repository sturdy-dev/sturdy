package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"getsturdy.com/api/pkg/installations/statistics"
)

type Publisher struct {
	httpClient *http.Client
}

func New() *Publisher {
	return &Publisher{
		httpClient: &http.Client{},
	}
}

func (p *Publisher) Publish(ctx context.Context, statistics *statistics.Statistic) error {
	payload, err := json.Marshal(statistics)
	if err != nil {
		return fmt.Errorf("failed to marshal statistics: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.getsturdy.com/v3/statistics", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send statistics: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send statistics: %s (status %d)", string(body), resp.StatusCode)
	}

	return nil
}
