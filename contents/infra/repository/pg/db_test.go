package pg_test

import (
	"testing"

	"play-ddd/contents/infra/repository/pg"
)

func TestInitDB(t *testing.T) {
	_, err := pg.InitDB(pg.DSN)
	if err != nil {
		t.Error(err)
	}
}
