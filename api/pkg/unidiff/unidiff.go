package unidiff

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	sdiff "getsturdy.com/api/vcs/diff"

	"github.com/sourcegraph/go-diff/diff"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type FileDiff struct {
	OrigName string `json:"orig_name"`
	NewName  string `json:"new_name"`
	// PreferredName is either OrigName or NewName. Comments and suggestions are reffering to this name.
	PreferredName string `json:"preferred_name"`

	IsDeleted bool `json:"is_deleted"`
	IsNew     bool `json:"is_new"`
	IsMoved   bool `json:"is_moved"`

	IsLarge bool `json:"is_large"`
	// LargeFileInfo is non nil when IsLarge is true
	LargeFileInfo *LargeFileInfo `json:"large_file_info"`

	// If this is true, all other fields are empty.
	// Used to represent file that is hidden from the logged-in user.
	IsHidden bool `json:"is_hidden"`

	Hunks []Hunk `json:"hunks"`
}

type LargeFileInfo struct {
	Size uint64 `json:"size"`
}

type Hunk struct {
	ID string `json:"id"`

	Patch string `json:"patch"`

	// Used only in suggestions
	IsOutdated  bool `json:"is_outdated"`
	IsApplied   bool `json:"is_applied"`
	IsDismissed bool `json:"is_dismissed"`
}

func NewHunk(patch string) Hunk {
	return Hunk{
		// TODO: It would be neat if the ID only took this hunk into account (and not the entire file-diff) that it's a
		// 		part of. Currently the hashes on the git-index row ("index 59e10d8..fc210a8 100644") are based on the
		//		entire version of the file. So if there are three hunks (1, 2, 3) and the contents of hunk 3 changes,
		//		the IDs of all 3 hunks will change, even if the "body" of hunk 1 and 2 is untouched.
		ID:    fmt.Sprintf("%x", sha256.Sum256([]byte(patch))),
		Patch: patch,
	}
}

type PatchReader interface {
	ReadPatch() (string, error)
}

type Unidiff struct {
	diffs             PatchReader
	hunksResolverFunc func(patch string, parsedDiff *diff.FileDiff) ([]Hunk, error)
	ignoreBinary      bool
	invertHunks       bool
	joinHunks         bool
	hunksFilter       map[string]struct{}
	filters           []FilterFunc
	allower           *Allower
	logger            *zap.Logger
}

// FilterFunc is a func that if it returns true, the diff that it's filtering will be removed
type FilterFunc func(diff *diff.FileDiff) (bool, error)

type Option func(*Unidiff)

func NewUnidiff(diffs PatchReader, logger *zap.Logger, options ...Option) *Unidiff {
	u := &Unidiff{
		diffs:             diffs,
		hunksResolverFunc: simpleHunkResolverFunc,
		logger:            logger,
	}
	for _, applyOption := range options {
		applyOption(u)
	}
	return u
}

func WithAllower(allower *Allower) Option {
	return func(unidiff *Unidiff) {
		unidiff.allower = allower
	}
}

func (u *Unidiff) WithAllower(allower *Allower) *Unidiff {
	u.allower = allower
	return u
}

func (u *Unidiff) WithExpandedHunks() *Unidiff {
	u.hunksResolverFunc = expandedHunkResolverFunc
	return u
}

func (u *Unidiff) WithJoiner() *Unidiff {
	u.joinHunks = true
	return u
}

func WithInverter() Option {
	return func(unidiff *Unidiff) {
		unidiff.invertHunks = true
	}
}

func (u *Unidiff) WithInverter() *Unidiff {
	u.invertHunks = true
	return u
}

func (u *Unidiff) WithIgnoreBinary() *Unidiff {
	u.ignoreBinary = true
	return u
}

func WithHunksFilter(id ...string) Option {
	return func(unidiff *Unidiff) {
		filter := make(map[string]struct{})
		for _, i := range id {
			filter[i] = struct{}{}
		}
		unidiff.hunksFilter = filter
	}
}

func (u *Unidiff) WithHunksFilter(id ...string) *Unidiff {
	filter := make(map[string]struct{})
	for _, i := range id {
		filter[i] = struct{}{}
	}
	u.hunksFilter = filter
	return u
}

func (u *Unidiff) WithFilterFunc(fn FilterFunc) *Unidiff {
	u.filters = append(u.filters, fn)
	return u
}

