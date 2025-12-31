package outbox

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

// EventRepo is interface between outbox processor and event out persisitent
// storage.
// (e.g. Postgres or MySQL)
// It is assumed that events are arranged in an orderly manner by time overall,
// and individual aggregation events are strictly sorted in the order they
// occur.
type EventRepo interface {
	// Process starts a transaction, executes an arbitrary number of commands
	// in fn within the transaction, commit the changes if no error, otherwise
	// rollback.
	Process(ctx context.Context, fn func(EventsBatch) error) error
}

// EventsBatch fetch and processes a batch of events, record process results by
// AdvanceCursor. PollEvents and AdvanceCursor may be called arbitrary times.
// The amount of Events and Result in each call are not guaranteed to be the
// same, but the order must be the same.
type EventsBatch interface {
	// Poll events limited by arg, in order. Events don't have to be in same
	// aggregation.
	PollEvents(Arg) ([]Event, error)

	//AdvanceCursor pass results, expected them to be persisitent.
	AdvanceCursor(...Result) error
}

type Arg struct {
	Max int
}

type Result struct {
	ID     ID
	Status status
	Reason string
}

type (
	ID    = ulid.ULID
	Event interface {
		ID() ID
		AggregateID() ID
		Kind() string
		AggregateKind() string
		EmitAt() time.Time
		Payload() []byte
	}
)
