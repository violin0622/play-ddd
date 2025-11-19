package datatypes

type Model[A comparable] struct {
	ID        A         `gorm:"primarykey"`
	CreatedAt sqlTime   `gorm:"type:bigint;not null;comment:autofilled create milliseconds unix timestamp."`
	UpdatedAt sqlTime   `gorm:"type:bigint;not null;comment:autofilled update milliseconds unix timestamp."`
	DeletedAt deletedAt `gorm:"type:bigint;not null;index;softDelete:milli;index;comment:autofilled delete milliseconds unix timestamp. Used as soft deletion."`
}