func simpleHunkResolverFunc(patch string, parsedDiff *diff.FileDiff) ([]Hunk, error) {
	d, err := diff.PrintFileDiff(parsedDiff, diff.WithQuotedNames())
	if err != nil {
		return nil, err
	}
	return []Hunk{NewHunk(string(d))}, nil
}

func expandedHunkResolverFunc(patch string, parsedDiff *diff.FileDiff) ([]Hunk, error) {
	expanded := expandHunks(parsedDiff)
	hunks := make([]Hunk, len(expanded), len(expanded))
	for k, ex := range expanded {
		d, err := diff.PrintFileDiff(ex, diff.WithQuotedNames())
		if err != nil {
			return nil, err
		}
		hunks[k] = NewHunk(string(d))
	}
	return hunks, nil
}

func invertHunks(in []Hunk) (out []Hunk, err error) {
	for _, h := range in {
		fd, err := diff.ParseFileDiffKeepCr([]byte(h.Patch))
		if err != nil {
			return nil, err
		}

		inverted, err := invertDiff(fd)
		if err != nil {
			return nil, err
		}

		invertedPatch, err := diff.PrintFileDiff(inverted, diff.WithQuotedNames())
		if err != nil {
			return nil, err
		}

		out = append(out, NewHunk(string(invertedPatch)))
	}
	return
}

// Decorate consumes the input and returns a []FileDiff
func (u *Unidiff) Decorate() ([]FileDiff, error) {
	var res []FileDiff

	for {
		hunks, parsedFileDiff, err := u.readWithHunkResolver()
		if errors.Is(err, io.EOF) {
			break
		} else if errors.Is(err, errHidden) {
			res = append(res, FileDiff{IsHidden: true})
			continue
		} else if errors.Is(err, errEmptyPatch) {
			continue
		} else if errors.Is(err, errParsePatch) {
			// log as error, but continue
			u.logger.Error("failed to parse patch", zap.Error(err))
			continue
		} else if err != nil {
			return nil, fmt.Errorf("could not read with hunk resolver: %w", err)
		}

		fileDiff, err := getFileDiffMeta(parsedFileDiff)
		if err != nil {
			return nil, fmt.Errorf("could not get meta: %w", err)
		}

		fileDiff.Hunks = hunks
		fileDiff.IsLarge, fileDiff.LargeFileInfo = largeData(hunks)
		res = append(res, fileDiff)
	}

	return res, nil
}

// DecorateSeparateBinary consumes the input and returns two separate []FileDiffs
func (u *Unidiff) DecorateSeparateBinary() (binaryDiffs, nonBinaryDiffs []FileDiff, err error) {
	for {
		hunks, parsedFileDiff, err := u.readWithHunkResolver()
		if errors.Is(err, io.EOF) {
			break
		} else if errors.Is(err, errHidden) {
			nonBinaryDiffs = append(nonBinaryDiffs, FileDiff{IsHidden: true})
			continue
		} else if errors.Is(err, errEmptyPatch) {
			continue
		} else if err != nil {
			return nil, nil, err
		}

		fileDiff, err := getFileDiffMeta(parsedFileDiff)
		if err != nil {
			return nil, nil, err
		}
		if parsedFileDiff.Hunks == nil {
			binaryDiffs = append(binaryDiffs, fileDiff)
		} else {
			fileDiff.Hunks = hunks
			nonBinaryDiffs = append(nonBinaryDiffs, fileDiff)
		}

	}
	return
}

func largeData(hunks []Hunk) (bool, *LargeFileInfo) {
	if len(hunks) != 1 {
		return false, nil
	}

	// diff --git a/steam.dmg b/steam.dmg
	// new file mode 100644
	// index 0000000..25b9d04
	// --- /dev/null
	// +++ b/steam.dmg
	// @@ -0,0 +1,3 @@
	// +version https://git-lfs.github.com/spec/v1
	// +oid sha256:dda4744327fe200e08d132ccbba9828b6bde8672080a69f69d52e72e9a6bda17
	// +size 4872474

	// diff --git a/steam.dmg b/steam.dmg
	// index 25b9d04..8d9d584 100644
	// --- a/steam.dmg
	// +++ b/steam.dmg
	// @@ -1,3 +1,3 @@
	// version https://git-lfs.github.com/spec/v1
	// -oid sha256:dda4744327fe200e08d132ccbba9828b6bde8672080a69f69d52e72e9a6bda17
	// -size 4872474
	// +oid sha256:6365d10c9e388ac7a91fe1e65d54694faad69149f421125eaddfff07d48763ea
	// +size 5901865

	if len(hunks[0].Patch) > 800 {
		return false, nil
	}

	rows := strings.Split(hunks[0].Patch, "\n")

	var isLarge bool
	if len(rows) > 5 && strings.Contains(rows[5], "version https://git-lfs.github.com/spec/v1") {
		isLarge = true
	} else if len(rows) > 6 && strings.Contains(rows[6], "version https://git-lfs.github.com/spec/v1") {
		isLarge = true
	}
	if !isLarge {
		return false, nil
	}

	var res LargeFileInfo
	var newSize, oldSize uint64

	for _, r := range rows {
		if strings.HasPrefix(r, "+size ") {
			if s, err := strconv.ParseUint(r[6:], 10, 64); err == nil {
				newSize = s
			}
		}
		if strings.HasPrefix(r, "-size ") {
			if s, err := strconv.ParseUint(r[6:], 10, 64); err == nil {
				oldSize = s
			}
		}
	}

	if newSize > 0 {
		res.Size = newSize
	} else {
		res.Size = oldSize
	}

	return true, &res
}

