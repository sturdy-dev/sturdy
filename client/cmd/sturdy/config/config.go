package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Remote         string       `json:"remote"` // gRPC API (unused)
	InsecureRemote bool         `json:"insecure-remote,omitempty"`
	APIRemote      string       `json:"api-remote"`  // HTTP API
	SyncRemote     string       `json:"sync-remote"` // Mutagen SSH API
	Auth           string       `json:"auth"`
	GitRemote      string       `json:"git-remote,omitempty"` // Git Server
	Views          []ViewConfig `json:"views"`
}

func (c Config) GetGitRemote() (proto, host string) {
	if c.GitRemote == "" {
		// Default remote
		return "https", "git.getsturdy.com"
	}

	return "http", c.GitRemote
}

type ViewConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

func ReadConfig(path string) (*Config, error) {
	configContents, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a default config
			// TODO: Easier defaults for local dev
			return &Config{
				Remote:         "fs.getsturdy.com:443",
				InsecureRemote: false,
				APIRemote:      "https://api.getsturdy.com",
				SyncRemote:     "sync.getsturdy.com",
			}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	var conf Config
	err = json.Unmarshal(configContents, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// For backwards compatibility
	if conf.SyncRemote == "" {
		conf.SyncRemote = "sync.getsturdy.com"
	}

	return &conf, nil
}

func WriteConfig(path string, conf *Config) error {
	data, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func SetAuth(configPath, auth string) error {
	c, err := ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	c.Auth = auth
	err = WriteConfig(configPath, c)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	return nil
}

func AddMount(configPath, viewID, mountPath string) (*Config, error) {
	c, err := ReadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	c.Views = append(c.Views, ViewConfig{
		ID:   viewID,
		Path: mountPath,
	})

	err = WriteConfig(configPath, c)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}
	return c, nil
}
