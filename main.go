package main

import (
	"bwastartup/auth"
	"bwastartup/entities/campaign"
	"bwastartup/entities/user"
	"bwastartup/handlers"
	"bwastartup/helpers"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)

	userHandler := handlers.NewUserHandler(userService, authService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/emailCheck", userHandler.CheckEmailAvailability)
	api.POST("/updateAvatar", authorize(authService, userService), userHandler.UpdateAvatar)

	api.GET("/me/campaigns", authorize(authService, userService), campaignHandler.GetOwnCampaigns)

	api.GET("/campaigns", campaignHandler.GetAllCampaigns)
	api.GET("/campaigns/:campaign_id", campaignHandler.GetCampaignByID)

	router.Run()

}

func authorize(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer ") {
			data := helpers.APIResponse("Missing or invalid access token.", http.StatusUnauthorized, "error", gin.H{"error": errors.New("Access token is either missing or invalid.")})
			c.AbortWithStatusJSON(http.StatusUnauthorized, data)

			return
		}

		accessToken := strings.Split(authHeader, " ")[1]

		validatedToken, err := authService.ValidateToken(accessToken)

		if err != nil || !validatedToken.Valid {
			data := helpers.APIResponse("Missing or invalid access token.", http.StatusUnauthorized, "error", gin.H{"error": errors.New("Access token is either missing or invalid.")})
			c.AbortWithStatusJSON(http.StatusUnauthorized, data)

			return
		}

		claims, ok := validatedToken.Claims.(jwt.MapClaims)

		if !ok {
			data := helpers.APIResponse("Missing or invalid access token.", http.StatusUnauthorized, "error", gin.H{"error": errors.New("Access token is either missing or invalid.")})
			c.AbortWithStatusJSON(http.StatusUnauthorized, data)

			return
		}

		userID := int(claims["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			data := helpers.APIResponse("Missing or invalid access token.", http.StatusUnauthorized, "error", gin.H{"error": errors.New("Access token is either missing or invalid.")})
			c.AbortWithStatusJSON(http.StatusUnauthorized, data)

			return
		}

		c.Set("authUser", user)
	}
}
