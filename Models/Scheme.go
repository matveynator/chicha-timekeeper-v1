package Models

import (
	"time"
	"gorm.io/gorm"
)

// Database locator
var DB *gorm.DB

// Save laps time labels
type Lap struct {
	ID          uint `gorm:"primaryKey"`
	TagID       string `gorm:"char(36)" json:"tag_id"`
	UnixTime    string `gorm:"char(20)" json:"unix_time"`
	Antenna     uint8 `gorm:"index" json:"antenna"`
	CreatedAt   time.Time `json:"created_at"`
}

// Laps time labels table name
func (u *Lap) TableName() string {
	return "laps"
}
