package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name         string        `json:"name" gorm:"unique"`
	Capacity     int           `json:"capacity"`
	Type         string        `json:"type"`
	Reservations []Reservation `json:"reservations" gorm:"foreignKey:RoomID"`
}
