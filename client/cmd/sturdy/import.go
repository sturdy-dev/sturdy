package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"getsturdy.com/client/cmd/sturdy/config"
	"getsturdy.com/client/pkg/api"
)

func importCodebase(conf *config.Config, args []string, apiClient *api.HttpApiClient) {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory", err)
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Printf("⚠️  Unexpected number of arguments. Expected: 'sturdy import $CODEBASE_ID'\n")
		os.Exit(1)
	}

	codebaseID := args[0]
	codebase, err := apiClient.GetCodebase(codebaseID)
	if err != nil {
		fmt.Printf("⚠️  Failed to get codebase. Did you enter the correct codebase ID?\n")
		os.Exit(1)
	}

	if !isGitRepo() {
		fmt.Printf("⚠️  The current directory is not within a Git repository. Run this command at the root of a repo to import it.\n")
		os.Exit(1)
	}

	fmt.Printf("✅ Importing git repo at %s to '%s'\n", workingDir, codebase.Name)

	gitProto, gitHost := conf.GetGitRemote()

	cmd := exec.Command("git", "push", fmt.Sprintf("%s://import:%s@%s/%s", gitProto, conf.Auth, gitHost, codebaseID), "HEAD:sturdytrunk", "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Import failed", err)
		fmt.Println(string(output))
		os.Exit(1)
	}

	fmt.Printf("✅ Imported successfully!\n")
}

func isGitRepo() bool {
	// Check if in a git repo
	cmd := exec.Command("git", "show", "HEAD")
	_, err := cmd.CombinedOutput()
	if err == nil {
		return true
	}
	return false
}
