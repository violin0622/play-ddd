package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	sd "gorm.io/plugin/soft_delete"
)

//nolint:lll
type Model[A comparable] struct {
	ID        A         `gorm:"primarykey"`
	CreatedAt sqlTime   `gorm:"type:bigint;not null;autoCreateTime:milli;comment:autofilled create milliseconds unix timestamp."`
	UpdatedAt sqlTime   `gorm:"type:bigint;not null;autoUpdateTime:milli;comment:autofilled update milliseconds unix timestamp."`
	DeletedAt deletedAt `gorm:"type:bigint;not null;index;softDelete:milli;index;comment:autofilled delete milliseconds unix timestamp. Used as soft deletion."`
}

func NewULID() (m Model[ulid.ULID]) {
	now := time.Now()
	m.ID = ulid.Make()
	m.CreatedAt.FromTime(now)
	m.UpdatedAt.FromTime(now)
	return m
}

// nolint:recvcheck
type sqlTime int64

func (t sqlTime) String() string {
	return t.AsTime().String()
}

func (t sqlTime) AsTime() time.Time {
	return time.UnixMilli(int64(t))
}

func (t *sqlTime) FromTime(tt time.Time) {
	*t = sqlTime(tt.UnixMilli())
}

// nolint:recvcheck
type deletedAt struct {
	sd.DeletedAt
}

// nolint:gosec
func (t deletedAt) Value() (driver.Value, error) {
	return int64(t.DeletedAt), nil
}

type scanError struct {
	v any
}

func (e scanError) Error() string {
	return fmt.Sprintf(`unable to convert value %T to deletedAt`, e.v)
}

// nolint:gosec
func (t *deletedAt) Scan(v any) error {
	if v == nil {
		t.DeletedAt = 0
		return nil
	}

	switch v := v.(type) {
	case int64:
		t.DeletedAt = sd.DeletedAt(v)
		return nil
	case time.Time:
		t.DeletedAt = sd.DeletedAt(v.UnixMilli())
		return nil
	default:
		return scanError{v}
	}
}

func (t *deletedAt) String() string {
	return t.AsTime().String()
}

// nolint:gosec
func (t *deletedAt) AsTime() time.Time {
	return time.UnixMilli(int64(t.DeletedAt))
}

// nolint:gosec
func (t *deletedAt) FromTime(tt time.Time) {
	*t = deletedAt{sd.DeletedAt(tt.UnixMilli())}
}
