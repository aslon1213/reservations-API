package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aslon1213/reservations-API/models"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) BookRoom(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Invalid room id"})
		return
	}
	room_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check if room exists
	room := models.Room{}
	res := h.DB.First(&room, room_id)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": res.Error.Error()})
		return
	}
	var reserv struct {
		Resident struct {
			Name string `json:"name"`
		} `json:"resident"`
		Start string `json:"start"`
		End   string `json:"end"`
	}

	err = c.BindJSON(&reserv)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	loc, _ := time.LoadLocation("Asia/Tashkent")
	// parse date
	starting_time, err := time.ParseInLocation("02-01-2006 15:04:05", reserv.Start, loc)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ending_time, err := time.ParseInLocation("02-01-2006 15:04:05", reserv.End, loc)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(ending_time.Day())

	fmt.Println(starting_time, ending_time)
	reservations := []models.Reservation{}
	h.DB.Where("room_id = ?", room_id).Where("start >= ?", starting_time).Where("end <= ?", ending_time).Find(&reservations)
	if len(reservations) > 0 {
		c.JSON(400, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
		return
	}

	// check if resident exists
	// resident := models.Resident{}
	// res = h.DB.Where("name = ?", reserv.Resident.Name).First(&resident)
	// if res.Error != nil {
	// 	// create resident
	// 	t := h.DB.Create(&models.Resident{
	// 		Name: reserv.Resident.Name,
	// 	})
	// 	if t.Error != nil {
	// 		c.JSON(400, gin.H{"error": t.Error.Error()})
	// 		return
	// 	}
	// }

	rese := models.Reservation{
		Reservator: reserv.Resident.Name,
		RoomID:     uint(room_id),
		Start:      starting_time,
		End:        ending_time,
	}

	t := h.DB.Create(&rese)
	if t.Error != nil {
		c.JSON(400, gin.H{"error": t.Error.Error()})
		return
	}

	// room.Reservations = append(room.Reservations,
	// h.DB.Commit()

	c.JSON(200, gin.H{"message": "xona muvaffaqiyatli band qilindi"})

}

func (h *Handlers) GetAllReservations(c *gin.Context) {

	reservations := []models.Reservation{}
	h.DB.Find(&reservations)
	c.JSON(200, reservations)

}

func (h *Handlers) DeleteAllReservations(c *gin.Context) {

	reservations := []models.Reservation{}
	h.DB.Find(&reservations)
	h.DB.Delete(&reservations)
	c.JSON(200, gin.H{"message": "All reservations deleted"})

}
