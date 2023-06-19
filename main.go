package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aslon1213/reservations-API/handlers"
	"github.com/aslon1213/reservations-API/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB_CLIENT *gorm.DB
var Handlers *handlers.Handlers

func Init() {

	// load .env fil
	// export variables
	ctx := context.Background()
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URI")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB_CLIENT = db
	//create tables
	// db.AutoMigrate(&models.ReservationsbyDate{})
	// db.AutoMigrate(&models.Reservation{})
	db.AutoMigrate(&models.Resident{}, &models.Room{}, &models.Reservation{})

	Handlers = handlers.NewHandlers(ctx, db)

}

func create_some_rooms() {

	http_client := &http.Client{}
	type raw_input struct {
		Name     string
		Type     string
		Capacity int
	}

	var raw_inputs = []raw_input{
		{"mytaxi", "focus", 1},
		{"workly", "team", 5},
		{"express24", "conference", 15},
		{"Focus Room 1", "focus", 1},
		{"Focus Room 2", "focus", 1},
		{"Focus Room 3", "focus", 1},
	}

	for _, raw_input := range raw_inputs {
		// wait for 10 seconds
		time.Sleep(5 * time.Second)

		// make request to create room

		raw_date := []byte(fmt.Sprintf(`{"name":"%s","type":"%s","capacity":%d}`, raw_input.Name, raw_input.Type, raw_input.Capacity))

		req, err := http.NewRequest("POST", "http://localhost:8080/api/rooms/new", bytes.NewBuffer(raw_date))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http_client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

	}
}

func main() {

	time.Sleep(time.Second * 30)

	Init()
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/rooms", Handlers.GetRooms)
	api.POST("/rooms/new", Handlers.CreateRoom)
	api.GET("/rooms/:id", Handlers.GetRoom)
	api.GET("/rooms/:id/availability", Handlers.GetRoomAvailability)
	api.POST("/rooms/:id/book", Handlers.BookRoom)
	api.GET("/allreservs", Handlers.GetAllReservations)
	api.DELETE("/reservs/delete/all", Handlers.DeleteAllReservations)
	api.GET("/rooms/:id/reservations_full", Handlers.GetAllReservationsFullInfo)
	api.DELETE("/rooms/:id/unbook", Handlers.UnbookRoom)
	go create_some_rooms()
	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}

}
