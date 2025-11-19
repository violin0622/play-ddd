package pg

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"play-ddd/contents/app"
	"play-ddd/contents/infra/repository/pg/novel"
)

const DSN = `host=127.0.0.1 port=15432 sslmode=disable user=postgres password=mysecretpassword dbname=playddd`

func InitDB(dsn string) (db *gorm.DB, err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		})
	db, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: newLogger,
		})
	if err != nil {
		return db, err
	}

	return db, novel.Init(db)
}

var _ app.Repo = postgresRepo{}

func New(db *gorm.DB) postgresRepo {
	return postgresRepo{
		db:    db,
		novel: novel.New(db),
	}
}

type postgresRepo struct {
	novel   app.NovelRepo
	chapter app.ChapterRepo

	db *gorm.DB
}

func (g postgresRepo) Chapter() app.ChapterRepo { return g.chapter }
func (g postgresRepo) Novel() app.NovelRepo     { return g.novel }
func (g postgresRepo) Tx(fn func(app.Repo) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn(New(tx))
	})
}
