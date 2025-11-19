package event

import (
	dt "play-ddd/contents/infra/repository/pg/datatypes"
	dtulid "play-ddd/contents/infra/repository/pg/datatypes/ulid"
)

type (
	ID    = dtulid.ULID
	Event struct {
		dt.Model[ID]
		AggregateID   ID     `gorm:"type:char(26);not null;index;"`
		AggregateKind string `gorm:"type:varchar(128);not null;index;comment:Aggregate Kind"`
		Kind          string `gorm:"type:varchar(128);not null;index;comment:Event Kind"`
		Payload       any    `gorm:"type:jsonb;not null;default:'{}';index;serializer:json;comment:Event Payload"`
	}
)
