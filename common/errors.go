package common

import (
	"errors"
	"fmt"
)

var (
	ErrMutateID         = errors.New(`id already exist`)
	ErrZeroID           = errors.New(`id is zero`)
	ErrInitialEvent     = errors.New(`invalid initial event`)
	ErrUnknownEventKind = func(k string) error { return fmt.Errorf(`unknown event kind: %s`, k) }
)

func NewInvariantsBrokenError(e error) error {
	return InvariantsBrokenError{e: e}
}

type InvariantsBrokenError struct {
	e error
}

func (e InvariantsBrokenError) Error() string {
	return fmt.Sprintf(`invariants broken: %s`, e.e)
}

func (e InvariantsBrokenError) Unwrap() error {
	return e.e
}

// WrapOnErr wrap extra message when err is not nil.
// Deprecated. Use xerr.WrapOn instead.
func WrapOnErr(err *error, msg string, a ...any) error {
	if err == nil || *err == nil {
		return nil
	}

	return fmt.Errorf(`%s: %w`, fmt.Sprintf(msg, a...), *err)
}

type TransientError interface {
	IsTransient(error) bool
}
