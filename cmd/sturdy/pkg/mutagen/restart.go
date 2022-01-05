package mutagen

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func RestartDaemon() error {
	stopCtx, stopCancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer stopCancelFunc()
	stopOutput, err := exec.CommandContext(stopCtx, "sturdy-sync", "daemon", "stop").CombinedOutput()

	if errors.Is(stopCtx.Err(), context.DeadlineExceeded) {
		log.Println("Timeout exceeded, trying to restart...")
	} else if err != nil {
		log.Println(string(stopOutput))
		return fmt.Errorf("failed to restart daemon: %w", err)
	}

	startCtx, startCancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer startCancelFunc()
	startOutput, err := exec.CommandContext(startCtx, "sturdy-sync", "daemon", "start").CombinedOutput()

	if errors.Is(startCtx.Err(), context.DeadlineExceeded) {
		log.Println("Timeout exceeded, was not able to start the daemon")
		return context.DeadlineExceeded
	} else if err != nil {
		log.Println(string(startOutput))
		return fmt.Errorf("failed to restart daemon: %w", err)
	}

	return nil
}
