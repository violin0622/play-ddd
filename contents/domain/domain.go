package domain

import (
	"context"
	"fmt"
	"time"
)

type Event[EID, AID comparable] interface {
	fmt.Stringer

	ID() EID
	Kind() string
	AggID() AID
	AggKind() string
	EmittedAt() time.Time
}

type EventRepo[EID, AID comparable] interface {
	// Fetch(context.Context, AID) ([]Event[EID, AID], error)
	Append(context.Context, ...Event[EID, AID]) error
}

type Aggregate[ID comparable] interface {
	ID() ID
	Kind() string
}

type AggregateRepo[ID comparable, A Aggregate[ID]] interface {
	Get(context.Context, ID) (A, error)
	Save(context.Context, A) error
	Update(context.Context, ID, func(context.Context, *A) error) error
}
