package unidiff

import "io"

type stringsPatchReader struct {
	idx int
	s   []string
}

func NewStringsPatchReader(s []string) PatchReader {
	return &stringsPatchReader{s: s}
}

func (s *stringsPatchReader) ReadPatch() (string, error) {
	if s.idx == len(s.s) {
		return "", io.EOF
	}
	r := s.s[s.idx]
	s.idx++
	return r, nil
}
