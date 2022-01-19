package main

import (
	"fmt"
	"os"
)

func isDirectory(d string) error {
	dStat, err := os.Stat(d)

	if os.IsNotExist(err) {
		return fmt.Errorf("missing packageRoot %s: %w", d, err)
	}

	if !dStat.IsDir() {
		return fmt.Errorf("packageRoot (%s) isn't a directory: %w", d, err)
	}

	return nil
}
