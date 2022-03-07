package vcs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func (r *repository) Path() string {
	defer getMeterFunc("Path")()
	return r.path
}

func (r *repository) LargeFilesClean(codebaseID string, paths []string) ([][]byte, error) {
	defer getMeterFunc("LargeFilesClean")()
	if r.lfsHostname == "" {
		return nil, fmt.Errorf("LFS not configured")
	}

	var objectIds []string
	var patches [][]byte

	for _, p := range paths {
		fPath := path.Join(r.path, p)
		// Open file
		fp, err := os.OpenFile(fPath, os.O_RDONLY, 0o644)
		if err != nil {
			return nil, err
		}

		// Clean
		cmd := exec.Command("git-lfs", "clean")
		cmd.Stdin = fp
		cmd.Dir = r.path

		ptr, err := cmd.CombinedOutput()
		if err != nil {
			fp.Close()
			return nil, fmt.Errorf("failed to run git-lfs clean: %w", err)
		}
		fp.Close()

		ptrRows := bytes.Split(ptr, []byte{'\n'})
		if bytes.HasPrefix(ptrRows[1], []byte("oid sha256:")) {
			objectIds = append(objectIds, string(ptrRows[1][len("oid sha256:"):]))
		} else {
			return nil, fmt.Errorf("could not find oid")
		}

		diff, err := buildLfsDiff(p, string(ptr))
		if err != nil {
			return nil, fmt.Errorf("failed to build lfs diff: %w", err)
		}

		patches = append(patches, []byte(diff))
	}

	err := r.configureLfs(codebaseID)
	if err != nil {
		return nil, err
	}

	// Upload
	args := []string{
		"push", "--object-id",
		fmt.Sprintf("http://%s/", r.lfsHostname),
	}
	args = append(args, objectIds...)
	pushCmd := exec.Command("git-lfs", args...)
	pushCmd.Dir = r.path
	output, err := pushCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git-lfs push failed ('%s'): %w", string(output), err)
	}

	return patches, nil
}

func buildLfsDiff(p string, ptr string) (string, error) {
	rows := strings.Split(strings.TrimSpace(ptr), "\n")
	for k, v := range rows {
		rows[k] = "+" + v
	}
	// File names are quoted to support names with spaces on Linux
	return fmt.Sprintf(`diff --git "a/%s" "b/%s"
new file mode 100644
--- /dev/null
+++ "b/%s"
@@ -0,0 +1,3 @@
%s%s`, p, p, p, strings.Join(rows, "\n"), "\n\000No newline at end of file\n"), nil
}

func (r *repository) LargeFilesPull() error {
	defer getMeterFunc("LargeFilesPull")()
	if r.lfsHostname == "" {
		return fmt.Errorf("LFS not configured")
	}

	if err := r.configureLfs(r.CodebaseID()); err != nil {
		return err
	}

	cmd := exec.Command("git-lfs", "pull")
	cmd.Dir = r.path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git-lfs pull failed ('%s'): %w", string(output), err)
	}

	if err := r.CleanStaged(); err != nil {
		return fmt.Errorf("failed to clean index after lfs pull: %w", err)
	}

	return nil
}

func (r *repository) configureLfs(codebaseID string) error {
	defer getMeterFunc("configureLfs")()
	configCmd := exec.Command("git", "config", "lfs.url", fmt.Sprintf("http://%s/api/sturdy/%s", r.lfsHostname, codebaseID))
	configCmd.Dir = r.path
	err := configCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to configure lfs.url: %w", err)
	}
	return nil
}
