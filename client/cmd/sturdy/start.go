package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"getsturdy.com/client/cmd/sturdy/config"
	"getsturdy.com/client/pkg/api"
	"getsturdy.com/client/pkg/edkey"
	"getsturdy.com/client/pkg/mutagen"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// SessionVersionNumber tracks versions of the sturdy configuration.
//
// Sessions are labeled with this number
// If the existing label is different from this number, the session will be terminated and re-created on the next
// "sturdy start" instead of resumed.
//
// [6-19] (inclusive) - Reserved for the CLI version of Sturdy
// 20 - Sturdy Electron 0.3.0
const SessionVersionNumber = "6"

func viewMutagenName(view config.ViewConfig) string {
	return "view-" + view.ID
}

func startMutagen(dotSturdyConfigPath string, conf *config.Config, apiClient api.SturdyAPI) {
	mutagenAgentDirPath, err := mutagenSturdyAgentDirPath()
	if err != nil {
		log.Fatalf("failed to get config: %s", err)
	}

	user, err := apiClient.GetUser()
	if err != nil {
		log.Fatalf("failed to get user: %s", err)
	}

	authorizedKey, privateKeyPath, err := generateKey(mutagenAgentDirPath, user.ID)
	if err != nil {
		log.Fatalf("failed to generate keypair: %s", err)
	}

	err = apiClient.AddPublicKey(authorizedKey)
	if err != nil {
		log.Fatalf("failed to establish a secure connection: %s", err)
	}

	err = ensureKnownHosts(conf.SyncRemote)
	if err != nil {
		log.Fatalf("failed to add trust: %s", err)
	}

	mutagenSessions, err := mutagen.Status()
	if err != nil {
		log.Fatalf("failed to get current status: %s", err)
	}

	sessionsByName := make(map[string]mutagen.SessionStatus)
	for _, s := range mutagenSessions {
		sessionsByName[s.Session.Name] = s
	}

	// removeViewConfIdx contains keys for view indexes that should be removed from the configuration
	// after this command has been successfully executed
	removeViewConfIdx := make(map[int]struct{})

	removeView := func(name string, viewIDx int, view config.ViewConfig) {
		// Mark for removal
		removeViewConfIdx[viewIDx] = struct{}{}

		// Terminate the connection if it exists
		if _, ok := sessionsByName[name]; ok {
			err = terminateMutagenForView(view)
			if err != nil {
				log.Fatalf("failed to terminate: %s", err)
			}
		}
	}

	if len(conf.Views) == 0 {
		fmt.Println("You don't have any codebases configured. Go to https://getsturdy.com to get started!")
		os.Exit(0)
	}

	// Create mutagen sync sessions
	for viewIDx, view := range conf.Views {
		name := viewMutagenName(view)

		apiView, err := apiClient.GetView(view.ID)
		if errors.Is(err, api.ErrUnauthorized) {
			// Broom emoji
			fmt.Printf("\U0001F9F9 Removing configuration for %s (you don't have access to the codebase)\n", view.Path)
			removeView(name, viewIDx, view)
			continue
		}
		if err != nil {
			log.Fatalf("failed to get view data from Sturdy, skipping: %s", err)
			continue
		}
		if apiView.CodebaseIsArchived {
			// Broom emoji
			fmt.Printf("\U0001F9F9 Removing configuration for %s (the codebase has been archived)\n", view.Path)
			removeView(name, viewIDx, view)
			continue
		}

		// Get ignores paths
		ignores, err := apiClient.GetIgnores(view.ID)
		if err != nil {
			log.Fatalf("failed to get syncing data from Sturdy, skipping: %s", err)
			continue
		}

		viewConfigPath, err := configPathForView(privateKeyPath, mutagenAgentDirPath, view, ignores.Paths)
		if err != nil {
			log.Fatalf("failed to get config: %s", err)
		}

		var didMigrate bool
		if existingSess, ok := sessionsByName[name]; ok {
			// Resume if the session version is the expected number
			// Otherwise terminate followed by start
			if versionLabel, ok := existingSess.Session.Labels["sessionVersion"]; ok && versionLabel == SessionVersionNumber {
				err = resumeMutagenForView(view)
				if err != nil {
					log.Fatalf("failed to resume: %s", err)
				}

				fmt.Printf("✅ Started %s\n", view.Path)
				continue
			}

			// Do terminate
			err = terminateMutagenForView(view)
			if err != nil {
				log.Fatalf("failed to terminate: %s", err)
			}
			didMigrate = true
		}

		remote := conf.APIRemote
		labelProto := ""
		if strings.HasPrefix(remote, "https://") {
			remote = remote[len("https://"):]
			labelProto = "https"
		}
		if strings.HasPrefix(remote, "http://") {
			remote = remote[len("http://"):]
			labelProto = "http"
		}
		remoteParts := strings.Split(remote, ":")
		labelHost := remoteParts[0]
		var labelHostPort string
		if len(remoteParts) > 1 {
			labelHostPort = remoteParts[1]
		}

		args := []string{
			"sync", "create",
			"--no-global-configuration",
			"-c", viewConfigPath,
			"--name", viewMutagenName(view),
			"--label", "sturdy=true",
			"--label", fmt.Sprintf("sessionVersion=%s", SessionVersionNumber),
			"--label", fmt.Sprintf("sturdyApiProto=%s", labelProto),
			"--label", fmt.Sprintf("sturdyApiHost=%s", labelHost),
			"--label", fmt.Sprintf("sturdyApiHostPort=%s", labelHostPort),
			"--label", fmt.Sprintf("sturdyViewId=%s", view.ID),
			"--stage-mode-beta=neighboring",

			// Alpha
			view.Path,

			// Beta
			fmt.Sprintf("%s@%s:/repos/%s/%s/",
				apiView.UserID,
				conf.SyncRemote,
				apiView.CodebaseID,
				view.ID,
			),
		}

		_, err = mutagen.RunMutagenCommandWithRestart(args...)
		if err != nil {
			log.Printf("failed to start sturdy: %s\n", err)
			os.Exit(1)
		}

		if didMigrate {
			fmt.Printf("✅ Started %s! (migrated to latest version)\n", view.Path)
		} else {
			fmt.Printf("✅ Started %s for the first time!\n", view.Path)
		}
	}

	// Filter the config, and remove views that should no longer exist
	var newViews []config.ViewConfig
	for idx, view := range conf.Views {
		if _, ok := removeViewConfIdx[idx]; ok {
			continue
		}
		newViews = append(newViews, view)
	}
	conf.Views = newViews

	err = config.WriteConfig(dotSturdyConfigPath, conf)
	if err != nil {
		fmt.Printf("Failed to save new configuration to %s: %s\n", dotSturdyConfigPath, err)
	}
}

