package outbox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
	"play-ddd/utils/xerr"
)

// ErrRelayStopped is the error returned when the relay is stopped.
var ErrRelayStopped = errors.New("relay stopped")

// Relay relays events from the event store to the event bus.
type Relay struct {
	bus  EventBus
	repo EventRepo

	// Configuration
	maxFetch    int
	maxPub      int
	interval    time.Duration
	tickTimeout time.Duration
	instance    ulid.ULID

	// Runtime
	ticker *time.Ticker
	notify <-chan struct{}
	log    logr.Logger

	mu       sync.Mutex
	running  bool
	ctx      context.Context
	cancelFn context.CancelCauseFunc
	wg       sync.WaitGroup

	// Health check
	lastTickAt atomic.Value // time.Time
	// lastTickErr atomic.Value // error
}

// NewRelay creates a new Relay with the given dependencies and options.
func NewRelay(
	bus EventBus,
	repo EventRepo,
	log logr.Logger,
	// notify <-chan struct{},
	opts ...Option,
) *Relay {
	r := &Relay{
		bus:         bus,
		repo:        repo,
		maxFetch:    100,
		maxPub:      100,
		interval:    time.Second,
		tickTimeout: 10 * time.Second,
		notify:      make(<-chan struct{}),
		log:         log,
		instance:    ulid.Make(),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.ticker = time.NewTicker(r.interval)
	r.ticker.Stop() // Stop until Start() is called

	return r
}

// Start starts the relay. It is idempotent and concurrent-safe.
func (r *Relay) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return nil
	}

	r.ctx, r.cancelFn = context.WithCancelCause(context.Background())
	r.ticker.Reset(r.interval)
	r.running = true

	r.wg.Go(r.run)

	return nil
}

// Stop stops the relay gracefully. It is idempotent and concurrent-safe.
// It waits for the current batch to complete before returning.
func (r *Relay) Stop() error {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return nil
	}
	r.running = false
	r.mu.Unlock()

	r.ticker.Stop()
	r.cancelFn(ErrRelayStopped)
	r.wg.Wait()

	return nil
}

// Healthy checks if the relay is healthy (last tick was within expected time).
func (r *Relay) Healthy() bool {
	r.mu.Lock()
	running := r.running
	r.mu.Unlock()

	if !running {
		return true // not started is considered healthy
	}

	t, ok := r.lastTickAt.Load().(time.Time)
	if !ok {
		return true // just started, no tick yet
	}

	// Last tick should be within 2x interval
	return time.Since(t) < 2*r.interval
}

// Ready checks if the relay is ready (running and healthy).
func (r *Relay) Ready() bool {
	r.mu.Lock()
	running := r.running
	r.mu.Unlock()

	return running && r.Healthy()
}

// LastError returns the last tick error, if any.
// func (r *Relay) LastError() error {
// 	if err, ok := r.lastTickErr.Load().(error); ok {
// 		return err
// 	}
// 	return nil
// }

// run is the main loop that listens for triggers and processes events.
func (r *Relay) run() {
NEXT:
	select {
	case <-r.ctx.Done():
		r.log.Info("Relay stopped", "cause", context.Cause(r.ctx))
		return
	case _, ok := <-r.notify:
		if !ok {
			// Channel closed, set to nil to prevent busy loop
			r.notify = nil
			goto NEXT
		}
		r.tick()
	case <-r.ticker.C:
		r.tick()
	}
	goto NEXT
}

// tick performs a single tick and updates health check state.
func (r *Relay) tick() {
	r.lastTickAt.Store(time.Now())

	ctx, cancel := context.WithTimeoutCause(
		r.ctx,
		r.tickTimeout,
		fmt.Errorf("max duration per tick is %s", r.tickTimeout))
	defer cancel()

	err := r.repo.Process(r.processEvents(ctx))

	if err != nil {
		r.log.Error(err, "Relay tick failed.")
		return
	}

	r.log.V(1).Info("Relay ticked.")
}

func (r *Relay) processEvents(
	ctx context.Context) (
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

		results := make([]Result, len(events))
		pbEvents := make([]*novelv1.Event, len(events))
		r.convertEvents(events, pbEvents, results)

		n, err := r.pubEvents(ctx, ef, events, pbEvents, results)
		if err != nil {
			r.log.Error(err, "Events processing aborted",
				"total", len(events),
				"processed", n)
		}

		return nil
	}
}

func (r *Relay) convertEvents(
	events []Event,
	pbEvents []*novelv1.Event,
	results []Result,
) {
	for i := range events {
		pe, err := fromRepo(events[i])
		if err == nil {
			pbEvents[i] = pe
		} else {
			results[i].Status = StatusFailed
			results[i].Reason = err.Error()
		}
	}
}

type msg struct {
	e   Evnt
	pbe *novelv1.Event
	r   Result
}

// pubEvents publishes events in batches, maintaining order within aggregates.
// If an event fails, subsequent events in the same aggregate are skipped.
func (r *Relay) pubEvents(
	ctx context.Context,
	ef EventsBatch,
	events []Event,
	pbEvents []*novelv1.Event,
	results []Result,
) (n int, err error) {
	defer xerr.Expect(&err, "publish events")

	// Track failed aggregates to maintain order within aggregate
	failedAggs := make(map[ID]struct{})

	for left := 0; left < len(pbEvents); left += r.maxPub {
		select {
		case <-ctx.Done():
			return left, context.Cause(ctx)
		default:
		}

		right := min(left+r.maxPub, len(pbEvents))
		batch := pbEvents[left:right]
		batchResults := results[left:right]
		batchEvents := events[left:right]

		// Filter out events from failed aggregates
		toSend := make([]*novelv1.Event, 0, len(batch))
		toSendIdx := make([]int, 0, len(batch))

		for j, e := range batchEvents {
			aggID := e.AggregateID()

			// Check if this aggregate has failed
			if _, failed := failedAggs[aggID]; failed {
				// Skip events from failed aggregates, keep as pending for retry
				batchResults[j].ID = e.ID()
				batchResults[j].Status = StatusPending
				batchResults[j].Reason = "skipped: previous event in same aggregate failed"
				continue
			}

			// Skip events that failed conversion
			if batch[j] == nil {
				continue
			}

			toSend = append(toSend, batch[j])
			toSendIdx = append(toSendIdx, j)
		}

		if len(toSend) == 0 {
			// Still need to advance cursor for skipped events
			if err := ef.AdvanceCursor(batchResults...); err != nil {
				return left, err
			}
			continue
		}

		// Publish batch
		pubErr := r.bus.Pub(ctx, toSend)

		// Update results
		for k, idx := range toSendIdx {
			batchResults[idx].ID = toSend[k].GetId().Into()
			if pubErr != nil {
				batchResults[idx].Status = StatusFailed
				batchResults[idx].Reason = pubErr.Error()
				// Mark this aggregate as failed
				failedAggs[batchEvents[idx].AggregateID()] = struct{}{}
			} else {
				batchResults[idx].Status = StatusPublished
			}
		}

		if err := ef.AdvanceCursor(batchResults...); err != nil {
			return left, err
		}
	}

	return len(pbEvents), nil
}

func fromRepo(e Event) (pe *novelv1.Event, err error) {
	defer xerr.Expect(&err, "from repo")

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
