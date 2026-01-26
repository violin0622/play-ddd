// package stderr provides a event bus that logs all
// consumed events to stderr in JSON, one per line.
package stderr

import (
	"context"
	"fmt"
	"slices"

	"google.golang.org/protobuf/encoding/protojson"

	"play-ddd/contents/infra/outbox"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	"play-ddd/utils/xslice"
)

var _ outbox.EventBus = bus{}

type bus struct{}

func New() bus {
	return bus{}
}

// Pub implements outbox.EventBus.
func (b bus) Pub(_ context.Context, es []*novelv1.Event) error {
	xslice.Foreach(slices.Values(es), func(e *novelv1.Event) {
		raw, err := protojson.Marshal(e)
		if err != nil {
			fmt.Println(`ERROR: `, err)
		}
		fmt.Println(string(raw))
	})

	return nil
}
