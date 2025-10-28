package teacher

import (
	"errors"
	"fmt"
)

var (
	MutateIDError       = errors.New(`id already exist`)
	ZeroIDError         = errors.New(`id is zero`)
	ErrInitialEvent     = errors.New(`invalid initial event`)
	ErrUnknownEventKind = func(k string) error { return fmt.Errorf(`unknown event kind: %s`, k) }
)
