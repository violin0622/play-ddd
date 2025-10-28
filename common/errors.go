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
