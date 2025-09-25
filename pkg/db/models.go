package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	// Change `type:uuid` to `type:varchar(36)`
	ID            uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Name          string    `gorm:"type:varchar(255);uniqueIndex"`
	APIToken      string    `gorm:"type:varchar(255);uniqueIndex"`
	APISecret     string    `gorm:"type:varchar(255)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Notifications []Notification `gorm:"foreignKey:ApplicationID"`
}

// BeforeCreate hook is correct and needs no changes
func (app *Application) BeforeCreate(tx *gorm.DB) (err error) {
	app.ID = uuid.New()
	return
}

type Notification struct {
	// Change `type:uuid` to `type:varchar(36)`
	ID                 uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	ApplicationID      uuid.UUID `gorm:"type:varchar(36);index"` // FK
	QueueID            string    `gorm:"type:varchar(100);uniqueIndex"`
	Type               string    `gorm:"type:text"`
	Channel            string    `gorm:"type:text"`
	Provider           string    `gorm:"type:text"`
	TemplateID         string    `json:"template_id,omitempty"`
	MessageContentType string    `json:"message_content_type,omitempty"`
	Recipient          string    `gorm:"type:text"`
	Subject            string    `gorm:"type:text"`
	Message            string    `gorm:"type:text"`
	Status             string    `gorm:"type:text"`
	Attempts           int
	// In-app read state (only used for InApp channel)
	Read   bool       `gorm:"default:false;index" json:"read"`
	ReadAt *time.Time `json:"read_at,omitempty"`

	CreatedAt   time.Time
	UpdatedAt   time.Time `json:"updated_at"`
	PersistedAt *time.Time
	ProcessedAt *time.Time
}

// BeforeCreate hook is correct and needs no changes
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.New()
	return
}

type WebPushSubscription struct {
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	UserID    string    `gorm:"index;size:36"`
	Endpoint  string    `gorm:"size:500;uniqueIndex;not null"`
	P256dh    string    `gorm:"size:255;not null"`
	Auth      string    `gorm:"size:255;not null"`
	Device    string    `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (wps *WebPushSubscription) BeforeCreate(tx *gorm.DB) (err error) {
	wps.ID = uuid.New()
	return
}
