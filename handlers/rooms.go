package handlers

import (
	"context"
	"fmt"
	"strconv"

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

	fmt.Println(rooms)
	rooms_with_pagination := map[int][]models.Room{}
	for i := 0; i <= len(rooms)/page_size_int; i++ {
		rooms_with_pagination[i] = rooms[i*page_size_int : (i+1)*page_size_int]
	}
	// c.JSON(200, gin.H{
	// 	"rooms": rooms,
	// })

	c.JSON(200, gin.H{"rooms": rooms_with_pagination[page_int-1]})

}

func (h *Handlers) GetRoom() {

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
