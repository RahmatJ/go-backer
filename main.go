package main

import (
	"fmt"
	"go-backer/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/backer_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to DB")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userInput := user.RegisterUserInput{
		Name:       "name from service",
		Email:      "mail@mail.com",
		Occupation: "occupation",
		Password:   "hai",
	}

	newUser, err := userService.RegisterUser(userInput)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(newUser)

	//router := gin.Default()
	//router.GET("/handler", handler)
	//
	//err := router.Run()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
}