func resumeMutagenForView(view config.ViewConfig) error {
	_, err := mutagen.RunMutagenCommandWithRestart("sync", "resume", viewMutagenName(view))
	if err != nil {
		return err
	}
	return nil
}

func terminateMutagenForView(view config.ViewConfig) error {
	_, err := mutagen.RunMutagenCommandWithRestart("sync", "terminate", viewMutagenName(view))
	if err != nil {
		return err
	}
	return nil
}

// hostWithOptionalPort is on format "host" or "host:1234"
func ensureKnownHosts(hostWithOptionalPort string) error {
	var port = "22"
	var host = hostWithOptionalPort

	parts := strings.Split(hostWithOptionalPort, ":")
	if len(parts) > 1 {
		host = parts[0]
		port = parts[1]
	}

	cmd := exec.Command("ssh-keyscan", "-p", port, host)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("could not get key for sync.getsturdy.com: %w", err)
	}

	var trustRows []string
	for _, sshKeyscanRow := range strings.Split(string(output), "\n") {
		// Ignore comments in the output
		if strings.HasPrefix(sshKeyscanRow, "#") {
			continue
		}

		sshKeyscanRow = strings.TrimSpace(sshKeyscanRow)

		// Collect rows
		trustRows = append(trustRows, sshKeyscanRow)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home dir: %w", err)
	}

	knownHostsDir := filepath.Join(homeDir, ".ssh")
	knownHostsFilePath := filepath.Join(homeDir, ".ssh", "known_hosts")

	// Create the .ssh dir if it does not exist
	// drwxr-xr-x
	err = os.MkdirAll(knownHostsDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to mkdir .ssh (to create known_hosts): %w", err)
	}

	existingKnownHosts, err := ioutil.ReadFile(knownHostsFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not open .ssh/known_hosts: %w", err)
	}

	newKnownHosts := addToKnownHosts(trustRows, existingKnownHosts, lineSeparator)

	err = ioutil.WriteFile(knownHostsFilePath, newKnownHosts, 0o644)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to write new keys to .ssh/known_hosts: %w", err)
	}

	return nil
}

