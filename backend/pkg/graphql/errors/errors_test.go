package errors

import (
	"database/sql"
	goerr "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errA = goerr.New("a")
var errB = goerr.New("b")

func TestClientSideError(t *testing.T) {
	cases := []struct {
		err error
		is  bool
	}{
		{fmt.Errorf("not found: %w", ErrNotFound), true},
		{fmt.Errorf("fobidden %w", ErrForbidden), true},
		{fmt.Errorf("bad request %w", ErrBadRequest), true},
		{fmt.Errorf("failed to query db %w", sql.ErrNoRows), false},
		{fmt.Errorf("random error"), false},
		{fmt.Errorf("internal error %w", ErrInternalServer), false},
	}
	for _, tc := range cases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			assert.Equal(t, tc.is, IsClientSideError(tc.err))
		})
	}
}

func TestErrorsIs(t *testing.T) {
	cases := []struct {
		err      error
		target   error
		expected bool
	}{
		{err: Error(errA), target: errA, expected: true},
		{err: Error(errA), target: errB, expected: false},
		{err: Error(fmt.Errorf("wrapped A: %w", errA)), target: errA, expected: true},
		{err: Error(fmt.Errorf("wrapped A: %w", errA)), target: ErrInternalServer, expected: true},
		{err: Error(fmt.Errorf("wrapped A: %w", errA)), target: ErrNotFound, expected: false},
		{err: Error(fmt.Errorf("wrapped sql: %w", sql.ErrNoRows)), target: ErrNotFound, expected: true},
		{err: Error(fmt.Errorf("wrapped sql: %w", sql.ErrNoRows)), target: ErrInternalServer, expected: false},
		{err: Error(fmt.Errorf("wrapped A: %w", errA)), target: errB, expected: false},
	}

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			assert.Equal(t, tc.expected, goerr.Is(tc.err, tc.target))
		})
	}
}
