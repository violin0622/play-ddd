package events

import (
	"fmt"

	"play-ddd/faculty/domain"

	"github.com/oklog/ulid/v2"
)

type Event = domain.Event[ulid.ULID, ulid.ULID]

func formatEvent(e Event) string {
	return fmt.Sprintf(`%s[%s]: %s[%s] @%s`,
		e.AggKind(), e.AggID(), e.Kind(), e.ID(), e.EmittedAt())
}
