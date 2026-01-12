package event

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"play-ddd/common"
	"play-ddd/contents/domain"
	"play-ddd/contents/domain/novel"
	dt "play-ddd/contents/infra/repository/pg/datatypes"
	"play-ddd/contents/infra/repository/pg/datatypes/ts"
	dtulid "play-ddd/contents/infra/repository/pg/datatypes/ulid"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

var _ domain.EventRepo[novel.ID, novel.ID] = eventRepo{}

type (
	ID    = dtulid.ULID
	Event struct {
		dt.Model[ID]
		AggregateID   ID
		AggregateKind string
		Kind          string
		Payload       []byte
		Version       uint64
		Status        string
		Reason        string
		Seq           uint64
	}
)

type eventRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) eventRepo { return eventRepo{db} }

// Append implements domain.EventRepo.
func (e eventRepo) Append(
	ctx context.Context, es ...common.Event[novel.ID, novel.ID],
) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		for i := range es {
			var ev Event
			var ver uint64
			var err error
			if err := ev.FromDomain(es[i]); err != nil {
				return fmt.Errorf(`append events: %w`, err)
			}

			err = tx.Model(Event{}).
				WithContext(ctx).
				Select(`COALESCE(MAX(version),0)`).
				Where(`"aggregate_id"=?`, ev.AggregateID).
				Take(&ver).
				Error
			if err != nil {
				return fmt.Errorf(`append events: %w`, err)
			}

			ev.Version = ver + 1

			err = tx.WithContext(ctx).Create(&ev).Error
			if err != nil {
				return fmt.Errorf(`append events: %w`, err)
			}
		}

		return nil
	})
}

// func (e eventRepo) Fetch(
// 	ctx context.Context, aid novel.ID) (
// 	[]common.Event[novel.ID, novel.ID], error,
// ) {
// 	var events []Event
// 	err := e.db.
// 		WithContext(ctx).
// 		Where(`"aggregate_id" = ?`, ID(aid)).
// 		Find(&events).
// 		Error
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// return convert.SliceIntoDomain[Event, novelevent.Event](events)
// 	return nil, nil
// }

// func (e eventRepo) Pull(
// 	ctx context.Context, seq uint64,
// ) ([]common.Event[novel.ID, novel.ID], error) {
// 	return nil, nil
// }

func (e *Event) FromDomain(de common.Event[novel.ID, novel.ID]) (err error) {
	if e == nil {
		return nil
	}

	e.ID = ID(de.ID())
	e.AggregateID = ID(de.AggID())
	e.Kind = de.Kind()
	e.AggregateKind = de.AggKind()
	e.CreatedTs = ts.From(de.EmittedAt())
	e.Payload, err = de.Payload()
	return err
}

func (e *Event) IntoPB() (pe *novelv1.Event, err error) {
	pe = &novelv1.Event{Payload: &anypb.Any{}}
	if err = protojson.Unmarshal(e.Payload, pe.GetPayload()); err != nil {
		return pe, err
	}

	pe.Id = ulidpb.From(e.ID.Into())
	pe.AggregateId = ulidpb.From(e.AggregateID.Into())
	pe.EmitAt = timestamppb.New(e.CreatedTs.Into())
	pe.Kind = e.Kind
	pe.AggregateKind = e.AggregateKind

	return pe, err
}