// DecorateSingle consumes one item from the input, and returns it as a FileDiff
func (u *Unidiff) DecorateSingle() (FileDiff, error) {
	hunks, parsedFileDiff, err := u.readWithHunkResolver()
	if errors.Is(err, errHidden) {
		return FileDiff{IsHidden: true}, nil
	} else if err != nil {
		return FileDiff{}, err
	}
	fileDiff, err := getFileDiffMeta(parsedFileDiff)
	if err != nil {
		return FileDiff{}, err
	}
	fileDiff.Hunks = hunks
	return fileDiff, nil
}

var (
	errEmptyPatch = errors.New("patch is empty")
	errHidden     = errors.New("patch is hidden")
)

// fixLargeFilesDiffs "fixes" the patches for binary files larger than MaxSize
// libgit2 generates invalid diffs for files that are larger than DiffOptions.gitMaxSize
//
// For example the diff contains "Binary files /dev/null and /dev/null differ" for new files, where the second file
// name should be the name of the file.
//
// fixLargeFilesDiffs fixes these diffs and returns a diff that's on a valid/expected format
func fixLargeFilesDiffs(patch string) (string, error) {
	if len(patch) > 2000 || !strings.Contains(patch, "Binary files /dev/null and /dev/null differ") {
		return patch, nil
	}

	expectedDiffPrefix := "diff --git " + sdiff.DiffOldPrefix + "/"

	firstLine := patch[0:strings.Index(patch, "\n")]

	if !strings.HasPrefix(patch, expectedDiffPrefix) || !strings.Contains(patch, sdiff.DiffNewPrefix) {
		return patch, nil
	}

	rows := strings.Split(patch, "\n")
	if len(rows) != 5 {
		return patch, nil
	}

	firstLine = firstLine[len(expectedDiffPrefix):]

	delim := " " + sdiff.DiffNewPrefix + "/"
	idx := strings.Index(firstLine, delim)

	oldName := firstLine[0:idx]
	newName := firstLine[idx+len(delim):]

	// Find mode
	var mode string
	for _, row := range rows {
		if strings.HasPrefix(row, "new mode ") {
			mode = row[len("new mode "):]
		}
	}
	if mode == "" {
		return "", fmt.Errorf("no file mode found")
	}

	return fmt.Sprintf(`diff --git %s/%s %s/%s
new file mode %s
index %s..%s
Binary files %s and %s/%s differ
`,

		sdiff.DiffOldPrefix, oldName,
		sdiff.DiffNewPrefix, newName,
		mode,
		"0000000", "0000000", // TODO?
		"/dev/null", sdiff.DiffNewPrefix, newName,
	), nil

	// Expected
	// diff --git sturdy-old/aaa-100MB.dmg sturdy-new/aaa-100MB.dmg
	// new file mode 100644
	// index 0000000..17b677c
	// Binary files /dev/null and sturdy-new/aaa-100MB.dmg differ

	// Actual
	// diff --git sturdy-old/aaa-100MB.dmg sturdy-new/aaa-100MB.dmg
	// old mode 0
	// new mode 100644
	// Binary files /dev/null and /dev/null differ
}

var errParsePatch = errors.New("could not parse patch")

