package client

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/machinebox/graphql"

	"mash/pkg/graphql/model"
	"mash/pkg/version"
)

type Client struct {
	client *graphql.Client
}

// LicenseValidationServer should not be a constant (it can be overwritten by build args)
var LicenseValidationServer = "https://api.getsturdy.com/graphql"

func New() *Client {
	return &Client{
		client: graphql.NewClient(LicenseValidationServer),
	}
}

//go:embed validateLicense.graphql
var validateLicenseQuery string

func (c *Client) Validate(ctx context.Context, licenseKey string, userCount int) (*model.LicenseValidation, error) {
	req := graphql.NewRequest(validateLicenseQuery)
	req.Var("key", licenseKey)
	req.Var("version", version.Version)
	req.Var("bootedAt", int32(version.BootedAt.Unix()))
	req.Var("userCount", userCount)
	req.Var("codebaseCount", 0)

	var res struct {
		ValidateLicense model.LicenseValidation `json:"validateLicense"`
	}

	if err := c.client.Run(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("(clientside) failed to validete license: %w", err)
	}

	return &res.ValidateLicense, nil
}
