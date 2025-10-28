package novel

import (
	"errors"
	"fmt"
)

type InvariantsBrokenError struct {
	e error
}

func newInvariantsBrokenError(e error) error {
	return InvariantsBrokenError{e: e}
}

func (e InvariantsBrokenError) Error() string {
	return fmt.Sprintf(`invariants broken: %s`, e.e)
}

func (e InvariantsBrokenError) Unwrap() error {
	return e.e
}

type NotfoundError struct {
	e error
}

func (n NotfoundError) Error() string {
	return fmt.Sprintf(`not found: %s`, n.e)
}

func (n NotfoundError) Unwarp() error {
	return n.e
}

func NewNotfoundError(e error) error {
	return NotfoundError{e}
}

var (
	ErrMutateCompletedNovel = errors.New(`mutate a completed novel is not allowed`)
	ErrTitleAlreadyExist    = errors.New(
		`the author has already created a novel of same title before`,
	)
)