func (u *Unidiff) readWithHunkResolver() ([]Hunk, *diff.FileDiff, error) {
	patch, err := u.diffs.ReadPatch()
	if errors.Is(err, io.EOF) {
		return nil, nil, io.EOF // propagate
	} else if err != nil {
		return nil, nil, fmt.Errorf("could not read patch: %w", err)
	}

	patch, err = fixLargeFilesDiffs(patch)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fix large files diffs: %w", err)
	}

	// Sometimes we're getting empty-ish diffs from git, ignore them from here
	if strings.Count(patch, "\n") < 2 {
		return nil, nil, errEmptyPatch
	}

	parsedDiff, err := diff.ParseFileDiffKeepCr([]byte(patch))
	if err != nil {

		// Add first line of diff to error message
		var firstLine string
		if idx := strings.IndexByte(patch, '\n'); idx > 0 {
			firstLine = patch[0:idx]
		}

		return nil, nil, multierr.Combine(
			fmt.Errorf("header='%s': %w", firstLine, errParsePatch),
			err)
	}

	if isHidden(u.allower, parsedDiff) {
		return nil, nil, errHidden
	}

	if hasBinaryFiles(parsedDiff) && u.ignoreBinary {
		return nil, nil, errEmptyPatch
	}

	// Apply filters
	for _, fn := range u.filters {
		ok, err := fn(parsedDiff)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return nil, nil, errEmptyPatch
		}
	}

	var hunks []Hunk

	if parsedDiff.Hunks != nil {
		// Ignore diffs with no changes
		st := parsedDiff.Stat()
		if st.Added+st.Changed+st.Deleted == 0 {
			return nil, nil, errEmptyPatch
		}

		// Resolve hunks (non-binary files only)
		hunks, err = u.hunksResolverFunc(patch, parsedDiff)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to resolve hunks: %w", err)
		}
	} else {
		// Binary diffs
		// Use the original patch

		formattedPatch, err := diff.PrintFileDiff(parsedDiff, diff.WithQuotedNames())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build binary patch: %w", err)
		}

		hunks = []Hunk{NewHunk(string(formattedPatch))}
	}

	// Apply filter
	if u.hunksFilter != nil {
		hunks, err = u.filter(hunks)
		if err != nil {
			return nil, nil, err
		}
		// nothing left after filtering
		if hunks == nil || len(hunks) == 0 {
			return nil, nil, errEmptyPatch
		}
	}

	if u.invertHunks {
		hunks, err = invertHunks(hunks)
		if err != nil {
			return nil, nil, err
		}
	}

	if u.joinHunks {
		hunks, err = joinHunks(hunks)
		if err != nil {
			return nil, nil, err
		}
	}

	return hunks, parsedDiff, nil
}

func (u *Unidiff) filter(hunks []Hunk) ([]Hunk, error) {
	var filteredHunks []Hunk
	var droppedAdditions int32

	for _, h := range hunks {
		fd, err := diff.ParseFileDiffKeepCr([]byte(h.Patch))
		if err != nil {
			return nil, err
		}

		// Binary diff that does not pass the filter
		if fd.Hunks == nil {
			if _, ok := u.hunksFilter[h.ID]; !ok {
				return nil, errEmptyPatch
			}
		}

		var firstHunk *diff.Hunk
		if fd.Hunks != nil && len(fd.Hunks) > 0 {
			firstHunk = fd.Hunks[0]
		}

		// This hunk does not match the filter
		if _, ok := u.hunksFilter[h.ID]; !ok {
			if firstHunk != nil {
				droppedAdditions += firstHunk.NewLines - firstHunk.OrigLines
			}
			continue
		}

		// No change to line numbers, we can short-circuit
		if droppedAdditions == 0 {
			filteredHunks = append(filteredHunks, h)
			continue
		}

		// Re-generate the hunk with new line numbers
		if firstHunk != nil {
			if u.invertHunks {
				firstHunk.OrigStartLine = firstHunk.OrigStartLine + droppedAdditions
			} else {
				firstHunk.NewStartLine = firstHunk.NewStartLine - droppedAdditions
			}
		}

		recalculatedHunk, err := diff.PrintFileDiff(fd, diff.WithQuotedNames())
		if err != nil {
			return nil, err
		}
		filteredHunks = append(filteredHunks, NewHunk(string(recalculatedHunk)))
	}

	return filteredHunks, nil
}

