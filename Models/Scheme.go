package Models

import (
	"time"
	"gorm.io/gorm"
)

// Database locator in memory (GORM is calling by Models.DB)
var DB *gorm.DB

// Races details
type Race struct {
    ID              uint `gorm:"primaryKey" json:"id"`
    Name            string `gorm:"varchar(255)" json:"name" form:"name" binding:"required"`
    Description     string `gorm:"text" json:"description" form:"description"`
    IsActive        bool   `gorm:"not null;default:false" json:"is_active"`
    ActualStart     time.Time `json:"actual_start_time"`
    ActualFinish    time.Time `json:"actual_finish_time"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// Laps time labels
type Lap struct {
    ID             uint `gorm:"primaryKey" json:"id"`
	TagID          string `gorm:"char(36);index" json:"tag_id" xml:"TagID"`
	DiscoveryTime  string `json:"-" xml:"DiscoveryTime"`
    DiscoveryTimePrepared time.Time `json:"discovery_time"`
	Antenna        uint8 `gorm:"index" json:"antenna" xml:"Antenna"`
    CreatedAt      time.Time `json:"created_at"`
}

// Users
type User struct {
    ID          uint `gorm:"primaryKey" json:"id"`
    FirstName   string `gorm:"varchar(255);index" json:"first_name" form:"first_name"`
    LastName    string `gorm:"varchar(255)" json:"last_name" form:"last_name"`
    MiddleName  string `gorm:"varchar(255)" json:"middle_name" form:"middle_name"`
    DateOfBirth string `gorm:"varchar(30)" json:"date_of_birth" form:"date_of_birth"`
    City        string `gorm:"varchar(100)" json:"city" form:"city"`
    Team        string `gorm:"varchar(100)" json:"team" form:"team"`
    Class       string `gorm:"varchar(100)" json:"class" form:"class"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Check-in (registration users on race)
type Checkin struct {
    ID          uint `gorm:"primaryKey" json:"id"`
    Number      string `gorm:"varchar(30)" json:"number" form:"number"` // Bib number of user
    TagID       string `gorm:"char(36);index" json:"tag_id" form:"tag_id"`
    UserId      uint `gorm:"index"`
    User        User `gorm:"foreignKey:UserId" json:"user"`
    RaceId      uint `gorm:"index"`
    Race        Race `gorm:"foreignKey:RaceId" json:"race"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Races table name
func (u *Race) TableName() string {
	return "races"
}

// Laps time labels table name
func (u *Lap) TableName() string {
	return "laps"
}

// Users table name
func (u *User) TableName() string {
	return "users"
}

// Laps time labels table name
func (u *Checkin) TableName() string {
	return "checkins"
}
