package api

import (
	"errors"
	"fmt"
)

var (
	// ErrNotAllowed is used to return an error when a user
	// isn't allowed to perform a specific operation
	ErrNotAllowed = errors.New("not allowed")

	// ErrCannotPerformOperation is used to return an error
	// when a operation is not possible to perform for some reason
	ErrCannotPerformOperation = errors.New("cannot perform this operation")
)

// Error wraps an error with an internal-error
//
// _example_
// if err := auth.Verify(token); err != nil {
//		return api.Error(err, api.ErrNotAllowed)
// }
func Error(source, internal error) error {
	return fmt.Errorf("%s: %s", internal.Error(), source.Error())
}
