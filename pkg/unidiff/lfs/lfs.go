package lfs

import (
	"bytes"
	"io"
	"log"
	"mash/pkg/unidiff"
	"mash/vcs"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sourcegraph/go-diff/diff"
)

var (
	lfsIgnoreSmuggedHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sturdy_git_lfs_smudged_filter_millis",
		Help:    "Duration in milliseconds",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 50, 100, 200, 500, 2000, 5000},
	}, []string{"outcome"})

	lfsCacheCurrentSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sturdy_git_lfs_smudged_cache_current_size",
		Help: "Count of entries in the cache",
	})
)

type cacheKey struct {
	path    string
	size    int64
	modTime time.Time
}

var cache *lru.Cache

// TODO: Dependency inject the cache instead?
func init() {
	var err error
	cache, err = lru.New(2048)
	if err != nil {
		log.Fatal(err)
	}
}

func NewIgnoreLfsSmudgedFilter(repo vcs.RepoReader) (unidiff.FilterFunc, error) {
	headCommit, err := repo.HeadCommit()
	if err != nil {
		return nil, err
	}

	return func(diff *diff.FileDiff) (bool, error) {
		t0 := time.Now()
		var metricOutcome string
		defer func() {
			lfsCacheCurrentSize.Set(float64(cache.Len()))
			lfsIgnoreSmuggedHistogram.WithLabelValues(metricOutcome).Observe(float64(time.Since(t0).Milliseconds()))
		}()

		isBinaryFile := diff.Hunks == nil
		if !isBinaryFile {
			metricOutcome = "non-binary"
			return false, nil
		}

		isNewFile := diff.OrigName == "/dev/null"
		if isNewFile {
			metricOutcome = "new-file"
			return false, nil
		}

		if !strings.HasPrefix(diff.NewName, "b/") {
			metricOutcome = "no-b-slash"
			return false, nil
		}
		// remove "b/" prefix
		name := diff.NewName[2:]

		blob, err := repo.FileBlobAtCommit(headCommit.Id().String(), name)
		if err != nil {
			if err.Error() == "file not found in tree" {
				metricOutcome = "file-not-found"
				return false, nil
			}
			metricOutcome = "get-blob-error"
			return false, err
		}

		// Fast path, this is not a LFS pointer (a LFS pointer is ~130 bytes)
		if blob.Size() < 100 || blob.Size() > 200 {
			metricOutcome = "ignored-size"
			return false, nil
		}

		fileDiskPath := path.Join(repo.Path(), name)

		fstat, err := os.Stat(fileDiskPath)
		if err != nil {
			metricOutcome = "stat-failed"
			return false, err
		}

		key := cacheKey{
			path:    fileDiskPath,
			size:    fstat.Size(),
			modTime: fstat.ModTime(),
		}

		// get from cache
		if val, ok := cache.Get(key); ok {
			if valBool, ok := val.(bool); ok {
				metricOutcome = "cached"
				return valBool, nil
			}
		}

		// exec and set cache
		res, err := shellExecGitLfsPointer(bytes.NewReader(blob.Contents()), fileDiskPath)
		if err != nil {
			// don't add to cache
			metricOutcome = "exec-err"
			return false, err
		}

		cache.Add(key, res)

		if res {
			metricOutcome = "true"
		} else {
			metricOutcome = "exec-one"
		}

		return res, nil
	}, nil
}

func shellExecGitLfsPointer(pointer io.Reader, fileDiskPath string) (bool, error) {
	cmd := exec.Command("git-lfs", "pointer", "--file", fileDiskPath, "--stdin")
	cmd.Stdin = pointer
	err := cmd.Run()

	if err != nil {
		if cmd.ProcessState.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
