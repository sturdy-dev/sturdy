package unidiff

/*
source: https://github.com/mutagen-io/mutagen/blob/master/pkg/synchronization/core/allow.go
licence: MIT (https://github.com/mutagen-io/mutagen/blob/master/LICENSE)
*/

import (
	"errors"
	"fmt"
	pathpkg "path"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
)

// allowPattern represents a single parsed allow pattern.
type allowPattern struct {
	// negated indicates whether or not the pattern is negated.
	negated bool
	// directoryOnly indicates whether or not the pattern should only match
	// directories.
	directoryOnly bool
	// matchLeaf indicates whether or not the pattern should be matched against
	// a path's base name in addition to the whole path.
	matchLeaf bool
	// pattern is the pattern to use in matching.
	pattern string
}

// newAllowPattern validates and parses a user-provided allow pattern.
func newAllowPattern(pattern string) (*allowPattern, error) {
	// Check for invalid patterns, or at least those that would leave us with an
	// empty string after parsing. Obviously we can't perform general complete
	// validation for all patterns, but if they pass this parsing, they should
	// be sane enough to at least try to match.
	if pattern == "" || pattern == "!" {
		return nil, errors.New("empty pattern")
	} else if pattern == "/" || pattern == "!/" {
		return nil, errors.New("root pattern")
	} else if pattern == "//" || pattern == "!//" {
		return nil, errors.New("root directory pattern")
	}

	// Check if this is a negated pattern. If so, remove the exclamation point
	// prefix, since it won't enter into pattern matching.
	negated := false
	if pattern[0] == '!' {
		negated = true
		pattern = pattern[1:]
	}

	// Check if this is an absolute pattern. If so, remove the forward slash
	// prefix, since it won't enter into pattern matching.
	absolute := false
	if pattern[0] == '/' {
		absolute = true
		pattern = pattern[1:]
	}

	// Check if this is a directory-only pattern. If so, remove the trailing
	// slash, since it won't enter into pattern matching.
	directoryOnly := false
	if pattern[len(pattern)-1] == '/' {
		directoryOnly = true
		pattern = pattern[:len(pattern)-1]
	}

	// Determine whether or not the pattern contains a slash.
	containsSlash := strings.IndexByte(pattern, '/') >= 0

	// Attempt to do a match with the pattern to ensure validity. We have to
	// match against a non-empty path (we choose something simple), otherwise
	// bad pattern errors won't be detected.
	if _, err := doublestar.Match(pattern, "a"); err != nil {
		return nil, fmt.Errorf("unable to validate pattern: %w", err)
	}

	// Success.
	return &allowPattern{
		negated:       negated,
		directoryOnly: directoryOnly,
		matchLeaf:     (!absolute && !containsSlash),
		pattern:       pattern,
	}, nil
}

// matches indicates whether or not the allow pattern matches the specified
// path and metadata.
func (i *allowPattern) matches(path string, directory bool) (bool, bool) {
	// If this pattern only applies to directories and this is not a directory,
	// then this is not a match.
	if i.directoryOnly && !directory {
		return false, false
	}

	// Check if there is a direct match. Since we've already validated the
	// pattern in the constructor, we know match can't fail with an error (it's
	// only return code is on bad patterns).
	if match, _ := doublestar.Match(i.pattern, path); match {
		return true, i.negated
	}

	// If it makes sense, attempt to match on the last component of the path,
	// assuming the path is non-empty (non-root).
	if i.matchLeaf && path != "" {
		if match, _ := doublestar.Match(i.pattern, pathpkg.Base(path)); match {
			return true, i.negated
		}
	}

	// No match.
	return false, false
}

// Allower is a collection of parsed allow patterns.
type Allower struct {
	// patterns are the underlying allow patterns.
	patterns []*allowPattern

	Patterns []string
}

func deduplicate(ss []string) []string {
	noDuplicates := make([]string, 0, len(ss))
	seen := make(map[string]bool, len(ss))
	for _, s := range ss {
		if seen[s] {
			continue
		}
		seen[s] = true
		noDuplicates = append(noDuplicates, s)
	}
	return noDuplicates
}

// NewAllower creates a new allowr given a list of user-provided allow
// patterns.
func NewAllower(patterns ...string) (*Allower, error) {
	patterns = deduplicate(patterns)
	patterns = append(patterns, "!.git", "!.git/**/*")
	// Parse patterns.
	allowPatterns := make([]*allowPattern, len(patterns))
	for i, p := range patterns {
		if ip, err := newAllowPattern(p); err != nil {
			return nil, fmt.Errorf("unable to parse pattern: %w", err)
		} else {
			allowPatterns[i] = ip
		}
	}

	// Success.
	return &Allower{
		patterns: allowPatterns,
		Patterns: patterns,
	}, nil
}

// IsAllowed determines whether or not the specified path should be allowd based
// on all provided allow patterns and their order.
func (i *Allower) IsAllowed(path string, directory bool) bool {
	// Nothing is initially allowed.
	allowed := false

	// Run through patterns, keeping track of the allowd state as we reach more
	// specific rules.
	for _, p := range i.patterns {
		if match, negated := p.matches(path, directory); !match {
			continue
		} else {
			allowed = !negated
		}
	}

	// Done.
	return allowed
}
