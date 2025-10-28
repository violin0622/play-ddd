// memdb is used for test.
package memdb

import (
	"context"
	"sync"

	"github.com/oklog/ulid/v2"

	teacher "play-ddd/faculty/domain/teacher"
)

var _ teacher.Repo[ulid.ULID, teacher.Record[ulid.ULID]] = (*memRepo[ulid.ULID, teacher.Record[ulid.ULID]])(
	nil,
)

type memRepo[A comparable, B teacher.Record[A]] struct {
	sync.RWMutex

	data map[A]B
}

func NewMemRepo[A comparable, B teacher.Record[A]](records ...B) (m memRepo[A, B]) {
	m.data = map[A]B{}
	for i := range records {
		m.data[records[i].ID()] = records[i]
	}
	return m
}

// Create implements domain.Repo.
func (m *memRepo[A, B]) Create(_ context.Context, record B) error {
	m.Lock()
	defer m.Unlock()

	m.data[record.ID()] = record
	return nil
}

// Delete implements domain.Repo.
func (m *memRepo[A, B]) Delete(context.Context, A) error {
	panic("unimplemented")
}

// Get implements domain.Repo.
func (m *memRepo[A, B]) Get(context.Context, A, ...teacher.GetOption) (B, error) {
	panic("unimplemented")
}

// Update implements domain.Repo.
func (m *memRepo[A, B]) Update(context.Context, B) error {
	panic("unimplemented")
}
