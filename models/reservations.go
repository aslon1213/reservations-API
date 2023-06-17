package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	RoomID     uint      `json:"room_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Reservator string    `json:"reservator"`
}

type Resident struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}
