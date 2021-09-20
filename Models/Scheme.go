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

// Laps main data of the race
type Lap struct {
	ID             		uint `gorm:"primaryKey" json:"id"`
	OwnerID    		uint `gorm:"primaryKey" json:"owner_id"`
	TagID          		string `gorm:"char(80);index" json:"tag_id" xml:"TagID"`
	DiscoveryUnixTime  	int64 `gorm:"char(80);index" json:"discovery_unix_time"`
	DiscoveryTime  		string `json:"-" xml:"DiscoveryTime"`
	DiscoveryTimePrepared 	time.Time `json:"discovery_time"`
	Antenna        		uint8 `gorm:"index" json:"antenna" xml:"Antenna"`
	AntennaIP      		string `gorm:"char(80);index" json:"antenna_ip"`
	CreatedAt      		time.Time `json:"created_at"`
	RaceID         		uint `gorm:"index" json:"race_id"`
	CurrentRacePosition   	uint `gorm:"index" json:"current_race_postition"`
	TimeBehindTheLeader 	int64 `gorm:"index" json:"time_behind_the_leader"`
	LapNumber      		int `gorm:"index" json:"lap_number"`
	LapTime        		int64 `gorm:"index" json:"lap_time"`
	LapPosition    		uint `gorm:"index" json:"lap_postition"`
	RaceTotalTime           int64 `gorm:"index" json:"race_total_time"`
	BetterOrWorseLapTime	int64 `gorm:"index" json:"better_or_worse_lap_time"`
}

// Laps time labels
type LapSmall struct {
        RaceID         uint `gorm:"index" json:"race_id"`
        LapNumber      uint `gorm:"index" json:"lap_number"`
        DiscoveryTime  string `json:"discovery_time`
        TagID          string `gorm:"char(36);index" json:"tag_id"`
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

// System administators
type Admin struct {
	ID		   uint `gorm:"primaryKey" json:"id"`
	Login	   string `gorm:varchar(100);index" json:"login"`
	Password   string `gorm:varchar(100);index" json:"login"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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

// Systems admin table name
func (u *Admin) TableName() string {
	return "admins"
}
