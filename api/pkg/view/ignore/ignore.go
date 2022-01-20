package ignore

import (
	"io/fs"
	"io/ioutil"
	"path"
	"strings"

	"github.com/tidwall/match"
)

const gitignore = ".gitignore"

func FindIgnore(root fs.FS) ([]string, error) {
	return findIgnore(root, ".", []string{})
}

// findIgnore traverses dir in the root fs.
// ignored is a list of all already discovered patterns to ignore.
// A updated ignored list is returned
func findIgnore(root fs.FS, dir string, ignored []string) ([]string, error) {
	ignoreFp, err := root.Open(path.Join(dir, gitignore))
	if err == nil {
		ignoreContent, err := ioutil.ReadAll(ignoreFp)
		if err != nil {
			return nil, err
		}

		// Add ignores
		lines := strings.Split(string(ignoreContent), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") || len(line) == 0 {
				continue
			}

			// Add as is if in the root
			if dir == "." {
				ignored = append(ignored, line)
			} else {
				ignored = append(ignored, "/"+path.Join(dir, line))
			}
		}
	}

	finfo, err := fs.ReadDir(root, dir)
	if err != nil {
		return nil, err
	}

filesInDir:
	for _, f := range finfo {
		if !f.IsDir() {
			continue
		}

		dirPath := path.Join(dir, f.Name())

		// If path matches any already ignored files, don't keep nesting
		for _, pattern := range ignored {
			if match.Match(dirPath, pattern) || match.Match("/"+dirPath, pattern) {
				continue filesInDir
			}
		}

		ignored, err = findIgnore(root, dirPath, ignored)
		if err != nil {
			return nil, err
		}
	}

	return ignored, nil
}
