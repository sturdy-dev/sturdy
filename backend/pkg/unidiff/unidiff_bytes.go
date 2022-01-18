package unidiff

import "io"

type bytesPatchReader struct {
	idx int
	s   [][]byte
}

func NewBytesPatchReader(s [][]byte) PatchReader {
	return &bytesPatchReader{s: s}
}

func (s *bytesPatchReader) ReadPatch() (string, error) {
	if s.idx == len(s.s) {
		return "", io.EOF
	}
	r := s.s[s.idx]
	s.idx++
	return string(r), nil
}
