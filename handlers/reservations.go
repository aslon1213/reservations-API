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

	// if starting_time.Hour() == ending_time.Hour() {
	// 	c.JSON(410, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
	// 	return
	// }

	if starting_time.Day() != ending_time.Day() {
		if starting_time.Hour() < ending_time.Hour() {
			c.JSON(400, gin.H{"error": "Session ni oxiri boshidan keyin bo'la olmaydi :_)"})
			return
		}
		c.JSON(400, gin.H{"error": "Kun bir xil bo'lishi kerak"})
		return
	}
	if starting_time.Day() < time.Now().Day() {
		c.JSON(400, gin.H{"error": "Bugundan oldingi sanalarga buyurtma qila olmaysiz"})
		return
	}

	if starting_time.Hour() < 9 {
		starting_time = time.Date(starting_time.Year(), starting_time.Month(), starting_time.Day(), 9, 0, 0, 0, loc)
		if ending_time.Hour() < 9 {
			c.JSON(400, gin.H{"error": "Kun boshidan oldin buyurtma qila olmaysiz"})
			return
		}
	}

	if ending_time.Hour() > 18 {

		ending_time = time.Date(ending_time.Year(), ending_time.Month(), ending_time.Day(), 18, 0, 0, 0, loc)
		if starting_time.Hour() > 18 {
			c.JSON(400, gin.H{"error": "Kun oxiridan keyin buyurtma qila olmaysiz"})
			return
		}
	}

	fmt.Println(starting_time, ending_time)
	reservations := []models.Reservation{}
	h.DB.Where("room_id = ?", room_id).Find(&reservations)

	// check if room is available

	// | xx |
	// x | x |

	for _, reserv := range reservations {
		if starting_time.After(reserv.Start) && starting_time.Before(reserv.End) {
			c.JSON(410, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
			return
		}
		if ending_time.After(reserv.Start) && ending_time.Before(reserv.Start) {
			c.JSON(410, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
			return
		}
		if ending_time.After(reserv.End) && starting_time.Before(reserv.Start) {
			c.JSON(410, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
			return
		}
		if ending_time.After(reserv.End) && ending_time.Before(reserv.Start) {
			c.JSON(410, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
			return
		}
		if starting_time.Equal(reserv.Start) || ending_time.Equal(reserv.End) {
			c.JSON(410, gin.H{"error": "uzr, siz tanlagan vaqtda xona band"})
			return
		}

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

	c.JSON(201, gin.H{"message": "xona muvaffaqiyatli band qilindi"})

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

func (h *Handlers) GetAllReservationsFullInfo(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "Invalid room id"})
		return
	}
	reservations := []models.Reservation{}
	t := h.DB.Find(&reservations, "room_id = ?", id)
	if t.Error != nil {
		c.JSON(400, gin.H{"error": t.Error.Error()})
		return
	}
	c.JSON(200, reservations)

}

func (h *Handlers) UnbookRoom(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "Invalid room id"})
		return
	}

	var reserv struct {
		Resident struct {
			Name string `json:"name"`
		} `json:"resident"`
		Start string `json:"start"`
		End   string `json:"end"`
	}

	err := c.BindJSON(&reserv)
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

	reservation := models.Reservation{}

	fmt.Println(starting_time, ending_time)

	h.DB.Where("start = ?", starting_time).Where("end = ?", ending_time).First(&reservation)
	if reservation.ID == 0 {
		c.JSON(400, gin.H{"error": "Reservation not found"})
		return
	}

	h.DB.Delete(&reservation)
	c.JSON(200, gin.H{"message": "Reservation deleted"})

}
