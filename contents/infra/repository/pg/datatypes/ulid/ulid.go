package uliddb

import (
	"database/sql/driver"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ULID ulid.ULID

func New() ULID { return ULID(ulid.Make()) }

func (ULID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Name() {
	case `mysql`, `postgres`, `sqlserver`:
		return `CHAR(26)`
	case `sqlite`:
		return `TEXT`
	default:
		return ``
	}
}

func (u *ULID) Scan(v any) error            { return (*ulid.ULID)(u).Scan(v) }
func (u ULID) Value() (driver.Value, error) { return u.String(), nil }
func (u ULID) String() string               { return ulid.ULID(u).String() }

func (u ULID) Equals(other ULID) bool {
	return ulid.ULID(u).Compare(ulid.ULID(other)) == 0
}

func (u ULID) IsEmpty() bool   { return ulid.ULID(u).IsZero() }
func (u ULID) Into() ulid.ULID { return ulid.ULID(u) }
