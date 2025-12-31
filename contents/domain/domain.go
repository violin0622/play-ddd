package domain

import (
	"context"

	"play-ddd/common"
)

type EventRepo[EID, AID comparable] interface {
	// Fetch(context.Context, AID) ([]Event[EID, AID], error)
	Append(context.Context, ...common.Event[EID, AID]) error
}
