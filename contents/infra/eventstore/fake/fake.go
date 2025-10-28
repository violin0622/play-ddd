package fake

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/puzpuzpuz/xsync/v4"

	"play-ddd/contents/domain"
	"play-ddd/contents/domain/novel"
)

var (
	_ domain.EventRepo[ulid.ULID, ulid.ULID] = (*fake[ulid.ULID, ulid.ULID])(
		nil,
	)
	_ domain.EventRepo[int64, int64]   = (*fake[int64, int64])(nil)
	_ domain.EventRepo[string, string] = (*fake[string, string])(nil)
	_ novel.EventRepo                  = (*fake[novel.ID, novel.ID])(nil)
)

// fake is a in memory event store for testing.
type fake[EID, AID comparable] struct {
	m *xsync.Map[AID, []domain.Event[EID, AID]]
}

func New[EID, AID comparable]() fake[EID, AID] {
	return fake[EID, AID]{
		m: xsync.NewMap[AID, []domain.Event[EID, AID]](),
	}
}

// Append implements domain.EventRepo.
func (f fake[EID, AID]) Append(
	_ context.Context,
	es ...domain.Event[EID, AID],
) error {
	for _, e := range es {
		f.m.Compute(e.AggID(), computeEvent(e))
	}

	return nil
}

// Fetch implements domain.EventRepo.
func (f fake[EID, AID]) Fetch(
	_ context.Context,
	k AID,
) ([]domain.Event[EID, AID], error) {
	if events, ok := f.m.Load(k); ok {
		return events, nil
	}

	return nil, novel.NewNotfoundError(errors.New(`aggregate not found`))
}

func computeEvent[EID, AID comparable](
	e domain.Event[EID, AID],
) func([]domain.Event[EID, AID], bool) ([]domain.Event[EID, AID], xsync.ComputeOp) {
	return func(
		ov []domain.Event[EID, AID], _ bool) (
		_ []domain.Event[EID, AID], op xsync.ComputeOp,
	) {
		return append(ov, e), xsync.UpdateOp
	}
}
