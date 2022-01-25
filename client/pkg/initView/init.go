package initView

import (
	"fmt"
	"os"

	"getsturdy.com/client/pkg/api"
)

func CreateWorkspaceAndView(host, authToken, codebaseID, mountPath string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	type createWorkspaceRequest struct {
		CodebaseID string `json:"codebase_id"`
	}

	type createWorkspaceResponse struct {
		ID string `json:"id"`
	}

	req := createWorkspaceRequest{
		CodebaseID: codebaseID,
	}

	var resp createWorkspaceResponse
	err = api.Request(host, "POST", "/v3/workspaces", authToken, req, &resp)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace: %w", err)
	}

	type createViewRequest struct {
		CodebaseID    string `json:"codebase_id"`
		WorkspaceID   string `json:"workspace_id"`
		MountPath     string `json:"mount_path"`
		MountHostname string `json:"mount_hostname"`
	}

	type createViewResponse struct {
		ID string `json:"id"`
	}

	viewReq := createViewRequest{
		CodebaseID:    codebaseID,
		WorkspaceID:   resp.ID,
		MountPath:     mountPath,
		MountHostname: hostname,
	}

	var viewResp createViewResponse
	err = api.Request(host, "POST", "/v3/views", authToken, viewReq, &viewResp)
	if err != nil {
		return "", fmt.Errorf("failed to create view: %w", err)
	}

	return viewResp.ID, nil
}
