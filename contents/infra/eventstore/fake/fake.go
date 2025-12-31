package fake

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/puzpuzpuz/xsync/v4"

	"play-ddd/common"
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
type fake[AID, EID comparable] struct {
	m *xsync.Map[AID, []common.Event[AID, EID]]
}

func New[EID, AID comparable]() fake[AID, EID] {
	return fake[AID, EID]{
		m: xsync.NewMap[AID, []common.Event[AID, EID]](),
	}
}

// Append implements domain.EventRepo.
func (f fake[AID, EID]) Append(
	_ context.Context,
	es ...common.Event[AID, EID],
) error {
	for _, e := range es {
		f.m.Compute(e.AggID(), computeEvent(e))
	}

	return nil
}

// Fetch implements domain.EventRepo.
func (f fake[AID, EID]) Fetch(
	_ context.Context,
	k AID,
) ([]common.Event[AID, EID], error) {
	if events, ok := f.m.Load(k); ok {
		return events, nil
	}

	return nil, novel.NewNotfoundError(errors.New(`aggregate not found`))
}

func computeEvent[EID, AID comparable](
	e common.Event[AID, EID],
) func([]common.Event[AID, EID], bool) ([]common.Event[AID, EID], xsync.ComputeOp) {
	return func(
		ov []common.Event[AID, EID], _ bool) (
		_ []common.Event[AID, EID], op xsync.ComputeOp,
	) {
		return append(ov, e), xsync.UpdateOp
	}
}
