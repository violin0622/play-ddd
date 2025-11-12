package novel

import (
	"errors"
	"fmt"

	"play-ddd/common"
)

var newInvariantsBrokenError = common.NewInvariantsBrokenError

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
	ErrMutateCompletedNovel = errors.New(
		`mutate a completed novel is not allowed`,
	)
	ErrTitleAlreadyExist = errors.New(
		`the author has already created a novel of same title before`,
	)
)
