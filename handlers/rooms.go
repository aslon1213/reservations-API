package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aslon1213/reservations-API/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHandlers(ctx context.Context, db *gorm.DB) *Handlers {
	return &Handlers{
		ctx: ctx,
		DB:  db,
	}
}

type Handlers struct {
	ctx context.Context
	DB  *gorm.DB
}

var Room_Types = []string{
	"focus",
	// "meeting",
	"conference",
	"team",
}

func (h *Handlers) GetRooms(c *gin.Context) {

	types := c.Query("type")
	page := c.Query("page")
	page_int, _ := strconv.Atoi(page)
	if page_int <= 0 {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}
	page_size := c.Query("page_size")
	page_size_int, err := strconv.Atoi(page_size)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	search := c.Query("search")
	fmt.Println(types, page_int, page_size, search)
	rooms := []models.Room{}
	if types == "" && search == "" {
		h.DB.Find(&rooms)
	} else {
		h.DB.Where("type = ?", types).Where("name LIKE ?", "%"+search+"%").Find(&rooms)

	}
	type output_room struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Capacity int    `json:"capacity"`
	}
	var output []output_room
	for _, room := range rooms {
		a := output_room{
			ID:       room.ID,
			Name:     room.Name,
			Type:     room.Type,
			Capacity: room.Capacity}
		output = append(output, a)
	}

	fmt.Println(output)
	rooms_with_pagination := map[int][]output_room{}
	for i := 0; i <= len(output)/page_size_int; i++ {
		if len(output) < (i+1)*page_size_int {
			rooms_with_pagination[i] = output[i*page_size_int:]
		} else {
			rooms_with_pagination[i] = output[i*page_size_int : (i+1)*page_size_int]
		}

	}
	// c.JSON(200, gin.H{
	// 	"rooms": rooms,
	// })

	c.JSON(200, gin.H{
		"page":      page_int,
		"count":     len(output)/page_size_int + 1,
		"page_size": page_size_int,
		"rooms":     rooms_with_pagination[page_int-1]})

}

func (h *Handlers) GetRoom(c *gin.Context) {

	id := c.Param("id")
	room := models.Room{}
	res := h.DB.First(&room, id)

	reservations := []models.Reservation{}
	h.DB.Where("room_id = ?", id).Find(&reservations)

	if res.Error != nil {
		c.JSON(400, gin.H{"error": "topilmadi"})
		return
	}
	room.Reservations = reservations
	c.JSON(200, gin.H{"id": room.ID,
		"name":     room.Name,
		"type":     room.Type,
		"capacity": room.Capacity})

}

type output struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (h *Handlers) GetRoomAvailability(c *gin.Context) {
	id := c.Param("id")
	room := models.Room{}
	res := h.DB.First(&room, id)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": res.Error.Error() + " Getting room failed"})
		return
	}
	err := error(nil)
	dat := c.Query("date")
	date := time.Time{}
	if dat != "" {
		date, err = time.Parse("02-01-2006", dat)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format please give in the format 02-01-2006 DD-MM-YYYY"})
			return
		}
	}
	if dat == "" {
		date = time.Now()
	}

	// get reservations
	fmt.Println(date)
	// get date to start of day
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	// get date to end of day
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)
	reservations := []models.Reservation{}
	h.DB.Where("room_id = ?", id).Where("start <= ?", end).Where("end >= ?", start).Find(&reservations)
	// h.DB.Where("room_id = ?", id).Find(&reservations)
	available_hours := getavailablehours(reservations)
	fmt.Println(available_hours)

	o := []output{}
	if len(available_hours) == 24 {
		o = append(o, output{Start: 0, End: 24})
	} else {
		prev, cur := 0, 1
		for cur != len(available_hours) {
			if available_hours[cur]-available_hours[cur-1] > 1 {
				o = append(o, output{Start: available_hours[prev], End: available_hours[cur-1] + 1})
				prev = cur
			}
			cur++
		}
		if prev-cur != 0 {
			o = append(o, output{Start: available_hours[prev], End: available_hours[cur-1] + 1})
		}
	}

	type output_ava struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	out := []output_ava{}

	for _, oo := range o {
		start, end, _ := oo.formatavailablehours(date.Day(), int(date.Month()), date.Year())
		out = append(out, output_ava{Start: start, End: end})
	}

	room.Reservations = reservations
	c.JSON(200, out)
}

func (a *output) formatavailablehours(day int, month int, year int) (string, string, error) {

	t1 := time.Date(year, time.Month(month), day, a.Start, 0, 0, 0, time.UTC)
	t2 := time.Date(year, time.Month(month), day, a.End, 0, 0, 0, time.UTC)

	start := t1.Format("02-01-2006 15:04:05")
	end := t2.Format("02-01-2006 15:04:05")
	return start, end, nil
}

func getavailablehours(reservations []models.Reservation) []int {
	available_hours := []int{}
	for i := 0; i < 24; i++ {
		available_hours = append(available_hours, i)
	}

	appointed_hours := []int{}
	for _, reservation := range reservations {
		for i := reservation.Start.Hour(); i < reservation.End.Hour(); i++ {
			appointed_hours = append(appointed_hours, i)
		}
	}

	for _, hour := range appointed_hours {
		available_hours = remove(available_hours, hour)
	}

	return available_hours
}

func remove(slice []int, s int) []int {
	for i, v := range slice {
		if v == s {
			slice = append(slice[:i], slice[i+1:]...)
			return slice
		}
	}
	return slice

}

func (h *Handlers) CreateRoom(c *gin.Context) {

	var input struct {
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
		Type     string `json:"type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	} else if !contains(Room_Types, input.Type) {
		c.JSON(400, gin.H{"error": "Invalid room type"})
		return
	}

	room := models.Room{
		Name:     input.Name,
		Capacity: input.Capacity,
		Type:     input.Type,
	}

	res := h.DB.Create(&room)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": res.Error.Error()})
		return
	}
	res.Commit()

	c.JSON(200, gin.H{"message": "Room created successfully",
		"room_id": room.ID,
	})

}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
