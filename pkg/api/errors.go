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

	// ErrInvalidDates is used when fromDate and toDate is invalid
	ErrInvalidDates = errors.New("from-date must be before the to-date")

	// ErrInvalidImportance is when the importance for an event is invalid
	ErrInvalidImportance = errors.New("importance must be a number between 1 - 5")

	// ErrInvalidEntityType is an error occuring when trying to use a non existing entity-type
	ErrInvalidEntityType = errors.New("invalid entity-type - list all available entity-types with: EntityService.Types")
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