func joinHunks(hunks []Hunk) ([]Hunk, error) {
	type namePair [2]string

	diffsByName := make(map[namePair][]*diff.FileDiff)

	for _, hunk := range hunks {
		fd, err := diff.ParseFileDiffKeepCr([]byte(hunk.Patch))
		if err != nil {
			return nil, err
		}
		np := namePair{fd.OrigName, fd.NewName}
		if existing, ok := diffsByName[np]; ok {
			diffsByName[np] = append(existing, fd)
		} else {
			diffsByName[np] = []*diff.FileDiff{fd}
		}
	}

	var newHunks []Hunk
	for _, fds := range diffsByName {
		// use the first one as the base, and add all hunks to it
		joined := fds[0]
		for _, fd := range fds[1:] {
			joined.Hunks = append(joined.Hunks, fd.Hunks...)
		}
		recalculatedHunk, err := diff.PrintFileDiff(joined, diff.WithQuotedNames())
		if err != nil {
			return nil, err
		}
		newHunks = append(newHunks, NewHunk(string(recalculatedHunk)))
	}

	return newHunks, nil
}

// Patches consumes the input, and returns the patches
func (u *Unidiff) Patches() ([]string, error) {
	var res []string

	for {
		hunks, _, err := u.readWithHunkResolver()
		if errors.Is(err, io.EOF) {
			break
		} else if errors.Is(err, errEmptyPatch) {
			continue
		} else if err != nil {
			return nil, err
		}

		for _, hunk := range hunks {
			res = append(res, hunk.Patch)
		}
	}

	return res, nil
}

// PatchesBytes consumes the input, and returns the patches
func (u *Unidiff) PatchesBytes() ([][]byte, error) {
	var res [][]byte

	for {
		hunks, _, err := u.readWithHunkResolver()
		if errors.Is(err, io.EOF) {
			break
		} else if errors.Is(err, errEmptyPatch) {
			continue
		} else if err != nil {
			return nil, err
		}

		for _, hunk := range hunks {
			res = append(res, []byte(hunk.Patch))
		}
	}

	return res, nil
}

func getFileDiffMeta(fd *diff.FileDiff) (FileDiff, error) {
	mode, origName, newName, err := DiffFileStat(fd)
	if err != nil {
		return FileDiff{}, fmt.Errorf("could not get diff stats: %w", err)
	}

	preferredName := newName
	if newName == "/dev/null" {
		preferredName = origName
	}

	return FileDiff{
		OrigName:      origName,
		NewName:       newName,
		PreferredName: preferredName,

		IsDeleted: mode == FileDiffModeDeleted,
		IsNew:     mode == FileDiffModeNew,
		IsMoved:   mode == FileDiffModeMoved,

		// Hunks is not set here
	}, nil
}

func expandHunks(in *diff.FileDiff) []*diff.FileDiff {
	if in == nil {
		return []*diff.FileDiff{}
	}
	if in.Hunks == nil {
		return []*diff.FileDiff{in}
	}
	var out []*diff.FileDiff
	for _, h := range in.Hunks {
		out = append(out, &diff.FileDiff{
			OrigName: in.OrigName,
			OrigTime: in.OrigTime,
			NewName:  in.NewName,
			NewTime:  in.NewTime,
			Extended: in.Extended,
			Hunks:    []*diff.Hunk{h},
		})
	}
	return out
}

func invertDiff(in *diff.FileDiff) (*diff.FileDiff, error) {
	var invertedHunks []*diff.Hunk
	for _, h := range in.Hunks {
		var invertedBody []byte
		s := bufio.NewScanner(bytes.NewReader(h.Body))
		s.Split(bufio.ScanRunes)
		var prev string = "\n"
		for s.Scan() {
			if prev == "\n" && s.Text() == "-" {
				invertedBody = append(invertedBody, []byte("+")...)
			} else if prev == "\n" && s.Text() == "+" {
				invertedBody = append(invertedBody, []byte("-")...)
			} else {
				invertedBody = append(invertedBody, s.Bytes()...)
			}
			prev = s.Text()
		}
		inv := diff.Hunk{
			OrigStartLine:   h.NewStartLine,
			OrigLines:       h.NewLines,
			OrigNoNewlineAt: 0,
			NewStartLine:    h.OrigStartLine,
			NewLines:        h.OrigLines,
			Section:         h.Section,
			StartPosition:   h.StartPosition,
			Body:            invertedBody,
		}
		invertedHunks = append(invertedHunks, &inv)
	}

	invertedExtended, err := invertExtended(in)
	if err != nil {
		return nil, err
	}

	out := diff.FileDiff{
		OrigName: cleanNameNewPrefix(in.NewName, "a/"),
		OrigTime: in.NewTime,
		NewName:  cleanNameNewPrefix(in.OrigName, "b/"),
		NewTime:  in.OrigTime,
		Extended: invertedExtended,
		Hunks:    invertedHunks,
	}
	return &out, nil
}

