package mutagen

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Status() ([]SessionStatus, error) {
	output, err := RunMutagenCommandWithRestart("sync", "list", "--json")
	if err != nil {
		return nil, err
	}

	var res []SessionStatus

	// This happens if there is no status
	if string(output) == "null" {
		return res, nil
	}

	// This is a nasty hack. Remove this from the mutagen output instead
	// Mutagen outputs: A lot of control characters (expected for interactive screens, followed by)
	// "Started Mutagen daemon in background (terminate with "mutagen daemon stop")\n"
	var idx int // in this scope so that we can log it
	if len(output) > 0 && output[0] != '[' {
		// Find first [
		idx = bytes.IndexByte(output, '[')
		if idx > -1 {
			output = output[idx:]
		} else {
			return res, nil
		}
	}

	err = json.Unmarshal(output, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from sturdy-sync: idx=%d: %w", idx, err)
	}

	return res, nil
}

type SessionStatus struct {
	Session struct {
		Identifier   string `json:"identifier"`
		Version      int    `json:"version"`
		CreationTime struct {
			Seconds int `json:"seconds"`
			Nanos   int `json:"nanos"`
		} `json:"creationTime"`
		CreatingVersionMinor int `json:"creatingVersionMinor"`
		Alpha                struct {
			Path string `json:"path"`
		} `json:"alpha"`
		Beta struct {
			Protocol int    `json:"protocol"`
			User     string `json:"user"`
			Host     string `json:"host"`
			Path     string `json:"path"`
		} `json:"beta"`
		Configuration struct {
			Ignores           []string `json:"ignores"`
			IgnoreVcsMode     int      `json:"ignoreVCSMode"`
			SshPrivateKeyPath string   `json:"sshPrivateKeyPath"`
		} `json:"configuration"`
		ConfigurationAlpha struct{}          `json:"configurationAlpha"`
		ConfigurationBeta  struct{}          `json:"configurationBeta"`
		Name               string            `json:"name"`
		Paused             bool              `json:"paused"`
		Labels             map[string]string `json:"labels"`
	} `json:"session"`
	Status                          int    `json:"status"`
	AlphaConnected                  bool   `json:"alphaConnected"`
	BetaConnected                   bool   `json:"betaConnected"`
	SuccessfulSynchronizationCycles int    `json:"successfulSynchronizationCycles"`
	LastError                       string `json:"lastError"`
}
