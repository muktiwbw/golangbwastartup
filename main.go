package main

import (
	"bwastartup/entities/user"
	"bwastartup/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "sbxmukti:password1234@tcp(127.0.0.1:3306)/go_bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// log.Println("Connection to database is established.")
	log.Println("Connected to database.")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/emailCheck", userHandler.CheckEmailAvailability)

	router.Run()

}