func invertExtended(fd *diff.FileDiff) ([]string, error) {
	in := fd.Extended

	hasLineWithPrefix := func(prefix string) (bool, int) {
		for idx, l := range in {
			if strings.HasPrefix(l, prefix) {
				return true, idx
			}
		}
		return false, -1
	}

	isNew, newFileModeIdx := hasLineWithPrefix("new file mode ")
	isDeleted, deletedFileModeIdx := hasLineWithPrefix("deleted file mode ")
	isRenamed, _ := hasLineWithPrefix("rename from ")
	hasSimilarityIndex, similarityIndexIdx := hasLineWithPrefix("similarity index")

	hasIndex, indexRowIdx := hasLineWithPrefix("index ")

	invertedOrigName := fd.NewName
	invertedNewName := fd.OrigName

	// go-diff parses the file name of the deleted side as "/dev/null"
	// even if the original diff has named both sides.
	// Restore the file name here.
	if isNew {
		invertedNewName = fd.NewName
	}
	if isDeleted {
		invertedOrigName = fd.OrigName
	}

	var out []string

	// Swapping the file order of the "diff --git" row
	out = append(out, fmt.Sprintf("diff --git %s %s", cleanNameNewPrefix(invertedOrigName, "a/"), cleanNameNewPrefix(invertedNewName, "b/")))

	// Keep similarity Index
	if hasSimilarityIndex {
		out = append(out, in[similarityIndexIdx])
	}

	// Swap "new file mode" with "deleted file mode"
	if isNew {
		out = append(out, "deleted"+in[newFileModeIdx][3:])
	} else if isDeleted {
		out = append(out, "new"+in[deletedFileModeIdx][7:])
	} else if isRenamed {
		// Swap rename from / rename to
		_, renameFromIdx := hasLineWithPrefix("rename from ")
		_, renameToIdx := hasLineWithPrefix("rename to ")
		out = append(out, "rename from "+in[renameToIdx][10:])
		out = append(out, "rename to "+in[renameFromIdx][12:])
	}

	// Swap indexes
	if hasIndex {
		indexTokens := strings.Split(in[indexRowIdx], " ")
		hashTokens := strings.Split(indexTokens[1], "..")

		// Swap order of indexes
		indexTokens[1] = strings.Join([]string{hashTokens[1], hashTokens[0]}, "..")

		out = append(out, strings.Join(indexTokens, " "))
	}

	return out, nil
}

func isHidden(allower *Allower, fd *diff.FileDiff) bool {
	if allower == nil {
		return false
	}
	if fd.NewName == "/dev/null" {
		return !allower.IsAllowed(fd.OrigName, false)
	}
	return !allower.IsAllowed(fd.NewName, false)
}

func hasBinaryFiles(fd *diff.FileDiff) bool {
	if fd == nil {
		return false
	}
	for _, e := range fd.Extended {
		if strings.Contains(e, "Binary files") {
			return true
		}
	}
	return false
}

type FileDiffMode int

var (
	FileDiffModeNew     FileDiffMode = 1
	FileDiffModeChanged FileDiffMode = 2
	FileDiffModeDeleted FileDiffMode = 3
	FileDiffModeMoved   FileDiffMode = 4
)

func cleanName(n string) string {
	idx := strings.IndexByte(n, '/')
	if idx < 0 {
		return n
	}
	return n[idx+1:]
}

func cleanNameNewPrefix(name, prefix string) string {
	if name == "/dev/null" {
		return name
	}
	return prefix + name[2:]
}

// TODO: Can we make this private?
func DiffFileStat(fd *diff.FileDiff) (mode FileDiffMode, origName, newName string, err error) {
	if fd.OrigName == "/dev/null" {
		return FileDiffModeNew, "/dev/null", cleanName(fd.NewName), nil
	} else if fd.NewName == "/dev/null" {
		return FileDiffModeDeleted, cleanName(fd.OrigName), "/dev/null", nil
	} else if fd.NewName != "/dev/null" && fd.OrigName != "/dev/null" && cleanName(fd.NewName) != cleanName(fd.OrigName) {
		return FileDiffModeMoved, cleanName(fd.OrigName), cleanName(fd.NewName), nil
	} else if fd.NewName != "/dev/null" && fd.OrigName != "/dev/null" && cleanName(fd.NewName) == cleanName(fd.OrigName) {
		return FileDiffModeChanged, cleanName(fd.NewName), cleanName(fd.NewName), nil
	} else {
		return 0, "", "", fmt.Errorf("could not detect diff file stat")
	}
}
