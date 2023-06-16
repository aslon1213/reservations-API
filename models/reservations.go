package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	Resident   string    `json:"resident"`
	RoomID     uint      `json:"room_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Reservator Resident  `json:"reservator" gorm:"foreignKey:Resident"`
}

type Resident struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}
