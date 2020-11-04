// Package errors provides helper functions and structures for more natural
// error handling in some situations.
package errors

import (
	"errors"
	"fmt"
)

// Convenience function to call errors.New() from the standard library.
func New(text string) error {
	return errors.New(text)
}

// Newf is a convenient function for creating formatted errors.
func Newf(format string, arguments ...interface{}) error {
	return fmt.Errorf(format, arguments...)
}
