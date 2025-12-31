package outbox

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "play-ddd/proto/gen/go/contents/novel/v1"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
	"play-ddd/utils/xerr"
	"play-ddd/utils/xslice"
)

type Relay struct {
	eb          EventBus
	er          EventRepo
	t           *time.Ticker
	maxPub      int
	maxFetch    int
	interval    time.Duration
	log         logr.Logger
	tickTimeout time.Duration
	notify      <-chan struct{}
	wg          sync.WaitGroup

	// These fields are for start/stop control.
	ctx      context.Context
	cancelFn context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewRelay(
	eb EventBus,
	er EventRepo,
	log logr.Logger,
	notify <-chan struct{},
	interval time.Duration,
) *Relay {
	r := &Relay{
		eb:          eb,
		er:          er,
		t:           time.NewTicker(interval),
		maxPub:      100,
		maxFetch:    100,
		interval:    interval,
		notify:      notify,
		log:         log,
		tickTimeout: 10 * time.Second,
	}
	r.stopped.Store(true)
	return r
}

func (r *Relay) Start() (err error) {
	// already started, noop.
	if !r.stopped.Swap(false) {
		return nil
	}

	defer xerr.Expect(&err, `start relay`)

	r.t.Reset(r.interval)
	r.ctx, r.cancelFn = context.WithCancelCause(context.Background())
	r.wg.Go(r.run)

	return nil
}

func (r *Relay) Stop() (err error) {
	// already stopped, noop.
	if r.stopped.Swap(true) {
		return nil
	}

	defer xerr.Expect(&err, `stop relay`)
	r.t.Stop()
	r.cancelFn(fmt.Errorf(`relay stopped`))
	r.wg.Wait()

	return nil
}

func (r *Relay) run() {
	for {
		select {
		case <-r.notify:
			if err := r.tick(r.ctx); err != nil {
				r.log.Error(err, `Relay tick failed.`)
			}
			r.log.V(1).Info(`Relay ticked.`)
		case <-r.t.C:
			if err := r.tick(r.ctx); err != nil {
				r.log.Error(err, `Relay tick failed.`)
			}
			r.log.V(1).Info(`Relay ticked.`)
		case <-r.ctx.Done():
			r.log.Info(`Relay tick stopped.`,
				`cause`, context.Cause(r.ctx))
			return
		}
	}
}

func (r *Relay) tick(ctx context.Context) (err error) {
	defer xerr.Expect(&err, `tick once`)

	ctx, cancel := context.WithTimeoutCause(
		ctx,
		r.tickTimeout,
		fmt.Errorf(`max duration per tick is %s`, r.tickTimeout))
	defer cancel()

	return r.er.Process(r.processEvents(ctx))
}

func (r *Relay) processEvents(ctx context.Context) (
	context.Context, func(ef EventsBatch) error,
) {
	return ctx, func(ef EventsBatch) error {
		events, err := ef.PollEvents(Arg{Max: r.maxFetch})
		if err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		pbEvents, results := convertEvents(events)

		n, err := r.pubEvents(ctx, ef, results, pbEvents)
		if err != nil {
			r.log.Error(err, `Events processing aborted.`,
				`total`, len(events),
				`processed`, n)
		}

		return nil
	}
}

func convertEvents(events []Event) (
	pbEvents []*novelv1.Event, results []Result,
) {
	results = make([]Result, len(events))
	pbEvents = make([]*novelv1.Event, 0, len(events))

	for i := range events {
		pe, err := fromRepo(events[i])
		if err == nil {
			pbEvents = append(pbEvents, pe)
		} else {
			results[i].ID = events[i].ID()
			results[i].Status = failed
			results[i].Reason = err.Error()
		}
	}

	return pbEvents, results
}

type status string

const (
	pending   status = `pending`
	failed    status = `failed`
	published status = `published`
)

func (r *Relay) pubEvents(
	ctx context.Context,
	ef EventsBatch,
	results []Result,
	events []*novelv1.Event) (
	n int, err error,
) {
	defer xerr.Expect(&err, `publish events`)

	left, right := 0, min(r.maxPub, len(events))
	for left < len(events) {
		select {
		case <-ctx.Done():
			return left, context.Cause(ctx)
		default:
		}

		err := r.eb.Pub(ctx, events[left:right])
		xslice.Parallel(
			results[left:right],
			events[left:right],
			func(i int) {
				results[left+i].ID = events[left+i].GetId().Into()
				results[left+i].Status = published
				if err != nil {
					results[left+i].Status = failed
					results[left+i].Reason = err.Error()
				}
			})

		if err = ef.AdvanceCursor(results[left:right]...); err != nil {
			return left, err
		}

		left, right = right, min(right+r.maxPub, len(events))
	}

	return left, nil
}

type (
	EventBus interface {
		Pub(context.Context, []*novelv1.Event) error
	}

	Stats struct {
		RelayedTotal  uint64
		RelayedEvents map[string]uint64
		FailedTotal   uint64
		RetriedTotal  uint64
	}
)

func fromRepo(e Event) (pe *novelv1.Event, err error) {
	defer xerr.Expect(&err, `from repo`)

	pe = &novelv1.Event{Payload: &anypb.Any{}}
	if err = protojson.Unmarshal(e.Payload(), pe.GetPayload()); err != nil {
		return pe, err
	}

	pe.Id = ulidpb.From(e.ID())
	pe.AggregateId = ulidpb.From(e.AggregateID())
	pe.EmitAt = timestamppb.New(e.EmitAt())
	pe.Kind = e.Kind()
	pe.AggregateKind = e.AggregateKind()

	return pe, err
}
