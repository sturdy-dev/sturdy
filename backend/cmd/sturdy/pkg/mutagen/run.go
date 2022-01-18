package mutagen

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func RunMutagenCommandWithRestart(args ...string) ([]byte, error) {
	// Execute with a context
	firstExecCtx, firstExecCancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer firstExecCancelFunc()

	firstOutput, err := exec.CommandContext(firstExecCtx, "sturdy-sync", args...).CombinedOutput()

	if errors.Is(firstExecCtx.Err(), context.DeadlineExceeded) {
		log.Println("Command was too slow, trying again...")
	} else if err == nil {
		// Everything worked on the first try
		return firstOutput, nil
	}

	// Restart the daemon
	err = RestartDaemon()
	if err != nil {
		return nil, fmt.Errorf("daemon restart failed: %w", err)
	}

	// Try again
	secondExecCtx, secondExecCancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer secondExecCancelFunc()
	secondOutput, err := exec.CommandContext(secondExecCtx, "sturdy-sync", args...).CombinedOutput()

	if errors.Is(secondExecCtx.Err(), context.DeadlineExceeded) {
		log.Println("Command was too slow again, giving up...")
		return nil, secondExecCtx.Err()
	} else if err != nil {
		log.Println(string(firstOutput))
		log.Println(string(secondOutput))
		return nil, fmt.Errorf("failed to run command after restart: %w", err)
	}

	return secondOutput, nil
}
