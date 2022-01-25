package api

import (
	"fmt"

	"getsturdy.com/client/cmd/sturdy/config"
)

type HttpApiClient struct {
	host      string
	authToken string
}

func NewHttpApiClient(c *config.Config) *HttpApiClient {
	return &HttpApiClient{
		host:      c.APIRemote,
		authToken: c.Auth,
	}
}

var _ SturdyAPI = (*HttpApiClient)(nil)

func (h *HttpApiClient) GetView(id string) (View, error) {
	var view View

	err := Request(
		h.host,
		"GET", fmt.Sprintf("/v3/mutagen/get-view/%s", id),
		h.authToken,
		nil,
		&view,
	)
	if err != nil {
		return View{}, fmt.Errorf("failed to load view info: %w", err)
	}

	return view, nil
}

func (h *HttpApiClient) GetCodebase(id string) (Codebase, error) {
	var codebase Codebase

	err := Request(
		h.host,
		"GET", fmt.Sprintf("/v3/codebases/%s", id),
		h.authToken,
		nil,
		&codebase,
	)
	if err != nil {
		return Codebase{}, fmt.Errorf("failed to load codebase info: %w", err)
	}
	return codebase, nil
}

func (h *HttpApiClient) AddPublicKey(publicKey string) error {
	type addPublicKeyRequest struct {
		PublicKey string `json:"public_key"`
	}

	req := addPublicKeyRequest{
		PublicKey: publicKey,
	}

	var res struct{}

	err := Request(h.host, "POST", "/v3/pki/add-public-key", h.authToken, req, &res)
	if err != nil {
		return fmt.Errorf("failed to add public key: %w", err)
	}

	return nil
}

type RenewAuthResponse struct {
	Token  string `json:"token"`
	HasNew bool   `json:"has_new"`
}

func (h *HttpApiClient) RenewAuth() (RenewAuthResponse, error) {
	var res RenewAuthResponse

	err := Request(
		h.host,
		"POST", fmt.Sprintf("/v3/auth/renew-token"),
		h.authToken,
		nil,
		&res,
	)
	if err != nil {
		return RenewAuthResponse{}, fmt.Errorf("failed to renew auth: %w", err)
	}
	return res, nil
}

type GetUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *HttpApiClient) GetUser() (GetUserResponse, error) {
	var res GetUserResponse

	err := Request(
		h.host,
		"GET", fmt.Sprintf("/v3/user"),
		h.authToken,
		nil,
		&res,
	)
	if err != nil {
		return GetUserResponse{}, fmt.Errorf("failed to load user: %w", err)
	}
	return res, nil
}

type GetIgnoresResponse struct {
	Paths []string `json:"paths"`
}

func (h *HttpApiClient) GetIgnores(viewID string) (GetIgnoresResponse, error) {
	var res GetIgnoresResponse

	err := Request(
		h.host,
		"GET", fmt.Sprintf("/v3/views/%s/ignores", viewID),
		h.authToken,
		nil,
		&res,
	)
	if err != nil {
		return GetIgnoresResponse{}, fmt.Errorf("failed to get ignores: %w", err)
	}
	return res, nil
}
