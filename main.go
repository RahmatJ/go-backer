package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-backer/auth"
	"go-backer/campaign"
	"go-backer/handler"
	"go-backer/helper"
	"go-backer/transaction"
	"go-backer/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/backer_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to DB")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()

	//add route for static file,
	// /images is for route, and ./images is for folder images
	router.Static("/images", "./images")

	//this for API versioning, should mind this
	//will automatically add '/api/v1' in front of each api
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)

	err = router.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}

//gin handler, using only gin.Context as parameter
//do some workaround for function that need input as gin handler
//that function still got another input, but return gin handler function
//with input only gin context

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("UNAUTHORIZED", http.StatusUnauthorized, "error", nil)
			//to abort this process when failed
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//token format: Bearer tokentokentoken
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("UNAUTHORIZED", http.StatusUnauthorized, "error", nil)
			//to abort this process when failed
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//get data from jwtToken
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("UNAUTHORIZED", http.StatusUnauthorized, "error", nil)
			//to abort this process when failed
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get userId from claim
		// by default int in claim will be converted to float64
		// so we need to cast it to int
		userID := int(claim["user_id"].(float64))

		authUser, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("UNAUTHORIZED", http.StatusUnauthorized, "error", nil)
			//to abort this process when failed
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", authUser)
	}
}