func addToKnownHosts(trustRows []string, existingKnownHosts []byte, lineSeparator string) []byte {
	knownHosts := bytes.TrimSpace(existingKnownHosts)
	if len(knownHosts) > 0 {
		knownHosts = append(knownHosts, []byte(lineSeparator)...)
	}

	// Append trust rows to existing known hosts if it doesn't already exist
	for _, trustRow := range trustRows {
		if bytes.Contains(knownHosts, []byte(trustRow)) {
			continue
		}

		trustRow = strings.TrimSpace(trustRow) + lineSeparator

		// Append
		knownHosts = append(knownHosts, []byte(trustRow)...)

		fmt.Printf("🔐 Adding trust: %s\n", strings.TrimSpace(trustRow))
	}

	return knownHosts
}

type mutagenConfig struct {
	Sync mutagenSyncConfig `json:"sync"`
}

type mutagenSyncConfig struct {
	Defaults mutagenSyncEndpointConfig `json:"defaults"`
}

type mutagenSyncEndpointConfig struct {
	Mode              string                          `json:"mode"`
	SshPrivateKeyPath string                          `json:"sshPrivateKeyPath"`
	Ignore            mutagenSyncEndpointIgnoreConfig `json:"ignore"`
}

type mutagenSyncEndpointIgnoreConfig struct {
	Paths []string `json:"paths"`
	VCS   bool     `json:"vcs"`
}

func configPathForView(privateKeyPath, sturdyAgentDir string, view config.ViewConfig, ignores []string) (string, error) {
	ignores = append(ignores, "node_modules", ".DS_Store", "*.swp")
	conf := mutagenConfig{
		Sync: mutagenSyncConfig{
			Defaults: mutagenSyncEndpointConfig{
				Mode:              "two-way-resolved",
				SshPrivateKeyPath: privateKeyPath,
				Ignore: mutagenSyncEndpointIgnoreConfig{
					Paths: ignores,
					VCS:   true,
				},
			},
		},
	}

	rawConfig, err := json.Marshal(conf)
	if err != nil {
		return "", fmt.Errorf("failed to construct sturdy-sync config: %w", err)
	}

	confPath := filepath.Join(sturdyAgentDir, view.ID+".yaml")

	err = ioutil.WriteFile(
		confPath,
		rawConfig,
		0o644,
	)
	if err != nil {
		return "", fmt.Errorf("failed to save sturdy-sync config: %w", err)
	}

	return confPath, nil
}

func mutagenSturdyAgentDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to lookup home dir: %w", err)
	}

	configPath := filepath.Join(homeDir, ".sturdy-agent")
	err = os.Mkdir(configPath, 0o744)
	if errors.Is(err, os.ErrExist) {
		return configPath, nil
	} else if err != nil {
		return "", fmt.Errorf("could not create .sturdy-agent dir: %w", err)
	}

	return configPath, nil
}

func generateKey(configPath string, userID string) (authorizedKey, privateKeyPath string, err error) {
	privateKeyPath = filepath.Join(configPath, "private-key-ed25519-"+userID+".pem")

	// Check if we have a private key already
	existingPrivateKey, err := os.Open(privateKeyPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", "", fmt.Errorf("could not read private key: %w", err)
	}

	// Parse existing key
	if err == nil {
		allKeyBytes, err := ioutil.ReadAll(existingPrivateKey)
		if err != nil {
			return "", "", fmt.Errorf("failed to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(allKeyBytes)
		if err != nil {
			return "", "", fmt.Errorf("could not parse existing key: %w", err)
		}

		// Use existing keypair
		return string(ssh.MarshalAuthorizedKey(signer.PublicKey())), privateKeyPath, nil
	}

	// Generate a new keypair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("could not generate new keypair: %w", err)
	}

	pub, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("could not parse new public key: %w", err)
	}

	encode, err := edkey.MarshalED25519PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal ed25519 private key")
	}

	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_RDWR|os.O_CREATE, 0o400)
	if err != nil {
		return "", "", fmt.Errorf("could open private key: %w", err)
	}
	defer privateKeyFile.Close()

	block := &pem.Block{Type: "OPENSSH PRIVATE KEY", Bytes: encode}

	if err := pem.Encode(privateKeyFile, block); err != nil {
		return "", "", fmt.Errorf("failed to encode private key: %w", err)
	}

	return fmt.Sprintf("%s", ssh.MarshalAuthorizedKey(pub)), privateKeyPath, nil
}
