package main

import (
	"fmt"
	"log"
	"mash/cmd/sturdy/config"
	"mash/cmd/sturdy/pkg/api"
	"mash/cmd/sturdy/pkg/initView"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// User friendly init
//
// If auth is missing it will run the auth flow.
// If the mount point path is already configured with a view it will not create a new one
// It will restart Sturdy daemon only if needed
func initSmart(conf *config.Config, configPath string, args []string, apiClient api.SturdyAPI) {
	if len(args) < 2 {
		log.Fatalln("❌ Unexpected number of arguments")
	}

	newConfig, _ := createView(conf, configPath, args)

	conf = newConfig

	fmt.Println("🔁 Starting Sturdy")
	startMutagen(configPath, conf, apiClient)

	printReady(args)
}

func printReady(args []string) {
	fmt.Printf("✅ Your new codebase view is now ready at: %s\n", args[1])
}

func createView(conf *config.Config, configPath string, args []string) (newConfig *config.Config, newlyCreated bool) {
	mountPath, err := absPath(args[1])
	if err != nil {
		log.Fatalf("Failed to convert to absolute path: %s\n", err)
	}

	if isInside, gitDirPath := pathIsInGitRepo(mountPath); isInside {
		// This message contains deliberate tabs, to make sure that the message is nicely aligned
		fmt.Printf(`❌	Setting up Sturdy inside a Git checkout is not allowed.
	"%s" is the path to a Git checkout ("%s").
	Re-run 'sturdy init' and replace the last argument ("%s") with a different path
`, mountPath, gitDirPath, args[1])
		os.Exit(1)
	}
	if isInside, viewPath, err := pathIsInsideOtherView(mountPath, conf.Views); isInside {
		// This message contains deliberate tabs, to make sure that the message is nicely aligned
		fmt.Printf(`❌	Setting up Sturdy inside another Sturdy directory is not allowed.
	"%s" is inside another Sturdy directory ("%s").
	Re-run 'sturdy init' and replace the last argument ("%s") with a different path
`, mountPath, viewPath, args[1])
		os.Exit(1)
	} else if err != nil {
		log.Fatalf("Failed to perform directory checks: %s\n", err)
	}

	exists := false
	for _, v := range conf.Views {
		// For backwards compatibility, absPath the left hand side
		p := v.Path
		if existingPath, err := absPath(v.Path); err == nil {
			p = existingPath
		}
		if p == mountPath {
			exists = true
		}
	}

	if exists {
		return conf, false
	}

	codebaseID := args[0]

	viewID, err := initView.CreateWorkspaceAndView(conf.APIRemote, conf.Auth, codebaseID, mountPath)
	if err != nil {
		log.Fatalln(err)
	}

	newConfig, err = config.AddMount(configPath, viewID, mountPath)
	if err != nil {
		log.Fatalln(err)
	}

	return newConfig, true
}

func absPath(mountPath string) (string, error) {
	if path.IsAbs(mountPath) {
		return mountPath, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, mountPath), nil
}

func pathIsInGitRepo(absPath string) (bool, string) {
	// Set absPath to a directory that exists (either absPath, or one of it's parents)
	for {
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			// try again with the parent
			parent := filepath.Dir(absPath)
			if parent == absPath {
				break
			}
			absPath = parent
			continue
		}
		break
	}

	// If this command finishes with exit code 0 it means that this path is inside a git repo
	// and we should not allow setting up Sturdy inside git
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = absPath
	output, err := cmd.CombinedOutput()
	if err == nil {
		gitDirPath := strings.TrimSpace(string(output))
		if !filepath.IsAbs(gitDirPath) {
			gitDirPath = filepath.Join(absPath, gitDirPath)
		}
		return true, filepath.Clean(gitDirPath)
	}
	return false, ""
}

func pathIsInsideOtherView(absPath string, views []config.ViewConfig) (bool, string, error) {
	for _, view := range views {
		if ok, err := pathIsSub(view.Path, absPath); err == nil && ok {
			return true, view.Path, nil
		} else if err != nil {
			return false, "", fmt.Errorf("failed to check if dir is child of: %w", err)
		}
	}
	return false, "", nil
}

// Returns true if (a is a child of b) or (b is a child of a)
func pathIsSub(a, b string) (bool, error) {
	a, _ = absPath(a)
	b, _ = absPath(b)
	directed := func(parent, sub string) (bool, error) {
		up := ".." + string(os.PathSeparator)
		rel, err := filepath.Rel(parent, sub)
		if err != nil {
			return false, nil // silence errors
		}
		if !strings.HasPrefix(rel, up) && rel != ".." {
			return true, nil
		}
		return false, nil
	}

	if ok, err := directed(a, b); err == nil && ok {
		return true, nil
	} else if err != nil {
		return false, err
	}

	if ok, err := directed(b, a); err == nil && ok {
		return true, nil
	} else if err != nil {
		return false, err
	}

	return false, nil
}
