package main

import (
	"context"
	"os"

	"github.com/aslon1213/reservations-API/handlers"
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
	// db.AutoMigrate(&models.Resident{}, &models.Room{}, &models.Reservation{})

	Handlers = handlers.NewHandlers(ctx, db)

}

func GetRooms(c *gin.Context) {

}

func main() {
	Init()
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/rooms", Handlers.GetRooms)
	api.POST("/rooms/new", Handlers.CreateRoom)

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}

}
