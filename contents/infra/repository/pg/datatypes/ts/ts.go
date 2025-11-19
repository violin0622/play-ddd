package ts

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Timestamp int64

func (Timestamp) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Name() {
	case `mysql`, `postgres`, `sqlserver`:
		return `BIGINT`
	case `sqlite`:
		return `INTEGER`
	default:
		return ``
	}
}

func (ts Timestamp) Into() time.Time { return time.UnixMilli(int64(ts)) }
func From(t time.Time) Timestamp     { return Timestamp(t.UnixMilli()) }

//
// func (u ULID) Value() (driver.Value, error) { return u.String(), nil }
//
// func (u ULID) String() string { return ulid.ULID(u).String() }
//
// func (u ULID) Equals(other ULID) bool {
// 	return ulid.ULID(u).Compare(ulid.ULID(other)) == 0
// }
//
// func (u ULID) IsEmpty() bool { return ulid.ULID(u).IsZero() }
//
// func (u ULID) Into() ulid.ULID { return ulid.ULID(u) }
