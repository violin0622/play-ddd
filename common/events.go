package common

import (
	"fmt"
	"time"
)

type Event[AID, EID comparable] interface {
	fmt.Stringer

	ID() EID
	Kind() string
	AggID() AID
	AggKind() string
	EmittedAt() time.Time
}

func FormatEvent[AID, EID comparable](e Event[AID, EID]) string {
	return fmt.Sprintf(`%s[%v]: %s[%v] @%s`,
		e.AggKind(), e.AggID(), e.Kind(), e.ID(), e.EmittedAt())
}
