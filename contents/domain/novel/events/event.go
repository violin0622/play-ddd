package events

import (
	"github.com/oklog/ulid/v2"

	"play-ddd/common"
	"play-ddd/contents/domain"
)

type (
	ID    = ulid.ULID
	Event = domain.Event[ID, ID]
)

var formatEvent = common.FormatEvent[ID, ID]
