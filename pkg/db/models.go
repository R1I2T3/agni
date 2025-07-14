package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string    `gorm:"type:text;uniqueIndex"`
	APIToken      string    `gorm:"type:text;uniqueIndex"`
	APISecret     string    `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Notifications []Notification `gorm:"foreignKey:ApplicationID"`
}

// BeforeCreate hook to set UUID before inserting
func (app *Application) BeforeCreate(tx *gorm.DB) (err error) {
	app.ID = uuid.New()
	return
}

type Notification struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;index"` // FK
	QueueID       string    `gorm:"type:text;uniqueIndex"`
	Type          string    `gorm:"type:text"`
	Recipient     string    `gorm:"type:text"`
	Message       string    `gorm:"type:text"`
	Status        string    `gorm:"type:text"`
	Attempts      int
	CreatedAt     time.Time
	PersistedAt   *time.Time
	ProcessedAt   *time.Time
}

// BeforeCreate hook to set UUID
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.New()
	return
}
