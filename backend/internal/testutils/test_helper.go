package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	FailNow()
	Helper()
}

// AssertTestHelper provides common assertion methods
type AssertTestHelper struct {
	T *testing.T
}

// NewAssertTestHelper creates a new AssertTestHelper
func NewAssertTestHelper(t *testing.T) *AssertTestHelper {
	return &AssertTestHelper{T: t}
}

// Context returns a context for testing
func (h *AssertTestHelper) Context() context.Context {
	return context.Background()
}

// RequireNoError checks if err is nil
func (h *AssertTestHelper) RequireNoError(err error, msgAndArgs ...interface{}) {
	require.NoError(h.T, err, msgAndArgs...)
}

// AssertNoError asserts that err is nil
func (h *AssertTestHelper) AssertNoError(err error, msgAndArgs ...interface{}) bool {
	return assert.NoError(h.T, err, msgAndArgs...)
}

// AssertEqual asserts that two values are equal
func (h *AssertTestHelper) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Equal(h.T, expected, actual, msgAndArgs...)
}

// AssertNotNil asserts that value is not nil
func (h *AssertTestHelper) AssertNotNil(value interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotNil(h.T, value, msgAndArgs...)
}

// AssertError asserts that err is not nil
func (h *AssertTestHelper) AssertError(err error, msgAndArgs ...interface{}) bool {
	return assert.Error(h.T, err, msgAndArgs...)
}
